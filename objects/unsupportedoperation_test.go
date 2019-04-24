package objects_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/opentradingnetworkfoundation/otn-go/objects"
	"github.com/opentradingnetworkfoundation/otn-go/objects/testdata"
)

func TestUnknownOperations(t *testing.T) {
	var ops objects.Operations
	err := json.Unmarshal([]byte(testdata.UnknownOperationsJSON), &ops)
	require.NoError(t, err)
	require.Len(t, ops, 2)
	assert.Equal(t, objects.OperationType(100), ops[0].Type())
	assert.Equal(t, objects.OperationType(101), ops[1].Type())

	data := struct {
		Order objects.GrapheneID
	}{}

	require.IsType(t, &objects.UnsupportedOperation{}, ops[0])

	o1 := ops[0].(*objects.UnsupportedOperation)
	json.Unmarshal(o1.RawJSON(), &data)

	assert.Equal(t, objects.NewGrapheneID("1.7.1000"), &data.Order)
}
