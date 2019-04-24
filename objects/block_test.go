package objects

import (
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const BlockJSON = `
{"previous":"000003e88d8b9842ee58b90a9470f752e2dbf821","timestamp":"2018-09-11T13:23:25","witness":"1.6.8","transaction_merkle_root":"c0b43bd1139de7bf587b5bdf27f7c49b92638a4b","extensions":[],
"witness_signature":"1f0e9d2cdeb6ea914a920e7b9e982bcf7f33fe43e4d1db1befa3a782a5007e001c3f5afdeebc0ad9a8c0a4cf33a458c624727f1937305c67123e5aba19a9b00087",
"transactions":[{"ref_block_num":1000,"ref_block_prefix":1117293453,"expiration":"2018-09-11T13:23:50","operations":[[2,{"fee":{"amount":0,"asset_id":"1.3.0"},"fee_paying_account":"1.2.19","order":"1.7.2568","extensions":[]}],[2,{"fee":{"amount":0,"asset_id":"1.3.0"},"fee_paying_account":"1.2.19","order":"1.7.2569","extensions":[]}],[2,{"fee":{"amount":0,"asset_id":"1.3.0"},"fee_paying_account":"1.2.19","order":"1.7.2570","extensions":[]}],[1,{"fee":{"amount":10000,"asset_id":"1.3.0"},"seller":"1.2.19","amount_to_sell":{"amount":"300000000000","asset_id":"1.3.0"},"min_to_receive":{"amount":"36328827974","asset_id":"1.3.8"},"expiration":"2018-09-11T13:25:21","fill_or_kill":false,"extensions":[]}],[1,{"fee":{"amount":10000,"asset_id":"1.3.0"},"seller":"1.2.19","amount_to_sell":{"amount":"300000000000","asset_id":"1.3.0"},"min_to_receive":{"amount":"36683255564","asset_id":"1.3.8"},"expiration":"2018-09-11T13:25:21","fill_or_kill":false,"extensions":[]}],[1,{"fee":{"amount":10000,"asset_id":"1.3.0"},"seller":"1.2.19","amount_to_sell":{"amount":"300000000000","asset_id":"1.3.0"},"min_to_receive":{"amount":"37037683154","asset_id":"1.3.8"},"expiration":"2018-09-11T13:25:21","fill_or_kill":false,"extensions":[]}]],"extensions":[],"signatures":["2061256e3cb9a124e625a121640bc63c392cc7462ba9ad4a16bcc0ef90afbfe5913f23fdcbdec0407a4acd88d011ef3e599d2936f95d8764076b7b78a4966aa89a"],"operation_results":[[2,{"amount":"300000000000","asset_id":"1.3.0"}],[2,{"amount":"300000000000","asset_id":"1.3.0"}],[2,{"amount":"300000000000","asset_id":"1.3.0"}],[1,"1.7.2598"],[1,"1.7.2599"],[1,"1.7.2600"]]}]}
`

func TestBlockJSONParser(t *testing.T) {
	var block Block
	err := json.Unmarshal([]byte(BlockJSON), &block)
	require.NoError(t, err)

	signature, _ := hex.DecodeString("1f0e9d2cdeb6ea914a920e7b9e982bcf7f33fe43e4d1db1befa3a782a5007e001c3f5afdeebc0ad9a8c0a4cf33a458c624727f1937305c67123e5aba19a9b00087")
	ts, _ := time.ParseInLocation("2006-01-02T15:04:05", "2018-09-11T13:23:25", time.UTC)
	t.Logf("Time: %s", ts)

	assert.Equal(t, uint32(1000), block.Previous.BlockNumber())
	assert.Equal(t, "8d8b9842ee58b90a9470f752e2dbf821", hex.EncodeToString(block.Previous.ID()))
	assert.Equal(t, ts.Unix(), block.Timestamp.Unix())
	assert.Equal(t, *NewGrapheneID("1.6.8"), block.Witness)
	assert.EqualValues(t, signature, block.WitnessSignature)
	assert.Equal(t, 1, len(block.Transactions))
	assert.Equal(t, UInt16(1000), block.Transactions[0].RefBlockNum)
	assert.Equal(t, 6, len(block.Transactions[0].Operations))
}
