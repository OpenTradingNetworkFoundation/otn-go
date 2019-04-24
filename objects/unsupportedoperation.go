package objects

import (
	"encoding/json"

	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

type UnsupportedOperation struct {
	opType OperationType
	value  json.RawMessage
}

func (o *UnsupportedOperation) ApplyFee(fee AssetAmount) {
	// NOT SUPPORTED
}

func (o *UnsupportedOperation) Type() OperationType {
	return o.opType
}

func (o *UnsupportedOperation) MarshalJSON() ([]byte, error) {
	return o.value, nil
}

func (o *UnsupportedOperation) UnmarshalJSON(data []byte) error {
	o.value = json.RawMessage(data)
	return nil
}

func (o *UnsupportedOperation) Marshal(enc *util.TypeEncoder) error {
	return errors.NotSupportedf("binary marshalling for operation type=%d", o.opType)
}

func (o *UnsupportedOperation) RawJSON() []byte {
	return o.value
}
