package objects

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadTestdata(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	return data
}

func TestFullAccountParse(t *testing.T) {
	js := loadTestdata("testdata/fullaccount.json")
	var acc FullAccount
	err := json.Unmarshal(js, &acc)
	require.NoError(t, err)

	assert.Equal(t, *NewGrapheneID("1.2.18"), acc.Account.ID)
	assert.Equal(t, "gateway", acc.Account.Name)
	assert.Len(t, acc.Balances, 3)
	assert.Equal(t, Int64(1000123), acc.Balances[0].Balance)
	assert.Equal(t, Int64(100), acc.Balances[1].Balance)
	assert.Equal(t, Int64(200), acc.Balances[2].Balance)

	// call orders
	assert.Len(t, acc.CallOrders, 2)
	assert.Equal(t, UInt64(1061318752610), acc.CallOrders[0].Collateral)
	assert.Equal(t, UInt64(652336746294), acc.CallOrders[1].Collateral)
}
