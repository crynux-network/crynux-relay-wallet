package relay_api

import (
	"crynux_relay_wallet/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type TaskFeeLog struct {
	ID        uint   `json:"id"`
	CreatedAt uint64 `json:"created_at"`
	Address   string `json:"address"`
	TaskFee   string `json:"task_fee"`
}

type GetTaskFeeLogsInput struct {
	StartID uint `query:"start_id" json:"start_id" description:"Start ID"`
	Limit   int  `query:"limit" json:"limit" description:"Limit"`
}

type GetTaskFeeLogsOutput struct {
	Data []TaskFeeLog `json:"data"`
}

func GetTaskFeeLogs(ctx context.Context, pivotTaskFeeLogID uint, limit int) ([]TaskFeeLog, error) {

	conf := config.GetConfig()

	req, err := http.NewRequest("GET", conf.Relay.Api.Host+"/v1/task_fee/logs", nil)
	if err != nil {
		return nil, err
	}

	input := GetTaskFeeLogsInput{
		StartID: pivotTaskFeeLogID,
		Limit:   limit,
	}

	timestamp, signature, err := SignData(input, conf.Relay.Api.PrivateKey)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("start_id", fmt.Sprintf("%d", pivotTaskFeeLogID))
	q.Add("limit", fmt.Sprintf("%d", limit))
	q.Add("timestamp", fmt.Sprintf("%d", timestamp))
	q.Add("signature", signature)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorln(err)
		}
	}(resp.Body)

	if err := processRelayResponse(resp); err != nil {
		return nil, err
	}

	var output GetTaskFeeLogsOutput

	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, err
	}

	return output.Data, nil
}
