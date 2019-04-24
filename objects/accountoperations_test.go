package objects

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opentradingnetworkfoundation/otn-go/util"
)

func TestAccountCreate(t *testing.T) {
	pubKey := NewPublicKey("OTN7n6PuAVYiXMSrjRGzCfkSroA67PK19EodsTib3Yztb5iNkLEpd")

	auth := Authority{}
	auth.AccountAuths = MapAccountAuths{}
	auth.KeyAuths = MapKeyAuths{}
	auth.KeyAuths[pubKey] = 1

	op := AccountCreateOperation{}
	op.Active = auth
	op.Owner = auth
	op.Options.MemoKey = pubKey

	var buf bytes.Buffer
	enc := util.NewTypeEncoder(&buf)

	assert.NoError(t, op.Marshal(enc))

	js, err := json.Marshal(op)
	t.Log(string(js))
	assert.NoError(t, err)
}
