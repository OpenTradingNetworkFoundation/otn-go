package objects

import (
	json "encoding/json"

	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

type Authority struct {
	WeightThreshold UInt32          `json:"weight_threshold"`
	AccountAuths    MapAccountAuths `json:"account_auths"`
	KeyAuths        MapKeyAuths     `json:"key_auths"`
	Extensions      Extensions      `json:"extensions"`
}

func (p Authority) Marshal(enc *util.TypeEncoder) error {
	return util.BinaryEncodeStruct(enc, &p)
}

type MapAccountAuths map[GrapheneID]uint16

func (p MapAccountAuths) MarshalJSON() ([]byte, error) {
	// convert map to array of pairs: [ [key, value], [key, value] ]
	data := make([]interface{}, len(p))
	idx := 0
	for key, value := range p {
		data[idx] = []interface{}{key, value}
		idx++
	}
	return json.Marshal(data)
}

func (p *MapAccountAuths) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		return nil
	}

	var res interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return errors.Annotate(err, "unmarshal MapAccountAuths")
	}

	(*p) = make(map[GrapheneID]uint16)
	accAuths := res.([]interface{})

	for _, aa := range accAuths {
		tk := aa.([]interface{})
		(*p)[*NewGrapheneID(ObjectID(tk[0].(string)))] = uint16(tk[1].(float64))
	}

	return nil
}

type MapKeyAuths map[PublicKey]uint16

func (p MapKeyAuths) MarshalJSON() ([]byte, error) {
	// convert map to array of pairs: [ [key, value], [key, value] ]
	data := make([]interface{}, len(p))
	idx := 0
	for key, value := range p {
		data[idx] = []interface{}{key, value}
		idx++
	}
	return json.Marshal(data)
}

func (p *MapKeyAuths) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		return nil
	}

	var res interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return errors.Annotate(err, "unmarshal MapKeyAuths")
	}

	(*p) = make(map[PublicKey]uint16)
	keyAuths := res.([]interface{})

	for _, ka := range keyAuths {
		tk := ka.([]interface{})
		(*p)[NewPublicKey(tk[0].(string))] = uint16(tk[1].(float64))
	}

	return nil
}
