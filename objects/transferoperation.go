package objects

import (
	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

type TransferOperation struct {
	From       GrapheneID  `json:"from"`
	To         GrapheneID  `json:"to"`
	Amount     AssetAmount `json:"amount"`
	Fee        AssetAmount `json:"fee"`
	Memo       *Memo       `json:"memo,omitempty"`
	Extensions Extensions  `json:"extensions"`
}

//implements Operation interface
func (p *TransferOperation) ApplyFee(fee AssetAmount) {
	p.Fee = fee
}

//implements Operation interface
func (p TransferOperation) Type() OperationType {
	return OperationTypeTransfer
}

//implements Operation interface
func (p TransferOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.From); err != nil {
		return errors.Annotate(err, "encode from")
	}

	if err := enc.Encode(p.To); err != nil {
		return errors.Annotate(err, "encode to")
	}

	if err := enc.Encode(p.Amount); err != nil {
		return errors.Annotate(err, "encode amount")
	}

	if p.Memo != nil {
		enc.Encode(uint8(1))
		if err := enc.Encode(p.Memo); err != nil {
			return errors.Annotate(err, "encode memo")
		}
	} else {
		enc.Encode(uint8(0))
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

//NewTransferOperation creates a new TransferOperation
func NewTransferOperation() *TransferOperation {
	tx := TransferOperation{
		Extensions: Extensions{},
	}
	return &tx
}
