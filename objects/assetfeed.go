package objects

import (
	"encoding/json"
)

type AssetFeed struct {
	ProviderID GrapheneID
	DateTime   Time
	FeedInfo   AssetFeedInfo
}

func (p *AssetFeed) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		return nil
	}

	res := [2]interface{}{
		&p.ProviderID,
		&[2]interface{}{
			&p.DateTime,
			&p.FeedInfo,
		},
	}

	return json.Unmarshal(data, &res)
}
