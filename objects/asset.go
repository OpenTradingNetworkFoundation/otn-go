package objects

import (
	"math"

	"github.com/juju/errors"
)

type Asset struct {
	ID                 GrapheneID   `json:"id"`
	Symbol             string       `json:"symbol"`
	Precision          int          `json:"precision"`
	Issuer             GrapheneID   `json:"issuer"`
	DynamicAssetDataID GrapheneID   `json:"dynamic_asset_data_id"`
	BitassetDataID     GrapheneID   `json:"bitasset_data_id"`
	Options            AssetOptions `json:"options"`
}

// CoreAssetID is a shortcut for id "1.3.0" which is the id for the core asset
var CoreAssetID = *NewGrapheneID("1.3.0")

//NewAsset creates a new Asset object
func NewAsset(id ObjectID) *Asset {
	ass := Asset{}
	if err := ass.ID.FromString(string(id)); err != nil {
		panic(errors.Annotate(err, "init GrapheneID"))
	}

	return &ass
}

func (a *Asset) CreateAmount(volume float64) AssetAmount {
	amount := volume * math.Pow10(a.Precision)
	return AssetAmount{
		Asset:  a.ID,
		Amount: Int64(amount),
	}
}

func (a *Asset) GetRate(amount AssetAmount) float64 {
	return amount.Rate(a.Precision)
}
