package objects_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/opentradingnetworkfoundation/otn-go/objects"
	"github.com/opentradingnetworkfoundation/otn-go/objects/testdata"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

func TestLimitOrderCancelOperation(t *testing.T) {
	var buf bytes.Buffer
	enc := util.NewTypeEncoder(&buf)

	tx := objects.NewTransaction()
	tx.RefBlockNum = 555
	tx.RefBlockPrefix = 3333333

	assert.NoError(t, tx.Expiration.UnmarshalJSON([]byte(`"2006-01-02T15:04:05"`)))

	op := objects.NewLimitOrderCancelOperation(
		*objects.NewGrapheneID("1.7.69314"),
		*objects.NewGrapheneID("1.2.456"),
	)

	op.Order = *objects.NewGrapheneID("1.7.123")
	op.Fee = objects.AssetAmount{
		Amount: 1000,
		Asset:  *objects.NewGrapheneID("1.3.789"),
	}

	tx.Operations = append(tx.Operations, op)

	js, _ := json.Marshal(tx)

	t.Log(string(js))
	assert.NoError(t, enc.Encode(tx))

	res := hex.EncodeToString(buf.Bytes())
	assert.Equal(t, "2b02d5dc3200e540b9430102e8030000000000009506c8037b0000", res)

	txDigest, _ := tx.Digest()
	txID := tx.ID()

	assert.Equal(t, "360adad70b08fffbb7dae94d239b04f778854f46e4b7b49660cda903d1b56f0b", hex.EncodeToString(txDigest))
	assert.Equal(t, "360adad70b08fffbb7dae94d239b04f778854f46", hex.EncodeToString(txID))
}

func TestOperationResults(t *testing.T) {
	var tx objects.Transaction
	err := json.Unmarshal([]byte(testdata.TransactionWithResults), &tx)
	require.NoError(t, err)
	require.Len(t, tx.OperationResults, 6)

	result := objects.AssetAmount{Amount: 100000000000, Asset: *objects.NewGrapheneID("1.3.0")}
	assert.Equal(t, objects.OperationResultType_Asset, tx.OperationResults[0].Type)
	assert.Equal(t, &result, tx.OperationResults[0].Result)

	assert.Equal(t, objects.OperationResultType_ObjectID, tx.OperationResults[5].Type)
	assert.Equal(t, objects.NewGrapheneID("1.7.4442342"), tx.OperationResults[5].Result)
}
