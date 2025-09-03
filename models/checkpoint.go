package models

type TaskFeeCheckpoint struct {
	ID uint `json:"id" gorm:"primarykey"`
	LatestTaskFeeLogID uint `json:"latest_task_fee_log_id"`
	LatestTaskFeeLogTimestamp uint64 `json:"latest_task_fee_log_timestamp"`
}

type WithdrawalRequestCheckpoint struct {
	ID uint `json:"id" gorm:"primarykey"`
	LatestWithdrawalRequestID uint `json:"latest_withdrawal_request_id"`
	LatestWithdrawalRequestTimestamp uint64 `json:"latest_withdrawal_request_timestamp"`
}