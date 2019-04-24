package crypto

import (
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/juju/errors"
)

type Signature struct {
	sig []byte
}

func NewSignature(sig []byte) *Signature {
	return &Signature{sig: sig}
}

func (p *Signature) ToHex() string {
	return hex.EncodeToString(p.Bytes())
}

func (p *Signature) Bytes() []byte {
	return p.sig
}

func isCanonical(data []byte) bool {
	return (data[1]&0x80) == 0 &&
		!(data[1] == 0 && (data[2]&0x80) == 0) &&
		(data[33]&0x80) == 0 &&
		!(data[33] == 0 && (data[34]&0x80) == 0)
}

func (p *Signature) IsCanonical() bool {
	return isCanonical(p.sig)
}

func Sign(data []byte, privKey *btcec.PrivateKey) (*Signature, error) {
	digest := sha256.Sum256(data)
	return signBufferSha256(digest[:], privKey)
}

func SignDigest(digest []byte, privKey *btcec.PrivateKey) (*Signature, error) {
	return signBufferSha256(digest, privKey)
}

func recoverPubKey(curve *btcec.KoblitzCurve, sig *btcec.Signature, msg []byte, iter int, doChecks bool) (*btcec.PublicKey, error) {
	// 1.1 x = (n * i) + r
	Rx := new(big.Int).Mul(curve.Params().N,
		new(big.Int).SetInt64(int64(iter/2)))
	Rx.Add(Rx, sig.R)
	if Rx.Cmp(curve.Params().P) != -1 {
		return nil, errors.New("calculated Rx is larger than curve P")
	}

	// convert 02<Rx> to point R. (step 1.2 and 1.3). If we are on an odd
	// iteration then 1.6 will be done with -R, so we calculate the other
	// term when uncompressing the point.
	Ry, err := decompressPoint(curve, Rx, iter%2 == 1)
	if err != nil {
		return nil, err
	}

	// 1.4 Check n*R is point at infinity
	if doChecks {
		nRx, nRy := curve.ScalarMult(Rx, Ry, curve.Params().N.Bytes())
		if nRx.Sign() != 0 || nRy.Sign() != 0 {
			return nil, errors.New("n*R does not equal the point at infinity")
		}
	}

	// 1.5 calculate e from message using the same algorithm as ecdsa
	// signature calculation.
	e := hashToInt(msg, curve)

	// Step 1.6.1:
	// We calculate the two terms sR and eG separately multiplied by the
	// inverse of r (from the signature). We then add them to calculate
	// Q = r^-1(sR-eG)
	invr := new(big.Int).ModInverse(sig.R, curve.Params().N)

	// first term.
	invrS := new(big.Int).Mul(invr, sig.S)
	invrS.Mod(invrS, curve.Params().N)
	sRx, sRy := curve.ScalarMult(Rx, Ry, invrS.Bytes())

	// second term.
	e.Neg(e)
	e.Mod(e, curve.Params().N)
	e.Mul(e, invr)
	e.Mod(e, curve.Params().N)
	minuseGx, minuseGy := curve.ScalarBaseMult(e.Bytes())

	// TODO: this would be faster if we did a mult and add in one
	// step to prevent the jacobian conversion back and forth.
	Qx, Qy := curve.Add(sRx, sRy, minuseGx, minuseGy)

	return &btcec.PublicKey{
		Curve: curve,
		X:     Qx,
		Y:     Qy,
	}, nil
}

func signBufferSha256(bufSha256 []byte, key *btcec.PrivateKey) (*Signature, error) {
	if len(bufSha256) != 32 {
		return nil, errors.New("buf_sha256: 32 byte buffer required")
	}

	for nonce := 1; ; nonce++ {
		sig, err := signRFC6979(key, bufSha256, nonce)

		if err != nil {
			return nil, errors.Annotate(err, "signRFC6979")
		}

		der := sig.Serialize()
		lenR := der[3]
		lenS := der[5+lenR]

		if lenR == 32 && lenS == 32 {
			//////////////////////////////////////
			// bitcoind checks the bit length of R and S here. The ecdsa signature
			// algorithm returns R and S mod N therefore they will be the bitsize of
			// the curve, and thus correctly sized.
			curve := btcec.S256()
			maxCounter := 4
			for i := 0; i < maxCounter; i++ {
				pk, err := recoverPubKey(curve, sig, bufSha256, i, true)

				if err == nil && pk.X.Cmp(key.X) == 0 && pk.Y.Cmp(key.Y) == 0 {
					byteSize := curve.BitSize / 8
					result := make([]byte, 1, 2*byteSize+1)
					result[0] = 27 + byte(i)
					if true { // isCompressedKey
						result[0] += 4
					}
					// Not sure this needs rounding but safer to do so.
					curvelen := (curve.BitSize + 7) / 8

					// Pad R and S to curvelen if needed.
					bytelen := (sig.R.BitLen() + 7) / 8
					if bytelen < curvelen {
						result = append(result, make([]byte, curvelen-bytelen)...)
					}
					result = append(result, sig.R.Bytes()...)

					bytelen = (sig.S.BitLen() + 7) / 8
					if bytelen < curvelen {
						result = append(result, make([]byte, curvelen-bytelen)...)
					}
					result = append(result, sig.S.Bytes()...)
					return NewSignature(result), nil
				}
			}
		}
	}
}

func decompressPoint(curve *btcec.KoblitzCurve, x *big.Int, ybit bool) (*big.Int, error) {
	// TODO: This will probably only work for secp256k1 due to
	// optimizations.

	// Y = +-sqrt(x^3 + B)
	x3 := new(big.Int).Mul(x, x)
	x3.Mul(x3, x)
	x3.Add(x3, curve.Params().B)

	// now calculate sqrt mod p of x2 + B
	// This code used to do a full sqrt based on tonelli/shanks,
	// but this was replaced by the algorithms referenced in
	// https://bitcointalk.org/index.php?topic=162805.msg1712294#msg1712294
	y := new(big.Int).Exp(x3, curve.QPlus1Div4(), curve.Params().P)

	if ybit != isOdd(y) {
		y.Sub(curve.Params().P, y)
	}
	if ybit != isOdd(y) {
		return nil, fmt.Errorf("ybit doesn't match oddness")
	}
	return y, nil
}

func isOdd(a *big.Int) bool {
	return a.Bit(0) == 1
}

func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
	orderBits := c.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	excess := len(hash)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}
