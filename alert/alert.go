package alert

import (
	"crynux_relay_wallet/config"
	"fmt"
	"os"
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

func SendHeartbeat() error {
	appConfig := config.GetConfig()
	alertLogFile := appConfig.Log.AlertLogFile

	f, err := os.OpenFile(alertLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString("Crynux Relay Wallet - heartbeat\n"); err != nil {
		return err
	}
	return nil
}
