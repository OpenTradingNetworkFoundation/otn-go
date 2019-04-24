package objects

import (
	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

type LimitOrderCancelOperation struct {
	FeePayingAccount GrapheneID  `json:"fee_paying_account"`
	Order            GrapheneID  `json:"order"`
	Fee              AssetAmount `json:"fee"`
	Extensions       Extensions  `json:"extensions"`
}

//implements Operation interface
func (p *LimitOrderCancelOperation) ApplyFee(fee AssetAmount) {
	p.Fee = fee
}

//implements Operation interface
func (p LimitOrderCancelOperation) Type() OperationType {
	return OperationTypeLimitOrderCancel
}

//implements Operation interface
func (p LimitOrderCancelOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.FeePayingAccount); err != nil {
		return errors.Annotate(err, "encode from")
	}

	if err := enc.Encode(p.Order); err != nil {
		return errors.Annotate(err, "encode to")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}
	return nil
}

func NewLimitOrderCancelOperation(order, account GrapheneID) *LimitOrderCancelOperation {
	op := LimitOrderCancelOperation{
		FeePayingAccount: account,
		Order:            order,
		Extensions:       Extensions{},
	}

	return &op
}
