package alert

import (
	"crynux_relay_wallet/config"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func SendAlert(taskName, alertMessage string) error {
	content := fmt.Sprintf("Crynux Relay Wallet - task: %s, alert: %s", taskName, alertMessage)

	appConfig := config.GetConfig()
	alertLogFile := appConfig.Log.AlertLogFile

	f, err := os.OpenFile(alertLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(content + "\n"); err != nil {
		return err
	}
	return nil
}

func SafeSendAlert(taskName, alertMessage string) {
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Errorf("[Task:%s] Alert send failed due to panic: %v", taskName, recovered)
		}
	}()

	if err := SendAlert(taskName, alertMessage); err != nil {
		log.Errorf("[Task:%s] Alert send failed: %v", taskName, err)
	}
}

func SendHeartbeat() error {
	appConfig := config.GetConfig()
	logFile := appConfig.Log.HeartbeatLogFile

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString("Crynux Relay Wallet - heartbeat\n"); err != nil {
		return err
	}
	return nil
}
