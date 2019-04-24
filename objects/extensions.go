package objects

import (
	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

type Extensions []Extension

//implements TypeMarshaller interface
func (p Extensions) Marshal(enc *util.TypeEncoder) error {
	// encode size
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return err
	}

	// write extensions
	for _, ex := range p {
		if err := enc.Encode(ex); err != nil {
			return errors.Annotate(err, "encode Extension")
		}
	}

	return nil
}

type Extension []interface{}

//implements TypeMarshaller interface
func (p Extension) Marshal(enc *util.TypeEncoder) error {
	// TODO: support extension encoding
	return nil
}
