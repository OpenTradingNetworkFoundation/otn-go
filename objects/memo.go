package objects

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/crypto"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

// Memo field
// Encoded as AES-256 CBC with checksum
type Memo struct {
	From    PublicKey `json:"from"`
	To      PublicKey `json:"to"`
	Nonce   UInt64    `json:"nonce"`
	Message Binary    `json:"message"`
}

const memoChecksumSize = 4
const nullPublicKeyString = "OTN1111111111111111111111111111111114T1Anm"

func makeKeys(nonce uint64, from *btcec.PrivateKey, to *btcec.PublicKey) ([]byte, []byte) {
	sharedSecret := crypto.GetSharedSecret(from, to)
	secret := fmt.Sprintf("%d%s", nonce, hex.EncodeToString(sharedSecret))
	secretHash := sha512.Sum512([]byte(secret))

	key := secretHash[:32]
	iv := secretHash[32:48]

	return key, iv
}

func decodeMessage(decryptedMessage []byte) (string, error) {
	checksum := decryptedMessage[:memoChecksumSize]
	digest := sha256.Sum256(decryptedMessage[memoChecksumSize:])
	if bytes.Compare(checksum, digest[:memoChecksumSize]) != 0 {
		return "", errors.New("Checksum mismatch")
	}

	message := string(decryptedMessage[memoChecksumSize:])
	return message, nil
}

func (p *Memo) Marshal(enc *util.TypeEncoder) error {
	return util.BinaryEncodeStruct(enc, p)
}

func (p *Memo) SetMessage(message, prefix string, from *btcec.PrivateKey, to *btcec.PublicKey, nonce uint64) error {
	if nonce == 0 {
		b := make([]byte, 8)
		rand.Read(b)
		nonce = binary.LittleEndian.Uint64(b)
	}

	messageBytes := []byte(message)
	checksum := sha256.Sum256(messageBytes)
	messageWithChecksum := append(checksum[:memoChecksumSize], messageBytes...)
	if from == nil || to == nil {
		p.To = NewPublicKey(nullPublicKeyString)
		p.From = NewPublicKey(nullPublicKeyString)
		p.Message = messageWithChecksum
		p.Nonce = UInt64(nonce)
		return nil
	}

	key, iv := makeKeys(nonce, from, to)

	// encrypted text is checksum + message
	plaintext := crypto.AddPadding(messageWithChecksum)
	encryptedMessage, err := crypto.EncryptCBC(plaintext, iv, key)
	if err != nil {
		return errors.Annotate(err, "EncryptCBC")
	}

	p.Message = Binary(encryptedMessage)
	p.From = NewPublicKey(crypto.GetPublicKeyString(prefix, from.PubKey()))
	p.To = NewPublicKey(crypto.GetPublicKeyString(prefix, to))
	p.Nonce = UInt64(nonce)
	return nil
}

func (p *Memo) GetMessage(to *btcec.PrivateKey) (string, error) {
	if !p.From.Valid() || p.From.String() == nullPublicKeyString {
		return decodeMessage(p.Message)
	}

	from := p.From.GetPublicKey()
	if from == nil {
		return "", errors.New("Invalid From key")
	}
	toPub := p.To.GetPublicKey()
	if toPub == nil {
		return "", errors.New("Invalid To key")
	}

	if !toPub.IsEqual(to.PubKey()) {
		return "", errors.New("'To' key mismatch")
	}

	key, iv := makeKeys(uint64(p.Nonce), to, from)
	decryptedMessage, err := crypto.DecryptCBC(p.Message, iv, key)
	if err != nil {
		return "", errors.Annotate(err, "DecryptCBC")
	}

	decryptedMessage, err = crypto.RemovePadding(decryptedMessage)
	if err != nil {
		return "", errors.Annotate(err, "RemovePadding")
	}

	return decodeMessage(decryptedMessage)
}
