package relay_api

type TaskFeeType string

const (
	Node    TaskFeeType = "node"
	Staking TaskFeeType = "staking"
)

type TaskFeeLog struct {
	ID        uint64      `json:"id"`
	Type      TaskFeeType `json:"type"`
	Address   string      `json:"address"`
	Amount    string      `json:"amount"`
	Timestamp uint64      `json:"timestamp"`
}
