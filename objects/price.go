package objects

import (
	"math"

	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

type Price struct {
	Base  AssetAmount `json:"base"`
	Quote AssetAmount `json:"quote"`
}

func (p Price) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.Base); err != nil {
		return errors.Annotate(err, "encode base")
	}

	if err := enc.Encode(p.Quote); err != nil {
		return errors.Annotate(err, "encode quote")
	}

	return nil
}

func (p Price) Rate(precBase, precQuote int) Rate {
	return Rate(p.Base.Rate(precBase) / p.Quote.Rate(precQuote))
}

func (p *Price) SetRate(precBase, precQuote int, price float64) {
	price = price * math.Pow10(precBase) / math.Pow10(precQuote)
	denominator := math.Pow10(precQuote)
	numerator := math.Round(price * math.Pow10(precQuote))
	p.Quote.Amount = Int64(denominator)
	p.Base.Amount = Int64(numerator)
}

func (p *Price) Set(base *Asset, quote *Asset, price float64) {
	p.Base.Asset = base.ID
	p.Quote.Asset = quote.ID
	p.SetRate(base.Precision, quote.Precision, price)
}

func (p Price) Valid() bool {
	return p.Base.Valid() && p.Quote.Valid()
}
