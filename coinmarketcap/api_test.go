package coinmarketcap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSingleCoin(t *testing.T) {
	client := NewClient(nil)
	info, err := client.Ticker(&TickerOptions{ID: "1", Convert: "EUR"})
	assert.Nil(t, err)
	assert.NotNil(t, info)

	assert.Equal(t, 1, info.ID)
	assert.Equal(t, "BTC", info.Symbol)
	assert.Len(t, info.Quotes, 2)
}

func TestTickers(t *testing.T) {
	client := NewClient(nil)
	info, err := client.Tickers(&TickersOptions{Limit: 10})
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Len(t, info, 10)

	btc, ok := info["1"]
	require.True(t, true, ok)
	assert.Equal(t, "BTC", btc.Symbol)
}
