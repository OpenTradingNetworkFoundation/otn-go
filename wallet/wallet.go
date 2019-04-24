package wallet

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/crypto"
)

type Wallet interface {
	AddPrivateKeys(wifKeys []string) error
	GetKeys() map[string]*btcec.PrivateKey
}

type otnWallet struct {
	keys map[string]*btcec.PrivateKey
}

func (w *otnWallet) AddPrivateKeys(wifKeys []string) error {
	for _, wif := range wifKeys {
		key, err := crypto.GetPrivateKey(wif)
		if err != nil {
			return errors.Annotate(err, "GetPrivateKey")
		}
		pubKey := crypto.GetPublicKeyString("OTN", key.PubKey())
		w.keys[pubKey] = key
	}
	return nil
}

func (w *otnWallet) GetKeys() map[string]*btcec.PrivateKey {
	return w.keys
}

func NewWallet() Wallet {
	return &otnWallet{
		keys: make(map[string]*btcec.PrivateKey),
	}
}
