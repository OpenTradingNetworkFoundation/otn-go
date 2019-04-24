package api

import (
	"encoding/hex"
	"log"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/objects"
)

func ConstructTransaction(p DatabaseAPI, feeAsset objects.GrapheneObject, ops ...objects.Operation) (*objects.Transaction, error) {
	operations := objects.Operations(ops)
	operations.SetFeeAsset(*objects.NewGrapheneID(feeAsset.Id()))
	fees, err := p.GetRequiredFees(operations, feeAsset)
	if err != nil {
		return nil, errors.Annotate(err, "GetRequiredFees")
	}

	if err := operations.ApplyFees(fees); err != nil {
		return nil, errors.Annotate(err, "ApplyFees")
	}

	props, err := p.GetDynamicGlobalProperties()
	if err != nil {
		return nil, errors.Annotate(err, "GetDynamicGlobalProperties")
	}

	tx, err := objects.NewTransactionWithBlockData(props)
	if err != nil {
		return nil, errors.Annotate(err, "NewTransaction")
	}

	tx.Operations = operations
	return tx, nil
}

func SignTransaction(p DatabaseAPI, pvtKeys map[string]*btcec.PrivateKey, tx *objects.Transaction) error {
	pubKeys, err := p.GetPotentialSignatures(tx)
	if err != nil {
		return errors.Annotate(err, "GetPotentialSignatures")
	}

	pubKeys, err = p.GetRequiredSignatures(tx, pubKeys)
	if err != nil {
		return errors.Annotate(err, "GetRequiredSignatures")
	}

	reqKeys := make([]*btcec.PrivateKey, len(pubKeys))
	for i, pk := range pubKeys {
		wif, ok := pvtKeys[pk.String()]
		if !ok {
			return errors.New("Don't have required key")
		}
		reqKeys[i] = wif
	}

	if err := tx.SignWithKeys(reqKeys, p.GetChainID()); err != nil {
		return errors.Annotate(err, "Sign")
	}

	return nil
}

func SignAndBroadcast(api BitsharesAPI, keys map[string]*btcec.PrivateKey, feeAsset objects.GrapheneObject, ops ...objects.Operation) (*objects.Transaction, error) {
	dbAPI, err := api.DatabaseAPI()
	if err != nil {
		return nil, err
	}

	broadcastAPI, err := api.BroadcastAPI()
	if err != nil {
		return nil, err
	}

	tx, err := ConstructTransaction(dbAPI, feeAsset, ops...)
	if err != nil {
		return nil, errors.Annotate(err, "ConstructTransaction")
	}

	// check if there is a transaction with the same txid
	for {
		txID := tx.ID()
		oldTx, err := dbAPI.GetRecentTransactionByID(txID)
		if err != nil {
			return nil, errors.Annotate(err, "GetRecentTransactionByID")
		}

		if oldTx == nil {
			break
		}

		log.Printf("Transaction with the same txid=%s exists, adjusting expiration time",
			hex.EncodeToString(txID))

		// change transaction to generate new txid
		tx.Expiration.Time = tx.Expiration.Time.Add(time.Second)
	}

	err = SignTransaction(dbAPI, keys, tx)
	if err != nil {
		return nil, errors.Annotate(err, "SignTransaction")
	}

	err = broadcastAPI.BroadcastTransaction(tx)
	if err != nil {
		return nil, errors.Annotate(err, "BroadcastTransaction")
	}

	return tx, err
}
