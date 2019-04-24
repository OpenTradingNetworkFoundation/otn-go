package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/chaincfg"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/opentradingnetworkfoundation/otn-go/crypto"
)

var prefix string

type Keypair struct {
	pvt    *btcec.PrivateKey
	prefix string
}

func (k *Keypair) GetPublicKeyString() string {
	return crypto.GetPublicKeyString(k.prefix, k.pvt.PubKey())
}

func (k *Keypair) GetAddressString() string {
	return crypto.GetAddressString(k.prefix, k.pvt.PubKey())
}

func (k *Keypair) GetPrivateKeyString() string {
	wif, err := btcutil.NewWIF(k.pvt, &chaincfg.MainNetParams, false)
	if err != nil {
		log.Fatal("Failed to create WIF:", err)
	}

	return wif.String()
}

func NewKeypair(prefix string) (*Keypair, error) {
	pvt, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}

	return &Keypair{pvt, prefix}, nil
}

func generateKeys(count int) {
	for i := 0; i < count; i++ {
		kp, err := NewKeypair(prefix)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\t%s\t%s\n",
			kp.GetAddressString(),
			kp.GetPublicKeyString(),
			kp.GetPrivateKeyString())
	}
}

func decodeKey(key string) {
	wif, err := btcutil.DecodeWIF(key)
	if err != nil {
		log.Fatal(err)
	}
	pubKey := wif.PrivKey.PubKey()
	pubstr := crypto.GetPublicKeyString(prefix, pubKey)
	address := crypto.GetAddressString(prefix, pubKey)
	fmt.Println(pubstr, address)
}

func main() {
	flag.StringVar(&prefix, "p", "OTN", "Address prefix")
	count := flag.Int("n", 1, "Number of items to generate")
	parseKey := flag.String("parse", "", "Parse private key")

	flag.Parse()

	if len(*parseKey) > 0 {
		decodeKey(*parseKey)
		return
	}

	generateKeys(*count)
}
