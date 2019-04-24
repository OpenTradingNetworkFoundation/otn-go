package objects

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/opentradingnetworkfoundation/otn-go/objects/testdata"
)

func TestOperationsJSONParse(t *testing.T) {
	var txs []Transaction
	err := json.Unmarshal([]byte(testdata.ManyOperationsJSON), &txs)
	require.NoError(t, err)
	assert.Len(t, txs, 28)

	// do roundtrip encode-decode
	data, err := json.Marshal(txs)
	require.NoError(t, err)

	var tx2 []Transaction
	require.NoError(t, json.Unmarshal(data, &tx2))
	assert.Len(t, tx2, 28)

	data2, err := json.Marshal(tx2)
	require.NoError(t, err)

	require.True(t, bytes.Compare(data, data2) == 0)
}

func BenchmarkOperationsJSONParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var txs []Transaction
		json.Unmarshal([]byte(testdata.ManyOperationsJSON), &txs)
	}
}
