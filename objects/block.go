package objects

type Block struct {
	Previous         BlockID       `json:"previous"`
	Timestamp        Time          `json:"timestamp"`
	Witness          GrapheneID    `json:"witness"`
	MerkleRoot       Binary        `json:"transaction_merkle_root"`
	Extensions       Extensions    `json:"extensions"`
	Transactions     []Transaction `json:"transactions"`
	WitnessSignature Binary        `json:"witness_signature"`
}
