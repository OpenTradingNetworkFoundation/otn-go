package objects

import (
	"crypto/sha256"
	"time"

	"github.com/btcsuite/btcd/btcec"

	"github.com/juju/errors"
	"github.com/opentradingnetworkfoundation/otn-go/crypto"
	"github.com/opentradingnetworkfoundation/otn-go/util"
)

type Signature string
type Signatures []Signature

const (
	TxExpirationDefault = 30 * time.Second
)

type Transaction struct {
	RefBlockNum      UInt16            `json:"ref_block_num"`
	RefBlockPrefix   UInt32            `json:"ref_block_prefix"`
	Expiration       Time              `json:"expiration"`
	Operations       Operations        `json:"operations"`
	Extensions       Extensions        `json:"extensions"`
	Signatures       Signatures        `json:"signatures,omitempty"`
	OperationResults []OperationResult `json:"operation_results,omitempty"`
}

//implements TypeMarshaller interface
func (p Transaction) Marshal(enc *util.TypeEncoder) error {

	if err := enc.Encode(p.RefBlockNum); err != nil {
		return errors.Annotate(err, "encode RefBlockNum")
	}

	if err := enc.Encode(p.RefBlockPrefix); err != nil {
		return errors.Annotate(err, "encode RefBlockPrefix")
	}

	if err := enc.Encode(p.Expiration); err != nil {
		return errors.Annotate(err, "encode Expiration")
	}

	if err := enc.Encode(p.Operations); err != nil {
		return errors.Annotate(err, "encode Operations")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extension")
	}

	return nil
}

func (p *Transaction) Digest() ([]byte, error) {
	h := sha256.New()
	enc := util.NewTypeEncoder(h)
	if err := enc.Encode(p); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func (p *Transaction) SignDigest(chainID ChainID) ([]byte, error) {
	h := sha256.New()
	enc := util.NewTypeEncoder(h)

	if err := enc.Encode(chainID); err != nil {
		return nil, errors.Annotate(err, "encode chainID")
	}

	if err := enc.Encode(p); err != nil {
		return nil, errors.Annotate(err, "encode transaction")
	}

	return h.Sum(nil), nil

}

func (p *Transaction) ID() []byte {
	digest, err := p.Digest()
	if err != nil {
		return nil
	}

	return digest[0:20]
}

//Sign signes a Transaction with the given private keys
func (p *Transaction) Sign(wifKeys []string, chainID ChainID) error {
	keys := make([]*btcec.PrivateKey, len(wifKeys))

	for idx, wif := range wifKeys {

		key, err := crypto.GetPrivateKey(wif)
		if err != nil {
			return errors.Annotate(err, "GetPrivateKey")
		}

		keys[idx] = key
	}

	return p.SignWithKeys(keys, chainID)
}

func (p *Transaction) SignWithKeys(pvtKeys []*btcec.PrivateKey, chainID ChainID) error {

	digest, err := p.SignDigest(chainID)
	if err != nil {
		return errors.Annotate(err, "SignDigest")
	}

	p.Signatures = make([]Signature, len(pvtKeys))

	for idx, key := range pvtKeys {
		sig, err := crypto.SignDigest(digest, key)
		if err != nil {
			return errors.Annotate(err, "Sign")
		}
		p.Signatures[idx] = Signature(sig.ToHex())
	}
	return nil
}

//AdjustExpiration extends expiration by given duration.
func (p *Transaction) AdjustExpiration(dur time.Duration) {
	p.Expiration = p.Expiration.Add(dur)
}

//NewTransactionWithBlockData creates a new Transaction and initialises
//relevant Blockdata fields and expiration.
func NewTransactionWithBlockData(props *DynamicGlobalProperties) (*Transaction, error) {
	tx := Transaction{
		Extensions:     Extensions{},
		Signatures:     Signatures{},
		RefBlockNum:    props.RefBlockNum(),
		Expiration:     props.Time.Add(TxExpirationDefault),
		RefBlockPrefix: props.RefBlockPrefix(),
	}
	return &tx, nil
}

//NewTransaction creates a new Transaction
func NewTransaction() *Transaction {
	tx := Transaction{
		Extensions: Extensions{},
		Signatures: Signatures{},
	}
	return &tx
}
