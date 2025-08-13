package models

type System struct {
	ID        uint64 `json:"id" gorm:"primarykey"`
	LatestTaskFeeLogID uint64 `json:"latest_task_fee_log_id"`
}
