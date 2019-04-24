package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/opentradingnetworkfoundation/otn-go/crypto"
	"github.com/opentradingnetworkfoundation/otn-go/objects"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

func hashWithChain(chainID string, data []byte) ([]byte, error) {
	rawChainID, err := hex.DecodeString(chainID)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Write(rawChainID)
	buf.Write(data)
	hash := sha256.Sum256(buf.Bytes())

	return hash[:], nil
}

func decodeFromFile(filename string, obj interface{}) error {
	var inputData io.Reader

	if filename == "-" {
		inputData = os.Stdin
	} else {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalln(err)
		}
		inputData = file
	}

	data, err := ioutil.ReadAll(inputData)

	if err != nil {
		log.Fatal(err)
	}

	// Decode json.
	err = json.Unmarshal([]byte(data), obj)
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func main() {
	inputFile := flag.String("f", "-", "Input file")
	encodeBinary := flag.Bool("b", false, "Output as binary (default is hex)")
	signKey := flag.String("s", "", "Sign transaction with this private key")
	memoKey := flag.String("m", "", "Decode memo with this key")
	chainID := flag.String("chain", "7f43b3c845d6b7b868eebcdd3df13f031eba425c573db44845f766205a909d74", "Chain id")

	flag.Parse()

	if len(*memoKey) > 0 {
		var memo objects.Memo
		decodeFromFile(*inputFile, &memo)

		key, err := crypto.GetPrivateKey(*memoKey)
		if err != nil {
			log.Fatal("Invalid key")
		}

		msg, err := memo.GetMessage(key)
		if err != nil {
			log.Fatal("Failed to decode: ", err)
		}

		fmt.Println("Message: ", string(msg))
		return
	}

	var tx objects.Transaction
	decodeFromFile(*inputFile, &tx)

	var buf bytes.Buffer
	enc := util.NewTypeEncoder(&buf)

	if err := tx.Marshal(enc); err != nil {
		log.Fatalln(err)
	}

	if len(*signKey) > 0 {
		keys := []string{*signKey}
		binChainID, _ := hex.DecodeString(*chainID)
		err := tx.Sign(keys, objects.ChainID(binChainID))
		if err != nil {
			log.Fatalln("Failed to sign:", err)
		}

		fmt.Println(tx.Signatures)
	}

	bin := buf.Bytes()

	if *encodeBinary {
		os.Stdout.Write(bin)
	} else {
		fmt.Println("Binary:", hex.EncodeToString(bin))
		hash := sha256.Sum256(bin)
		fmt.Println("SHA256:", hex.EncodeToString(hash[:]))

		hashWithChainID, err := hashWithChain(*chainID, bin)
		if err != nil {
			log.Fatalln("Failed to hash with chain id")
		}

		fmt.Println("SHA256 with chain:", hex.EncodeToString(hashWithChainID))
	}
}
