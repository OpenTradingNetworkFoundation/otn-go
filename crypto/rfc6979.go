package crypto

import (
	"bytes"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"hash"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/juju/errors"
)

var (
	// Used in RFC6979 implementation when testing the nonce for correctness
	one = big.NewInt(1)

	// oneInitializer is used to fill a byte slice with byte 0x01.  It is provided
	// here to avoid the need to create it multiple times.
	oneInitializer = []byte{0x01}
)

// signRFC6979 generates a deterministic ECDSA signature according to RFC 6979 and BIP 62.
func signRFC6979(privateKey *btcec.PrivateKey, hash []byte, nonce int) (*btcec.Signature, error) {
	privkey := privateKey.ToECDSA()
	N := btcec.S256().N
	halfOrder := new(big.Int).Rsh(btcec.S256().N, 1) // btcec.S256().halfOrder
	k := nonceRFC6979(privkey.D, hash, nonce)
	inv := new(big.Int).ModInverse(k, N)
	r, _ := privkey.Curve.ScalarBaseMult(k.Bytes())
	if r.Cmp(N) == 1 {
		r.Sub(r, N)
	}

	if r.Sign() == 0 {
		return nil, errors.New("calculated R is zero")
	}

	e := hashToInt(hash, privkey.Curve)
	s := new(big.Int).Mul(privkey.D, r)
	s.Add(s, e)
	s.Mul(s, inv)
	s.Mod(s, N)

	if s.Cmp(halfOrder) == 1 {
		s.Sub(N, s)
	}
	if s.Sign() == 0 {
		return nil, errors.New("calculated S is zero")
	}
	return &btcec.Signature{R: r, S: s}, nil
}

// nonceRFC6979 generates an ECDSA nonce (`k`) deterministically according to RFC 6979.
// It takes a 32-byte hash as an input and returns 32-byte nonce to be used in ECDSA algorithm.
func nonceRFC6979(privkey *big.Int, hash []byte, nonce int) *big.Int {
	if nonce > 0 {
		nonceA := make([]byte, 4)
		binary.BigEndian.PutUint32(nonceA, uint32(nonce))
		hs := sha256.New()
		hs.Write(append(hash, nonceA...))
		hash = hs.Sum(nil)
	}

	curve := btcec.S256()
	q := curve.Params().N
	x := privkey
	alg := sha256.New

	qlen := q.BitLen()
	holen := alg().Size()
	rolen := (qlen + 7) >> 3
	bx := append(int2octets(x, rolen), bits2octets(hash, curve, rolen)...)

	// Step B
	v := bytes.Repeat(oneInitializer, holen)

	// Step C (Go zeroes the all allocated memory)
	k := make([]byte, holen)

	// Step D
	k = mac(alg, k, append(append(v, 0x00), bx...))

	// Step E
	v = mac(alg, k, v)

	// Step F
	k = mac(alg, k, append(append(v, 0x01), bx...))

	// Step G
	v = mac(alg, k, v)

	// Step H
	for {
		// Step H1
		var t []byte

		// Step H2
		for len(t)*8 < qlen {
			v = mac(alg, k, v)
			t = append(t, v...)
		}

		// Step H3
		secret := hashToInt(t, curve)
		if secret.Cmp(one) >= 0 && secret.Cmp(q) < 0 {
			return secret
		}
		k = mac(alg, k, append(v, 0x00))
		v = mac(alg, k, v)
	}
}

// mac returns an HMAC of the given key and message.
func mac(alg func() hash.Hash, k, m []byte) []byte {
	h := hmac.New(alg, k)
	h.Write(m)
	return h.Sum(nil)
}

// https://tools.ietf.org/html/rfc6979#section-2.3.3
func int2octets(v *big.Int, rolen int) []byte {
	out := v.Bytes()

	// left pad with zeros if it's too short
	if len(out) < rolen {
		out2 := make([]byte, rolen)
		copy(out2[rolen-len(out):], out)
		return out2
	}

	// drop most significant bytes if it's too long
	if len(out) > rolen {
		out2 := make([]byte, rolen)
		copy(out2, out[len(out)-rolen:])
		return out2
	}

	return out
}

// https://tools.ietf.org/html/rfc6979#section-2.3.4
func bits2octets(in []byte, curve elliptic.Curve, rolen int) []byte {
	z1 := hashToInt(in, curve)
	z2 := new(big.Int).Sub(z1, curve.Params().N)
	if z2.Sign() < 0 {
		return int2octets(z1, rolen)
	}
	return int2octets(z2, rolen)
}
