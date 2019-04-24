package objects

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const propsJSON = `{"id":"2.1.0","head_block_number":124814,"head_block_id":"000012cec0d640c21813cd6995b72b4372200740","time":"2018-09-17T15:24:15","current_witness":"1.6.5","next_maintenance_time":"2018-09-17T15:30:00","last_budget_time":"2018-09-17T15:20:00","witness_budget":1440000,"accounts_registered_this_interval":0,"recently_missed_count":0,"current_aslot":1746165,"recent_slots_filled":"340282366920938463463374607431768211455","dynamic_flags":0,"last_irreversible_block_num":4805}`

func TestDynamicProperties(t *testing.T) {
	var props DynamicGlobalProperties
	err := json.Unmarshal([]byte(propsJSON), &props)
	require.NoError(t, err)

	assert.Equal(t, UInt16(124814%65536), props.RefBlockNum())
	prefix := props.RefBlockPrefix()
	assert.Equal(t, UInt32(0xc240d6c0), prefix)
}
