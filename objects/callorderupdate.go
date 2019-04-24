package objects

import (
	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

// CallOrderUpdate can be used to add collateral, cover, and adjust the margin call price for a particular user
type CallOrderUpdate struct {
	Fee             AssetAmount `json:"fee"`
	FundingAccount  GrapheneID  `json:"funding_account"`
	DeltaCollateral AssetAmount `json:"delta_collateral"`
	DeltaDebt       AssetAmount `json:"delta_debt"`
	Extensions      Extensions  `json:"extensions"`
}

// ApplyFee implements Operation interface
func (p *CallOrderUpdate) ApplyFee(fee AssetAmount) {
	p.Fee = fee
}

// Type implements Operation interface
func (p CallOrderUpdate) Type() OperationType {
	return OperationTypeCallOrderUpdate
}

// Marshal implements Operation interface
func (p CallOrderUpdate) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.FundingAccount); err != nil {
		return errors.Annotate(err, "encode borrower")
	}

	if err := enc.Encode(p.DeltaCollateral); err != nil {
		return errors.Annotate(err, "encode collateral")
	}

	if err := enc.Encode(p.DeltaDebt); err != nil {
		return errors.Annotate(err, "encode debt")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

//NewCallOrderUpdate creates a new CallOrderUpdate
func NewCallOrderUpdate() *CallOrderUpdate {
	tx := CallOrderUpdate{}
	return &tx
}
