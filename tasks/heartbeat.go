package tasks

import (
	"context"
	"time"

	"crynux_relay_wallet/alert"

	log "github.com/sirupsen/logrus"
)

func StartHeartbeat(ctx context.Context) {
	log.Infoln("[Heartbeat] Starting heartbeat...")
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Infoln("[Heartbeat] Heartbeat stopped")
			return
		case <-ticker.C:
			alert.SendHeartbeat()
		}
	}
}
