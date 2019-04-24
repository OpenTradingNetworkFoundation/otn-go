package objects

type MarketTicker struct {
	Time          Time   `json:"time"`
	Base          string `json:"base"`
	Quote         string `json:"quote"`
	Latest        string `json:"latest"`
	LowestAsk     string `json:"lowest_ask"`
	HighestBid    string `json:"highest_bid"`
	PercentChange string `json:"percent_change"`
	BaseVolume    string `json:"base_volume"`
	QuoteVolume   string `json:"quote_volume"`
}
