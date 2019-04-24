package objects

// CommitteeMember represents committee_member_object type
type CommitteeMember struct {
	ID                     GrapheneID `json:"id"`
	CommitteeMemberAccount GrapheneID `json:"committee_member_account"`
	VoteID                 string     `json:"vote_id"`
	TotalVotes             UInt64     `json:"total_votes"`
	URL                    string     `json:"url"`
}
