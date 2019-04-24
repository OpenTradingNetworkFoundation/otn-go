package crypto

import (
	"bytes"
	"crypto/sha512"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"

	"github.com/juju/errors"
)

func ripemd160Sum(data []byte) []byte {
	h := ripemd160.New()
	h.Write(data)
	return h.Sum(nil)
}

func GetSharedSecret(from *btcec.PrivateKey, to *btcec.PublicKey) []byte {
	x, _ := to.Curve.ScalarMult(to.X, to.Y, from.D.Bytes())
	data := x.Bytes()
	// pad with zeroes to 32 bytes
	if len(data) < 32 {
		prefix := make([]byte, 32-len(data))
		data = append(prefix, data...)
	}
	digest := sha512.Sum512(data)
	return digest[:]
}

// Decode can be used to turn WIF into a raw private key (32 bytes).
func Decode(wif string) ([]byte, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, errors.Annotate(err, "failed to decode WIF")
	}

	return w.PrivKey.Serialize(), nil
}

func GetPrivateKey(wif string) (*btcec.PrivateKey, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, errors.Annotate(err, "failed to decode WIF")
	}

	return w.PrivKey, nil
}

// GetPublicKey returns the public key associated with the given WIF
// in the 33-byte compressed format.
func GetPublicKey(wif string) ([]byte, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, errors.Annotate(err, "failed to decode WIF")
	}

	return w.PrivKey.PubKey().SerializeCompressed(), nil
}

// ParsePublicKey parses address to public key, returns nil on error
func ParsePublicKey(key string) *btcec.PublicKey {
	if len(key) < 8 {
		return nil
	}
	pkn1 := key[3:]
	b58 := base58.Decode(pkn1)
	chs := b58[len(b58)-4:]
	pkn2 := b58[0 : len(b58)-4]
	checksumHash := ripemd160.New()
	checksumHash.Write(pkn2)
	nchs := checksumHash.Sum(nil)[0:4]
	if bytes.Equal(chs, nchs) {
		pkn3, err := btcec.ParsePubKey(pkn2, btcec.S256())
		if err == nil {
			return pkn3
		}
	}
	return nil
}

func GetPublicKeyString(prefix string, key *btcec.PublicKey) string {
	keyData := key.SerializeCompressed()
	nchs := ripemd160Sum(keyData)[0:4]
	keyData = append(keyData, nchs...)
	return prefix + base58.Encode(keyData)
}

func GetAddressString(prefix string, key *btcec.PublicKey) string {
	longHash := sha512.Sum512(key.SerializeCompressed())
	shortHash := ripemd160Sum(longHash[:])
	checksum := ripemd160Sum(shortHash)[0:4]
	return prefix + base58.Encode(append(shortHash, checksum...))
}
