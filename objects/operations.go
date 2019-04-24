package objects

import (
	"encoding/json"

	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

type Operations []Operation

//implements TypeMarshaller interface
func (p Operations) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode Operations length")
	}

	for _, op := range p {
		if err := enc.Encode(op); err != nil {
			return errors.Annotate(err, "encode Operation")
		}
	}

	return nil
}

func (p Operations) MarshalJSON() ([]byte, error) {
	env := make([]OperationEnvelope, len(p))
	for idx, op := range p {
		env[idx] = OperationEnvelope{
			Type:      op.Type(),
			Operation: op,
		}
	}

	return json.Marshal(env)
}

func (p *Operations) UnmarshalJSON(data []byte) error {
	var envs []OperationEnvelope
	if err := json.Unmarshal(data, &envs); err != nil {
		return err
	}

	ops := make(Operations, len(envs))
	for idx, env := range envs {
		ops[idx] = env.Operation
	}

	*p = ops
	return nil
}

func (p Operations) ApplyFees(fees []AssetAmount) error {
	if len(p) != len(fees) {
		return errors.New("count of fees must match count of operations")
	}

	for idx, op := range p {
		op.ApplyFee(fees[idx])
	}

	return nil
}

func (p Operations) SetFeeAsset(asset GrapheneID) {
	fee := AssetAmount{Asset: asset, Amount: 0}
	for _, op := range p {
		op.ApplyFee(fee)
	}
}

func (p Operations) Types() [][]OperationType {
	ret := make([][]OperationType, len(p))
	for idx, op := range p {
		ret[idx] = []OperationType{op.Type()}
	}

	return ret
}
