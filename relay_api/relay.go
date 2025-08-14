package relay_api

import (
	"crynux_relay_wallet/config"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func GetTaskFeeLogs(pivotTaskFeeLogID uint64, limit int) ([]TaskFeeLog, error) {

	conf := config.GetConfig()
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", conf.Relay.Api.Host+"/v1/task_fee_logs", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("pivot", strconv.FormatUint(pivotTaskFeeLogID, 10))
	q.Add("limit", strconv.Itoa(limit))
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", "Bearer "+conf.Relay.Api.Key)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorln(err)
		}
	}(resp.Body)

	var logs []TaskFeeLog

	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		return nil, err
	}

	return logs, nil
}

func GetWithdrawalRequests(pivotWithdrawalRequestID uint, limit int) ([]WithdrawalRequest, error) {
	return nil, nil
}

func FullfillWithdrawalRequest(withdrawalRequest WithdrawalRequest, txHash string) error {
	return nil
}

func RejectWithdrawalRequest(withdrawalRequest WithdrawalRequest, errorMessage string) error {
	return nil
}
