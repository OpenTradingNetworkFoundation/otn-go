package objects

import (
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/opentradingnetworkfoundation/otn-go/util"
)

func createVBTx(op Operation) *Transaction {
	tx := &Transaction{
		RefBlockNum:    30266,
		RefBlockPrefix: 1338368302,
		Operations:     Operations{op},
		Extensions:     Extensions{},
	}

	tx.Expiration.UnmarshalJSON([]byte("\"2018-05-10T08:42:20\""))
	return tx
}

func createVBOp(policy VestingPolicyInitializer) *VestingBalanceCreateOperation {
	op := &VestingBalanceCreateOperation{
		Creator: *NewGrapheneID(ObjectID("1.2.20")),
		Owner:   *NewGrapheneID(ObjectID("1.2.21")),
		Fee: AssetAmount{
			Asset:  CoreAssetID,
			Amount: 500,
		},
		Amount: AssetAmount{
			Asset:  CoreAssetID,
			Amount: 10003,
		},
		Policy: policy,
	}

	return op
}

func checkVBTx(t *testing.T, tx *Transaction, expected string) {
	js, err := json.Marshal(tx)
	assert.NoError(t, err)
	t.Log(string(js))

	bin, err := util.EncodeToBytes(tx)
	assert.NoError(t, err)

	binhex := hex.EncodeToString(bin)
	t.Logf("Binary: %s", binhex)

	assert.Equal(t, expected, binhex)
}

func TestVBCreateCDD(t *testing.T) {
	policy := VestingPolicyInitializer{
		Type: VestingPolicyTypeCDD,
		Value: CDDVestingPolicyInitializer{
			StartClaim:     NewTime(time.Date(2020, 1, 1, 12, 00, 0, 0, time.UTC)),
			VestingSeconds: 9000,
		},
	}

	tx := createVBTx(createVBOp(policy))
	checkVBTx(t, tx, "3a762ee1c54fec05f45a0120f40100000000000000141513270000000000000001c0890c5e2823000000")
}

func TestVBCreateLinear(t *testing.T) {
	policy := VestingPolicyInitializer{
		Type: VestingPolicyTypeLinear,
		Value: LinearVestingPolicyInitializer{
			BeginTimestamp:         NewTime(time.Date(2025, 1, 1, 12, 00, 0, 0, time.UTC)),
			VestingCliffSeconds:    600,
			VestingDurationSeconds: 3200,
		},
	}
	tx := createVBTx(createVBOp(policy))
	checkVBTx(t, tx, "3a762ee1c54fec05f45a0120f40100000000000000141513270000000000000000402e756758020000800c000000")
}
