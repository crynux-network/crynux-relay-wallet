package relay_api

import (
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}
