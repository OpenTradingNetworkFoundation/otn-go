package crypto

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAddress(t *testing.T) {
	pubKey := "OTN79eaPCGjohXpu459EKjCn1EWpVuNYFj9LHtyZPQZ21QPTEQY49"
	pvtKey, _ := GetPrivateKey("5JNYPCUHiS3itj4tnyWFKkgnM1VtEJfzjoGGTsZFYxCTWm6tDTC")
	key := ParsePublicKey(pubKey)
	if key == nil {
		t.Error("Failed to parse address")
	}

	t.Log(pvtKey.PubKey())
	t.Log(key)

	if bytes.Compare(pvtKey.PubKey().SerializeCompressed(), key.SerializeCompressed()) != 0 {
		t.Error("Keys differ")
	}

	newPubKey := GetPublicKeyString("OTN", pvtKey.PubKey())

	if newPubKey != pubKey {
		t.Error("Public keys differ")
	}
}

func TestInvalidKeys(t *testing.T) {
	// check broken checksum
	_, err := GetPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvF53")
	assert.NotEqual(t, nil, err)

	// check broken value
	_, err = GetPrivateKey("5KQXrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	assert.NotEqual(t, nil, err)
}

func TestConversions(t *testing.T) {
	// parse correct key
	pvtKey, _ := GetPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	pubKey := pvtKey.PubKey()

	pubKeyStr := GetPublicKeyString("OTN", pubKey)
	assert.Equal(t, "OTN6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV", pubKeyStr)

	address := GetAddressString("OTN", pubKey)
	assert.Equal(t, "OTNFAbAx7yuxt725qSZvfwWqkdCwp9ZnUama", address)
}

func TestSharedSecret(t *testing.T) {
	fromKey, _ := GetPrivateKey("5JB331gjg4AheTKNuiUYWL37rC6Vbi943agL4o1fH4uoJmDtJiB")
	toKey, _ := GetPrivateKey("5HzbeKJPi1DJGC6ShPju1Muj2TTMRkc8Cs1DHJuyA38YjTkfNgB")

	expectedSecret := "b6038ff9c435238a097981a90e90921c5bf99f11d719491bd6613a68f15d29ca5ac627d07b26c98c4a92dc51b1fe991f02f7e705488df1838b5d62a86acf7e9f"

	secret := GetSharedSecret(fromKey, toKey.PubKey())
	secretHex := hex.EncodeToString(secret)
	assert.Equal(t, expectedSecret, secretHex)

	secretReverse := GetSharedSecret(toKey, fromKey.PubKey())
	secretHex = hex.EncodeToString(secretReverse)
	assert.Equal(t, expectedSecret, secretHex)
}
