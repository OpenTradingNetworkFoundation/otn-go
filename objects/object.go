package objects

import (
	"encoding/json"

	"github.com/juju/errors"
	"github.com/tidwall/gjson"
)

func UnmarshalObject(js []byte) (ret interface{}, err error) {
	id := GrapheneID{}

	idField := gjson.Get(string(js), "id")
	if !idField.Exists() {
		return nil, errors.New("no 'id' field")
	}

	if err = id.FromString(idField.String()); err != nil {
		return
	}

	switch id.Space() {
	case SpaceTypeProtocol:
		switch id.Type() {
		case ObjectTypeAsset:
			ass := Asset{}
			if err := json.Unmarshal(js, &ass); err != nil {
				return nil, errors.Annotate(err, "unmarshal Asset")
			}
			ret = ass

		case ObjectTypeAccount:
			acc := Account{}
			if err := json.Unmarshal(js, &acc); err != nil {
				return nil, errors.Annotate(err, "unmarshal Account")
			}
			ret = acc

		case ObjectTypeForceSettlement:
			set := SettleOrder{}
			if err := json.Unmarshal(js, &set); err != nil {
				return nil, errors.Annotate(err, "unmarshal SettleOrder")
			}
			ret = set

		case ObjectTypeLimitOrder:
			lim := LimitOrder{}
			if err := json.Unmarshal(js, &lim); err != nil {
				return nil, errors.Annotate(err, "unmarshal LimitOrder")
			}
			ret = lim

		case ObjectTypeCallOrder:
			cal := CallOrder{}
			if err := json.Unmarshal(js, &cal); err != nil {
				return nil, errors.Annotate(err, "unmarshal CallOrder")
			}
			ret = cal

		case ObjectTypeWitness:
			wit := Witness{}
			if err := json.Unmarshal(js, &wit); err != nil {
				return nil, errors.Annotate(err, "unmarshal Witness")
			}
			ret = wit

		case ObjectTypeVestingBalance:
			vb := VestingBalance{}
			if err := json.Unmarshal(js, &vb); err != nil {
				return nil, errors.Annotate(err, "unmarshal VestingBalance")
			}
			ret = vb

		default:
			return nil, errors.Errorf("unable to parse GrapheneObject with ID %s", id)
		}
	case SpaceTypeImplementation:
		switch id.Type() {
		case ObjectTypeAssetBitAssetData:
			bit := BitAssetData{}
			if err := json.Unmarshal(js, &bit); err != nil {
				return nil, errors.Annotate(err, "unmarshal BitAssetData")
			}
			ret = bit

		default:
			return nil, errors.Errorf("unable to parse GrapheneObject with ID %s", id)
		}
	}

	return
}
