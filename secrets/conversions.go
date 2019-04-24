package secrets

import (
	"crypto/rsa"
	"encoding/base64"

	"github.com/juju/errors"
)

func GetPrivateKey(storage SecretStorage, path string) (*rsa.PrivateKey, error) {
	secret, err := storage.ReadStringValue(path)
	if err != nil {
		return nil, err
	}

	data, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, errors.Annotate(err, "base64 decode")
	}

	return LoadPrivateKey(data)
}

func GetPublicKey(storage SecretStorage, path string) (*rsa.PublicKey, error) {
	secret, err := storage.ReadStringValue(path)
	if err != nil {
		return nil, err
	}

	data, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, errors.Annotate(err, "base64 decode")
	}

	return LoadPublicKey(data)
}
