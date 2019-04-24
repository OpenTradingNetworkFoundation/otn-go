package secrets

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

func LoadKeyFromBlob(data []byte) (interface{}, error) {
	block, _ := pem.Decode(data)

	if block == nil {
		return nil, fmt.Errorf("Key not found")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "RSA PUBLIC KEY":
		return x509.ParsePKIXPublicKey(block.Bytes)
	default:
		return nil, errors.New("Unsupported key type")
	}
}

func LoadPublicKey(data []byte) (*rsa.PublicKey, error) {
	key, err := LoadKeyFromBlob(data)
	if err != nil {
		return nil, err
	}

	v, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("Unsupported key type")
	}

	return v, nil
}

func LoadPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	key, err := LoadKeyFromBlob(data)
	if err != nil {
		return nil, err
	}

	v, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("Unsupported key type")
	}

	return v, nil
}
