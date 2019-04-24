package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignatureCanonical(t *testing.T) {
	pvtKey, _ := GetPrivateKey("5JNYPCUHiS3itj4tnyWFKkgnM1VtEJfzjoGGTsZFYxCTWm6tDTC")

	for i := 0; i < 1000; i++ {
		data := []byte(fmt.Sprintf("pass%d", i))
		sig, err := Sign(data, pvtKey)
		assert.NoErrorf(t, err, "Failed to sign")
		t.Log(sig.ToHex())
		assert.True(t, sig.IsCanonical(), "IsCanonical")
	}

	data := []byte("pass")
	sig, err := Sign(data, pvtKey)
	assert.NoErrorf(t, err, "Failed to sign")
	fmt.Println(sig.ToHex())
}
