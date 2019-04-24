package objects

type DynamicGlobalProperties struct {
	ID                             GrapheneID `json:"id"`
	CurrentWitness                 GrapheneID `json:"current_witness"`
	LastBudgetTime                 Time       `json:"last_budget_time"`
	Time                           Time       `json:"time"`
	NextMaintenanceTime            Time       `json:"next_maintenance_time"`
	AccountsRegisteredThisInterval UInt32     `json:"accounts_registered_this_interval"`
	DynamicFlags                   UInt32     `json:"dynamic_flags"`
	HeadBlockID                    BlockID    `json:"head_block_id"`
	RecentSlotsFilled              string     `json:"recent_slots_filled"`
	HeadBlockNumber                UInt32     `json:"head_block_number"`
	LastIrreversibleBlockNum       UInt32     `json:"last_irreversible_block_num"`
	CurrentAslot                   UInt64     `json:"current_aslot"`
	WitnessBudget                  UInt64     `json:"witness_budget"`
	RecentlyMissedCount            UInt64     `json:"recently_missed_count"`
}

func (p *DynamicGlobalProperties) RefBlockNum() UInt16 {
	return UInt16(p.HeadBlockNumber)
}

func (p *DynamicGlobalProperties) RefBlockPrefix() UInt32 {
	return p.HeadBlockID.RefBlockPrefix()
}
