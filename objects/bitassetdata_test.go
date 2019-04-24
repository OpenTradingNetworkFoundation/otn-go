package objects_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/opentradingnetworkfoundation/otn-go/objects"
	"github.com/opentradingnetworkfoundation/otn-go/objects/testdata"
)

func TestBitassetData(t *testing.T) {
	var info objects.BitAssetData
	err := json.Unmarshal([]byte(testdata.BitAssetData), &info)
	require.NoError(t, err)
	require.Len(t, info.Feeds, 3)

	settlementPrice := objects.Price{
		Base: objects.AssetAmount{
			Amount: objects.Int64(19424),
			Asset:  *objects.NewGrapheneID("1.3.2"),
		},
		Quote: objects.AssetAmount{
			Amount: objects.Int64(100000000),
			Asset:  *objects.NewGrapheneID("1.3.0"),
		},
	}

	assert.Equal(t, *objects.NewGrapheneID("1.2.6"), info.Feeds[0].ProviderID)
	assert.EqualValues(t, settlementPrice, info.Feeds[0].FeedInfo.SettlementPrice)
	assert.Equal(t, *objects.NewGrapheneID("1.2.14"), info.Feeds[1].ProviderID)
	assert.Equal(t, *objects.NewGrapheneID("1.2.16"), info.Feeds[2].ProviderID)

	settlementPrice.Base.Amount = objects.Int64(19426)
	assert.EqualValues(t, settlementPrice, info.CurrentFeed.SettlementPrice)
}
