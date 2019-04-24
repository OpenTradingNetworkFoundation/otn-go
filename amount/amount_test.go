package amount_test

import (
	"testing"

	"github.com/opentradingnetworkfoundation/otn-go/amount"
	"github.com/stretchr/testify/assert"
)

func TestAmountParse(t *testing.T) {
	a := assert.New(t)

	amount := amount.New(2)

	a.NoError(amount.Parse("1"))
	a.Equal(amount.Value().Uint64(), uint64(1))

	a.NoError(amount.Parse("100009"))
	a.Equal(amount.Value().Uint64(), uint64(100009))
}

func TestAmountParseDecimal(t *testing.T) {
	a := assert.New(t)
	v, err := amount.FromDecimalString("0")
	a.NoError(err)
	a.Equal(0, v.Precision())
	a.Equal(uint64(0), v.Uint64())

	v, err = amount.FromDecimalString("0.19")
	a.NoError(err)
	a.Equal(2, v.Precision())
	a.Equal(19, int(v.Uint64()))

	v, err = amount.FromDecimalString("0.0210")
	a.NoError(err)
	a.Equal(3, v.Precision())
	a.Equal(21, int(v.Uint64()))

	v, err = amount.FromDecimalString("10.0210")
	a.NoError(err)
	a.Equal(3, v.Precision())
	a.Equal(10021, int(v.Uint64()))
}

func TestAmountFloat(t *testing.T) {
	a := assert.New(t)
	amount := amount.New(5)
	a.NoError(amount.Parse("30000"))

	f, _ := amount.Float().Float64()
	a.Equal(0.3, f)
}

func TestAmountScale(t *testing.T) {
	a := assert.New(t)
	amount := amount.New(5)
	a.NoError(amount.Parse("30001"))

	scaled := amount.Scaled(8)
	a.Equal(uint64(30001000), scaled.Value().Uint64())

	scaledSame := amount.Scaled(amount.Precision())
	a.Equal(amount.Uint64(), scaledSame.Uint64())
}

func TestAmountAdd(t *testing.T) {
	a := assert.New(t)

	x := amount.NewValue(2, 10)
	y := amount.NewValue(2, 2)

	r := amount.New(3)
	r.Add(x, y)

	a.Equal(uint64(120), r.Uint64())

	z := amount.NewValue(0, 3)
	r.Add(x, z)
	a.Equal(uint64(3100), r.Uint64())

	r.Add(r, y)
	a.Equal(uint64(3120), r.Uint64())

	z.Add(amount.NewValue(1, 99), amount.New(0))
	a.Equal("9", z.String())

	z.Add(amount.NewValue(0, 10), amount.NewValue(2, 201))
	a.Equal("12", z.String())
}
