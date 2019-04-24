package tests

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"

	"github.com/opentradingnetworkfoundation/otn-go/api"
	"github.com/opentradingnetworkfoundation/otn-go/objects"
	"github.com/opentradingnetworkfoundation/otn-go/util"

	"github.com/stretchr/testify/suite"
)

type operationsAPITest struct {
	suite.Suite
	TestAPI api.BitsharesAPI
}

func (suite *operationsAPITest) SetupTest() {

	api := api.New(wsTestApiUrl, rpcApiUrl)

	if err := api.Connect(); err != nil {
		suite.Fail(err.Error(), "Connect")
	}

	api.OnError(func(err error) {
		suite.Fail(err.Error(), "OnError")
	})

	suite.TestAPI = api
}

func (suite *operationsAPITest) TearDown() {
	if err := suite.TestAPI.Close(); err != nil {
		suite.Fail(err.Error(), "Close")
	}
}

func (suite *operationsAPITest) Test_LimitOrderCancelOperation() {
	time.Sleep(1 * time.Second)
	var buf bytes.Buffer
	enc := util.NewTypeEncoder(&buf)

	tx := objects.NewTransaction()
	tx.RefBlockNum = 555
	tx.RefBlockPrefix = 3333333

	if err := tx.Expiration.UnmarshalJSON([]byte(`"2006-01-02T15:04:05"`)); err != nil {
		suite.Fail(err.Error(), "Unmarshal time")
	}

	op := objects.NewLimitOrderCancelOperation(
		*objects.NewGrapheneID("1.7.69314"),
	)

	op.Order = *objects.NewGrapheneID("1.7.123")
	op.FeePayingAccount = *objects.NewGrapheneID("1.2.456")
	op.Fee = objects.AssetAmount{
		Amount: 1000,
		Asset:  *objects.NewGrapheneID("1.3.789"),
	}

	tx.Operations = append(tx.Operations, op)
	if err := enc.Encode(tx); err != nil {
		suite.Fail(err.Error(), "Encode")
	}

	res := hex.EncodeToString(buf.Bytes())
	suite.Equal("2b02d5dc3200e540b9430102e8030000000000009506c8037b000000", res)
}

func TestOperations(t *testing.T) {
	testSuite := new(operationsAPITest)
	suite.Run(t, testSuite)
	testSuite.TearDown()
}
