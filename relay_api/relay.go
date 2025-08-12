package relay_api

func GetTaskFeeLogs(pivotTaskFeeLogID uint, limit int) ([]TaskFeeLog, error) {
	return nil, nil
}

func GetWithdrawalRequests(pivotWithdrawalRequestID uint, limit int) ([]WithdrawalRequest, error) {
	return nil, nil
}

func FullfillWithdrawalRequest(withdrawalRequest WithdrawalRequest, txHash string) (error) {
	return nil
}

func RejectWithdrawalRequest(withdrawalRequest WithdrawalRequest, errorMessage string) (error) {
	return nil
}
