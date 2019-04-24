package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/opentradingnetworkfoundation/otn-go/api"
	"github.com/opentradingnetworkfoundation/otn-go/objects"
)

type walletAPITest struct {
	suite.Suite
	TestAPI api.BitsharesAPI
}

func (suite *walletAPITest) SetupTest() {

	api := api.New(wsTestApiUrl, rpcApiUrl)

	if err := api.Connect(); err != nil {
		suite.Fail(err.Error(), "Connect")
	}

	api.OnError(func(err error) {
		suite.Fail(err.Error(), "OnError")
	})

	suite.TestAPI = api
}

func (suite *walletAPITest) TearDown() {
	if err := suite.TestAPI.Close(); err != nil {
		suite.Fail(err.Error(), "Close")
	}
}

// func (suite *walletAPITest) Test_ListAssets() {
// 	res, err := suite.TestAPI.ListAssets("PEG.FAKEUSD", 2)
// 	if err != nil {
// 		suite.Fail(err.Error(), "ListAssets")
// 	}

// 	suite.NotNil(res)
// 	suite.Len(res, 2)
// 	util.Dump("assets >", res)
// }

/* func (suite *walletAPITest) Test_GetBlock() {
	res, err := suite.TestAPI.GetBlock(10454132)
	if err != nil {
		suite.Fail(err.Error(), "GetBlock")
	}

	suite.NotNil(res)
	util.Dump("get_block >", res)
} */

func (suite *walletAPITest) Test_ChainConfig() {
	res, err := suite.TestAPI.GetChainID()
	if err != nil {
		suite.Fail(err.Error(), "GetChainID")
	}

	suite.Equal(ChainIDBitSharesTest, res)
}

/*
func (suite *walletAPITest) Test_Buy() {

	res, err := suite.TestAPI.Buy(AccountBuySell, AssetUSD, AssetBTS, 1111, 15, true)
	if err != nil {
		suite.Fail(err.Error(), "Buy")
	}

	util.Dump("buy <", res)
	suite.NotNil(res)
}
*/
/*
func (suite *walletAPITest) Test_GetAccountByName() {

	res, err := suite.TestAPI.GetAccountByName("denk-haus")
	if err != nil {
		suite.Fail(err.Error(), "GetAccountByName")
	}

	suite.NotNil(res)
	util.Dump("accounts >", res)
} */

/* func (suite *walletAPITest) Test_GetLimitOrders() {

	res, err := suite.TestAPI.GetLimitOrders(AssetTEST, AssetPEGFAKEUSD, 50)
	if err != nil {
		suite.Fail(err.Error(), "GetLimitOrders")
	}

	suite.NotNil(res)
	util.Dump("limitorders >", res)
}
*/

func (suite *walletAPITest) Test_CancelOrder() {

	op := objects.NewLimitOrderCancelOperation(
		*objects.NewGrapheneID("1.7.69314"),
	)
	op.FeePayingAccount = *TestAccount1ID

	_, err := suite.TestAPI.Broadcast([]string{TestAccount1PrivKeyActive}, AssetTEST, op)
	if err != nil {
		suite.Fail(err.Error(), "broadcast")
	}

}

/* func (suite *walletAPITest) Test_Transfer() {

	am := objects.AssetAmount{
		Amount: 1000,
		Asset:  *AssetTEST,
	}

	op := operations.TransferOperation{
		Extensions: []objects.Extension{},
		Amount:     am,
	}

	op.From.FromObjectID(TestAccount1ID.Id())
	op.To.FromObjectID(TestAccount2ID.Id())

	priv, err := crypto.Decode(TestAccount1PrivKey)
	if err != nil {
		suite.Fail(err.Error(), "decode wif key")
	}

	privKeys := [][]byte{priv}
	if err := suite.TestAPI.Broadcast(privKeys, op.Amount.Asset, &op); err != nil {
		suite.Fail(err.Error(), "broadcast")
	}

} */

func TestWalletApi(t *testing.T) {
	testSuite := new(walletAPITest)
	suite.Run(t, testSuite)
	testSuite.TearDown()
}
