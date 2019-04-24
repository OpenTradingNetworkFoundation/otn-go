package objects

import (
	"math"

	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

type AssetAmount struct {
	Asset  GrapheneID `json:"asset_id"`
	Amount Int64      `json:"amount"`
}

func AmountToFloat(amount Int64, prec int) float64 {
	return float64(amount) / math.Pow(10, float64(prec))
}

func (p *AssetAmount) Rate(prec int) float64 {
	return AmountToFloat(p.Amount, prec)
}

func (p AssetAmount) Valid() bool {
	return p.Asset.Valid() && p.Amount != 0
}

//implements Operation interface
func (p AssetAmount) Marshal(enc *util.TypeEncoder) error {

	if err := enc.Encode(p.Amount); err != nil {
		return errors.Annotate(err, "encode amount")
	}

	if err := enc.Encode(p.Asset); err != nil {
		return errors.Annotate(err, "encode asset")
	}

	return nil
}
