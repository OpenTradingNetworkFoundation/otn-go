package crypto

import (
	"bytes"
	"crypto/aes"
	"errors"
)

func RemovePadding(b []byte) ([]byte, error) {
	l := int(b[len(b)-1])
	if l > aes.BlockSize {
		return nil, errors.New("Invalid padding")
	}

	return b[:len(b)-l], nil
}

func AddPadding(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}
