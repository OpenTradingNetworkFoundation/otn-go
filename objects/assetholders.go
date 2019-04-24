package objects

type AssetHolder struct {
	Name      string     `json:"name"`
	AccountID GrapheneID `json:"account_id"`
	Amount    Int64      `json:"amount"`
}
