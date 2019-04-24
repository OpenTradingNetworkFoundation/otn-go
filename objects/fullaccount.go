package objects

import "encoding/json"

// struct full_account
// {
//    account_object                   account;
//    account_statistics_object        statistics;
//    string                           registrar_name;
//    string                           referrer_name;
//    string                           lifetime_referrer_name;
//    vector<variant>                  votes;
//    optional<vesting_balance_object> cashback_balance;
//    vector<account_balance_object>   balances;
//    vector<vesting_balance_object>   vesting_balances;
//    vector<limit_order_object>       limit_orders;
//    vector<call_order_object>        call_orders;
//    vector<force_settlement_object>  settle_orders;
//    vector<proposal_object>          proposals;
//    vector<asset_id_type>            assets;
//    vector<withdraw_permission_object> withdraws;
// };

type AccountBalance struct {
	Owner   GrapheneID `json:"owner"`
	AssetID GrapheneID `json:"asset_type"`
	Balance Int64      `json:"balance"`
}

type FullAccount struct {
	Account              Account `json:"account"`
	RegistrarName        string  `json:"registrar_name"`
	ReferrerName         string  `json:"referrer_name"`
	LifetimeReferrerName string  `json:"lifetime_referrer_name"`

	Balances    []AccountBalance `json:"balances"`
	LimitOrders []LimitOrder     `json:"limit_orders"`
	CallOrders  []CallOrder      `json:"call_orders"`
	Assets      []GrapheneID     `json:"assets"`
}

type FullAccountResult struct {
	NameOrID    string
	FullAccount FullAccount
}

func (f *FullAccountResult) UnmarshalJSON(data []byte) error {
	res := [2]interface{}{&f.NameOrID, &f.FullAccount}
	return json.Unmarshal(data, &res)
}

func (f *FullAccountResult) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{f.NameOrID, f.FullAccount})
}
