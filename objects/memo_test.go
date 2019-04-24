package objects

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/opentradingnetworkfoundation/otn-go/crypto"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

func TestMemo(t *testing.T) {
	memo := &Memo{}
	text := "Hello, world!"

	fromKey, _ := crypto.GetPrivateKey("5JB331gjg4AheTKNuiUYWL37rC6Vbi943agL4o1fH4uoJmDtJiB")
	toKey, _ := crypto.GetPrivateKey("5HzbeKJPi1DJGC6ShPju1Muj2TTMRkc8Cs1DHJuyA38YjTkfNgB")

	memo.SetMessage(text, "OTN", fromKey, toKey.PubKey(), 100)
	messsage, err := memo.GetMessage(toKey)
	assert.NoError(t, err, "GetMessage")
	assert.Equal(t, text, messsage)
	assert.Equal(t, "cc3c71df79507f62000b8ce3999bbf8a40d772e10f6982b45ee5a581fde1d094", hex.EncodeToString(memo.Message))

	const encoded = "02052c560f3e907f927734e32b7bbd61b07be11085755456003581aab1afcec62b02e3dec77ec0dbbcf43795475fb14c3cd2be8cd4377fa279aeb79e5e6e6f267f4c640000000000000020cc3c71df79507f62000b8ce3999bbf8a40d772e10f6982b45ee5a581fde1d094"

	bin, err := util.EncodeToBytes(memo)
	assert.NoError(t, err, "EncodeToBytes")
	assert.Equal(t, encoded, hex.EncodeToString(bin))

	// check that memo is serialized without failure
	js, err := json.Marshal(memo)
	assert.NoError(t, err)
	t.Log(string(js))

	memo2 := &Memo{}
	err = json.Unmarshal(js, memo2)
	assert.NoError(t, err, "Unmarshal")
	bin2, err := util.EncodeToBytes(memo2)
	assert.Equal(t, bin, bin2)

	// broken message should be detected using checksum
	memo.Message[10] = 'c'
	messsage, err = memo.GetMessage(toKey)
	assert.Error(t, err)
}

func TestPlaintextMemo(t *testing.T) {
	memo := &Memo{}
	text := "Hello, world!"
	toKey, _ := crypto.GetPrivateKey("5HzbeKJPi1DJGC6ShPju1Muj2TTMRkc8Cs1DHJuyA38YjTkfNgB")

	// set message without key
	memo.SetMessage(text, "OTN", nil, nil, 0)
	// try to decode it
	msg, err := memo.GetMessage(toKey)
	require.NoError(t, err)
	require.Equal(t, text, msg)
}
