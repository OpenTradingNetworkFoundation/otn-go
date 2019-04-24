package amount

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/juju/errors"
)

type Amount struct {
	amount    *big.Int
	precision int
}

func pow10(exp int) *big.Int {
	if exp == 0 {
		return big.NewInt(1)
	}

	return new(big.Int).Exp(
		big.NewInt(10),
		big.NewInt(int64(exp)), nil)
}

func (a *Amount) Valid() bool {
	return a.amount != nil
}

func (a *Amount) Value() *big.Int {
	return a.amount
}

func (a *Amount) Precision() int {
	return a.precision
}

func (a *Amount) Float() *big.Float {
	v := new(big.Float).SetInt(a.amount)
	return v.Quo(v, new(big.Float).SetInt(pow10(a.precision)))
}

func (a *Amount) Float64() float64 {
	v, _ := a.Float().Float64()
	return v
}

func (a *Amount) Uint64() uint64 {
	return a.amount.Uint64()
}

func New(precision int) Amount {
	return Amount{
		amount:    new(big.Int),
		precision: precision,
	}
}

func NewValue(precision int, value uint64) Amount {
	return Amount{
		amount:    new(big.Int).SetUint64(value),
		precision: precision,
	}
}

func FromString(precision int, value string) (Amount, error) {
	amount := New(precision)
	err := amount.Parse(value)
	return amount, err
}

func FromDecimalString(value string) (Amount, error) {
	var integer string
	var prec int

	parts := strings.Split(value, ".")
	if len(parts) == 1 {
		integer = value
	} else if len(parts) == 2 {
		// strip the insignificant digits for more accurate comparisons.
		frac := strings.TrimRight(parts[1], "0")
		prec = len(frac)
		integer = parts[0] + frac
	} else {
		return Amount{}, fmt.Errorf("can't parse '%s' as decimal", value)
	}

	intValue := new(big.Int)
	intValue, ok := intValue.SetString(integer, 10)
	if !ok {
		return Amount{}, fmt.Errorf("can't parse '%s' as decimal", value)
	}

	return Amount{precision: prec, amount: intValue}, nil
}

func (a *Amount) Parse(value string) error {
	intValue := new(big.Int)
	_, err := fmt.Sscan(value, intValue)
	if err != nil {
		return errors.Annotate(err, "parse amount")
	}

	a.amount = intValue
	return nil
}

func (a *Amount) String() string {
	return a.amount.String()
}

func (a *Amount) Scaled(targetPrecision int) Amount {
	if targetPrecision == a.precision {
		return Amount{amount: a.amount, precision: a.precision}
	}

	v := new(big.Int)

	if targetPrecision > a.precision {
		v = v.Mul(a.amount, pow10(targetPrecision-a.precision))
	} else {
		v = v.Div(a.amount, pow10(a.precision-targetPrecision))
	}

	return Amount{amount: v, precision: targetPrecision}
}

func (a *Amount) Add(x, y Amount) {
	result := New(a.precision)
	result.amount.Add(
		x.Scaled(a.precision).amount,
		y.Scaled(a.precision).amount)
	*a = result
}
