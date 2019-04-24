package objects

import (
	"time"
)

type VestingBalance struct {
	ID      GrapheneID    `json:"id"`
	Owner   GrapheneID    `json:"owner"`
	Balance AssetAmount   `json:"balance"`
	Policy  VestingPolicy `json:"policy"`
}

const (
	VestingPolicyTypeLinear = 0
	VestingPolicyTypeCDD    = 1
)

type VestingPolicy Variant

func (p *VestingPolicy) UnmarshalJSON(data []byte) error {
	return unpackTypedObject(data, &p.Type, &p.Value, func(objType int8) interface{} {
		switch objType {
		case VestingPolicyTypeLinear:
			return &LinearVestingPolicy{}
		case VestingPolicyTypeCDD:
			return &CDDVestingPolicy{}
		default:
			return nil
		}
	})
}

/*
   /// This is the time at which funds begin vesting.
   fc::time_point_sec begin_timestamp;
   /// No amount may be withdrawn before this many seconds of the vesting period have elapsed.
   uint32_t vesting_cliff_seconds = 0;
   /// Duration of the vesting period, in seconds. Must be greater than 0 and greater than vesting_cliff_seconds.
   uint32_t vesting_duration_seconds = 0;
   /// The total amount of asset to vest.
   share_type begin_balance;
*/

type LinearVestingPolicy struct {
	BeginTimestamp         Time   `json:"begin_timestamp"`
	VestingCliffSeconds    UInt32 `json:"vesting_cliff_seconds"`
	VestingDurationSeconds UInt32 `json:"vesting_duration_seconds"`
	BeginBalance           Int64  `json:"begin_balance"`
}

func (p *LinearVestingPolicy) AllowedToWithdraw(now time.Time, balance int64) int64 {
	elapsedSeconds := now.Unix() - p.BeginTimestamp.Unix()
	if elapsedSeconds <= int64(p.VestingCliffSeconds) {
		return 0
	}

	var totalVested int64
	if elapsedSeconds < int64(p.VestingDurationSeconds) {
		totalVested = int64(p.BeginBalance) * elapsedSeconds / int64(p.VestingDurationSeconds)
	} else {
		totalVested = int64(p.BeginBalance)
	}

	withdrawn := int64(p.BeginBalance) - balance

	return totalVested - withdrawn
}

type CDDVestingPolicy struct {
	VestingSeconds              UInt32 `json:"vesting_seconds"`
	CoinSecondsEarned           string `json:"coin_seconds_earned"` // uint128
	StartClaim                  Time   `json:"start_claim"`
	CoinSecondsEarnedLastUpdate Time   `json:"coin_seconds_earned_last_update"`
}
