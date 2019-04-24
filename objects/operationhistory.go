package objects

import (
	"encoding/json"

	"github.com/juju/errors"
)

// OperationResult types
const (
	OperationResultType_Void     int8 = 0
	OperationResultType_ObjectID int8 = 1
	OperationResultType_Asset    int8 = 2
)

type OperationResult struct {
	Type   int8
	Result interface{}
}

type OperationResultVoid struct{}

type OperationHistory struct {
	ID         GrapheneID        `json:"id"`
	Op         OperationEnvelope `json:"op"`
	Result     OperationResult   `json:"result"`
	BlockNum   uint32            `json:"block_num"`
	TrxInBlock uint16            `json:"trx_in_block"`
	OpInTrx    uint16            `json:"op_in_trx"`
	VirtualOp  uint16            `json:"virtual_op"`
}

type TypeResolver func(objType int8) interface{}

func unpackTypedObject(data []byte, objType *int8, objValue *interface{}, resolver TypeResolver) error {
	raw := make([]json.RawMessage, 2)
	if err := json.Unmarshal(data, &raw); err != nil {
		return errors.Annotate(err, "Unmarshal raw object")
	}

	if len(raw) != 2 {
		return errors.Errorf("Invalid operation data: %v", string(data))
	}

	if err := json.Unmarshal(raw[0], objType); err != nil {
		return errors.Annotate(err, "Unmarshal OperationType")
	}

	value := resolver(*objType)
	err := json.Unmarshal(raw[1], &value)

	if err != nil {
		return err
	}

	*objValue = value
	return nil
}

func resolveOperationResult(objType int8) interface{} {
	switch objType {
	case OperationResultType_Void:
		return &OperationResultVoid{}
	case OperationResultType_ObjectID:
		return &GrapheneID{}
	case OperationResultType_Asset:
		return &AssetAmount{}
	}

	return nil
}

func (p *OperationResult) UnmarshalJSON(data []byte) error {
	return unpackTypedObject(data, &p.Type, &p.Result, resolveOperationResult)
}

func (p *OperationResult) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{p.Type, p.Result})
}
