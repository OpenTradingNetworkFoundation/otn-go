package objects_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opentradingnetworkfoundation/otn-go/objects"
)

func TestPriceSet(t *testing.T) {
	price := objects.Price{}
	price.SetRate(4, 6, 1.5)
	assert.Equal(t, price.Rate(4, 6), objects.Rate(1.5))
}
