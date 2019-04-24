package objects

import (
	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

const (
	DefaultMaintenanceCollateralRatio = 1750
	DefaultMaximumShortSqueezeRatio   = 1500
)

type PriceFeed struct {
	SettlementPrice            Price  `json:"settlement_price"`
	MaintenanceCollateralRatio UInt16 `json:"maintenance_collateral_ratio"`
	MaximumShortSqueezeRatio   UInt16 `json:"maximum_short_squeeze_ratio"`
	CoreExchangeRate           Price  `json:"core_exchange_rate"`
}

//implements Operation interface
func (p PriceFeed) Marshal(enc *util.TypeEncoder) error {

	if err := enc.Encode(p.SettlementPrice); err != nil {
		return errors.Annotate(err, "encode settlement_price")
	}

	if err := enc.Encode(p.MaintenanceCollateralRatio); err != nil {
		return errors.Annotate(err, "encode maintenance_collateral_ratio")
	}

	if err := enc.Encode(p.MaximumShortSqueezeRatio); err != nil {
		return errors.Annotate(err, "encode maximum_short_squeeze_ratio")
	}

	if err := enc.Encode(p.CoreExchangeRate); err != nil {
		return errors.Annotate(err, "encode core_exchange_rate")
	}

	return nil
}

type AssetPublishFeedOperation struct {
	Fee        AssetAmount `json:"fee"`
	Publisher  GrapheneID  `json:"publisher"`
	AssetID    GrapheneID  `json:"asset_id"`
	Feed       PriceFeed   `json:"feed"`
	Extensions Extensions  `json:"extensions"`
}

func NewAssetPublishFeedOperation() *AssetPublishFeedOperation {
	op := &AssetPublishFeedOperation{Extensions: Extensions{}}
	op.Feed.MaintenanceCollateralRatio = DefaultMaintenanceCollateralRatio
	op.Feed.MaximumShortSqueezeRatio = DefaultMaximumShortSqueezeRatio
	return op
}

//implements Operation interface
func (p *AssetPublishFeedOperation) ApplyFee(fee AssetAmount) {
	p.Fee = fee
}

//implements Operation interface
func (p *AssetPublishFeedOperation) Type() OperationType {
	return OperationTypeASSET_PUBLISH_FEED_OPERATION
}

//implements Operation interface
func (p *AssetPublishFeedOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Publisher); err != nil {
		return errors.Annotate(err, "encode publisher")
	}

	if err := enc.Encode(p.AssetID); err != nil {
		return errors.Annotate(err, "encode asset_id")
	}

	if err := enc.Encode(p.Feed); err != nil {
		return errors.Annotate(err, "encode feed")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}
