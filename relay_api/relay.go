package relay_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}

type RelayError struct {
	StatusCode   int
	Method       string
	URL          string
	ErrorMessage string
}

func (e RelayError) Error() string {
	return fmt.Sprintf("RelayError: %s %s error code %d, %s", e.Method, e.URL, e.StatusCode, e.ErrorMessage)
}

func processRelayResponse(resp *http.Response) error {
	method := resp.Request.Method
	url := resp.Request.URL.RequestURI()
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	} else if resp.StatusCode == 400 {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		content := make(map[string]interface{})
		if err := json.Unmarshal(respBytes, &content); err != nil {
			return err
		}
		if data, ok := content["data"]; ok {
			msgBytes, err := json.Marshal(data)
			if err != nil {
				return err
			}
			msg := string(msgBytes)
			return RelayError{
				StatusCode:   resp.StatusCode,
				Method:       method,
				URL:          url,
				ErrorMessage: msg,
			}
		}
		if message, ok := content["message"]; ok {
			if msg, ok1 := message.(string); ok1 {
				return RelayError{
					StatusCode:   resp.StatusCode,
					Method:       method,
					URL:          url,
					ErrorMessage: msg,
				}
			}
		}
		return RelayError{
			StatusCode:   resp.StatusCode,
			Method:       method,
			URL:          url,
			ErrorMessage: string(respBytes),
		}
	} else {
		return RelayError{
			StatusCode:   resp.StatusCode,
			Method:       method,
			URL:          url,
			ErrorMessage: resp.Status,
		}
	}
}
