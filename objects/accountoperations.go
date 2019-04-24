package objects

import (
	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

type AccountCreateOperation struct {
	Fee             AssetAmount    `json:"fee"`
	Registrar       GrapheneID     `json:"registrar"`
	Referrer        GrapheneID     `json:"referrer"`
	ReferrerPercent UInt16         `json:"referrer_percent"`
	Name            string         `json:"name"`
	Owner           Authority      `json:"owner"`
	Active          Authority      `json:"active"`
	Options         AccountOptions `json:"options"`
	Extensions      Extensions     `json:"extensions"`
}

//implements Operation interface
func (o *AccountCreateOperation) ApplyFee(fee AssetAmount) {
	o.Fee = fee
}

//implements Operation interface
func (o *AccountCreateOperation) Type() OperationType {
	return OperationTypeACCOUNT_CREATE_OPERATION
}

//implements Operation interface
func (o *AccountCreateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(o.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
	}

	return util.BinaryEncodeStruct(enc, o)
}
