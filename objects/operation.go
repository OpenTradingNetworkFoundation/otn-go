package objects

import (
	"encoding/json"

	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

type Operation interface {
	util.TypeMarshaller
	ApplyFee(fee AssetAmount)
	Type() OperationType
}

type OperationEnvelope struct {
	Type      OperationType
	Operation Operation
}

func (p OperationEnvelope) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		p.Type,
		p.Operation,
	})
}

func (p *OperationEnvelope) UnmarshalJSON(data []byte) error {
	raw := make([]json.RawMessage, 2)
	if err := json.Unmarshal(data, &raw); err != nil {
		return errors.Annotate(err, "Unmarshal raw object")
	}

	if len(raw) != 2 {
		return errors.Errorf("Invalid operation data: %v", string(data))
	}

	if err := json.Unmarshal(raw[0], &p.Type); err != nil {
		return errors.Annotate(err, "Unmarshal OperationType")
	}

	p.Operation = ResolveOperationType(p.Type)
	if err := json.Unmarshal(raw[1], &p.Operation); err != nil {
		return errors.Annotate(err, "Unmarshal Operation")
	}
	return nil
}

func ResolveOperationType(op OperationType) Operation {
	switch op {
	case OperationTypeLimitOrderCreate:
		return &LimitOrderCreateOperation{}
	case OperationTypeTransfer:
		return &TransferOperation{}
	case OperationTypeCallOrderUpdate:
		return &CallOrderUpdate{}
	case OperationTypeLimitOrderCancel:
		return &LimitOrderCancelOperation{}
	case OperationTypeASSET_PUBLISH_FEED_OPERATION:
		return &AssetPublishFeedOperation{}
	case OperationTypeACCOUNT_CREATE_OPERATION:
		return &AccountCreateOperation{}
	case OperationTypeVESTING_BALANCE_CREATE_OPERATION:
		return &VestingBalanceCreateOperation{}
	case OperationTypeVESTING_BALANCE_WITHDRAW_OPERATION:
		return &VestingBalanceWithdrawOperation{}
	default:
		// fallback to special struct for unsupported operation types
		return &UnsupportedOperation{opType: op}
	}
}
