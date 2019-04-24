package objects

import (
	"encoding/json"

	"github.com/btcsuite/btcd/btcec"

	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/crypto"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

type PublicKey struct {
	key string
}

func (p *PublicKey) GetPublicKey() *btcec.PublicKey {
	return crypto.ParsePublicKey(p.key)
}

func (p *PublicKey) String() string {
	return p.key
}

func (p *PublicKey) Valid() bool {
	return len(p.key) > 0
}

func (p *PublicKey) UnmarshalJSON(s []byte) error {
	str := string(s)

	if len(str) > 0 && str != "null" {
		q, err := util.SafeUnquote(str)
		if err != nil {
			return errors.Annotate(err, "SafeUnquote")
		}

		p.key = q
		return nil
	}

	return errors.Errorf("unmarshal PublicKey from %s", str)
}

// implements TypeMarshaller interface
func (p PublicKey) Marshal(enc *util.TypeEncoder) error {
	pk := p.GetPublicKey()
	if pk == nil {
		return errors.Errorf("Failed to decode public key '%s'", p.key)
	}

	return enc.EncodeBytes(pk.SerializeCompressed())
}

func (p PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.key)
}

func NewPublicKey(key string) PublicKey {
	return PublicKey{key: key}
}
