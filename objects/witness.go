package objects

// Witness represents witness object type
type Witness struct {
	ID                    GrapheneID  `json:"id"`
	WitnessAccount        GrapheneID  `json:"witness_account"`
	LastAslot             UInt64      `json:"last_aslot"`
	SigningKey            PublicKey   `json:"signing_key"`
	PayVb                 *GrapheneID `json:"pay_vb"`
	VoteID                string      `json:"vote_id"`
	TotalVotes            UInt64      `json:"total_votes"`
	URL                   string      `json:"url"`
	TotalMissed           Int64       `json:"total_missed"`
	LastConfirmedBlockNum UInt32      `json:"last_confirmed_block_num"`
}
