package models

type System struct {
	ID                        uint   `json:"id" gorm:"primarykey"`
	LatestTaskFeeLogID        uint   `json:"latest_task_fee_log_id"`
	LatestTaskFeeLogTimestamp uint64 `json:"latest_task_fee_log_timestamp"`
}
