package objects

import (
	"encoding/json"

	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

type VestingBalanceCreateOperation struct {
	Fee     AssetAmount              `json:"fee"`
	Creator GrapheneID               `json:"creator"`
	Owner   GrapheneID               `json:"owner"`
	Amount  AssetAmount              `json:"amount"`
	Policy  VestingPolicyInitializer `json:"policy"`
}

type VestingBalanceWithdrawOperation struct {
	Fee            AssetAmount `json:"fee"`
	VestingBalance GrapheneID  `json:"vesting_balance"`
	Owner          GrapheneID  `json:"owner"`
	Amount         AssetAmount `json:"amount"`
}

type LinearVestingPolicyInitializer struct {
	BeginTimestamp         Time   `json:"begin_timestamp"`
	VestingCliffSeconds    UInt32 `json:"vesting_cliff_seconds"`
	VestingDurationSeconds UInt32 `json:"vesting_duration_seconds"`
}

func (p LinearVestingPolicyInitializer) Marshal(enc *util.TypeEncoder) error {
	return util.BinaryEncodeStruct(enc, &p)
}

type CDDVestingPolicyInitializer struct {
	StartClaim     Time   `json:"start_claim"`
	VestingSeconds UInt32 `json:"vesting_seconds"`
}

func (p CDDVestingPolicyInitializer) Marshal(enc *util.TypeEncoder) error {
	return util.BinaryEncodeStruct(enc, &p)
}

type VestingPolicyInitializer Variant

func (p *VestingPolicyInitializer) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{p.Type, p.Value})
}

func (p VestingPolicyInitializer) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.Type); err != nil {
		return errors.Annotate(err, "encode operation id")
	}
	return enc.Encode(p.Value)
}

func (o *VestingBalanceCreateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(o.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
	}

	return util.BinaryEncodeStruct(enc, o)
}

func (o *VestingBalanceCreateOperation) ApplyFee(fee AssetAmount) {
	o.Fee = fee
}

func (o *VestingBalanceCreateOperation) Type() OperationType {
	return OperationTypeVESTING_BALANCE_CREATE_OPERATION
}

func (o *VestingBalanceWithdrawOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(o.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
	}

	return util.BinaryEncodeStruct(enc, o)
}

func (o *VestingBalanceWithdrawOperation) ApplyFee(fee AssetAmount) {
	o.Fee = fee
}

func (o *VestingBalanceWithdrawOperation) Type() OperationType {
	return OperationTypeVESTING_BALANCE_WITHDRAW_OPERATION
}
