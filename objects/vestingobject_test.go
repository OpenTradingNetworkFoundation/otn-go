package objects

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadTestJsonObject(t *testing.T, path string, obj interface{}) {
	js := loadTestdata(path)
	err := json.Unmarshal(js, obj)
	require.NoError(t, err)
}

func TestVBLoadCDD(t *testing.T) {
	var vb VestingBalance
	loadTestJsonObject(t, "testdata/vb_cdd.json", &vb)

	assert.Equal(t, int8(VestingPolicyTypeCDD), vb.Policy.Type)
	assert.IsType(t, &CDDVestingPolicy{}, vb.Policy.Value)

	cdd := vb.Policy.Value.(*CDDVestingPolicy)
	assert.Equal(t, "145604171808000", cdd.CoinSecondsEarned)
	assert.Equal(t, UInt32(31536000), cdd.VestingSeconds)
	assert.Equal(t, int64(0), cdd.StartClaim.Unix())
	assert.Equal(t, int64(1539960000), cdd.CoinSecondsEarnedLastUpdate.Unix())
}

func TestVBLoadLinear(t *testing.T) {
	var vb VestingBalance
	loadTestJsonObject(t, "testdata/vb_linear.json", &vb)

	assert.Equal(t, AssetAmount{CoreAssetID, Int64(120000000000)}, vb.Balance)

	assert.Equal(t, int8(VestingPolicyTypeLinear), vb.Policy.Type)
	assert.IsType(t, &LinearVestingPolicy{}, vb.Policy.Value)

	p := vb.Policy.Value.(*LinearVestingPolicy)
	assert.Equal(t, int64(1540209600), p.BeginTimestamp.Unix())
	assert.Equal(t, Int64(120000000000), p.BeginBalance)

	// Test AllowedToWithdraw()

	assert.Equal(t, int64(0), p.AllowedToWithdraw(p.BeginTimestamp.Time, int64(vb.Balance.Amount)))

	beforeCliff := p.BeginTimestamp.Add(time.Duration(p.VestingCliffSeconds-1) * time.Second)
	assert.Equal(t, int64(0), p.AllowedToWithdraw(beforeCliff.Time, int64(vb.Balance.Amount)))

	endTime := p.BeginTimestamp.Add(time.Duration(p.VestingDurationSeconds) * time.Second)
	assert.Equal(t, int64(vb.Balance.Amount), p.AllowedToWithdraw(endTime.Time, int64(vb.Balance.Amount)))

	halfTime := endTime.Add(-time.Duration(p.VestingDurationSeconds/2) * time.Second)
	assert.Equal(t, int64(vb.Balance.Amount/2), p.AllowedToWithdraw(halfTime.Time, int64(vb.Balance.Amount)))
	assert.Equal(t, int64(0), p.AllowedToWithdraw(halfTime.Time, int64(vb.Balance.Amount/2)))
	assert.Equal(t, int64(p.BeginBalance/4), p.AllowedToWithdraw(halfTime.Time, int64(vb.Balance.Amount/4*3)))
}
