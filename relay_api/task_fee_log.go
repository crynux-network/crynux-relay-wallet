package relay_api

type TaskFeeLog struct {
	ID        uint64
	Type      string
	Address   string
	Amount    string
	Timestamp uint64
}
