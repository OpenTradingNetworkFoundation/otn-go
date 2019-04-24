package objects

import "github.com/opentradingnetworkfoundation/otn-go/util"

type AccountOptions struct {
	MemoKey       PublicKey  `json:"memo_key"`
	VotingAccount GrapheneID `json:"voting_account"`
	NumWitness    UInt16     `json:"num_witness"`
	NumComittee   UInt16     `json:"num_comittee"`
	Votes         []Vote     `json:"votes"`
	Extensions    Extensions `json:"extensions"`
}

func (p AccountOptions) Marshal(enc *util.TypeEncoder) error {
	return util.BinaryEncodeStruct(enc, &p)
}
