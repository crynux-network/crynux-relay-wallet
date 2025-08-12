package main

import (
	"context"
	"crynux_relay_wallet/config"
	"crynux_relay_wallet/migrate"
	"crynux_relay_wallet/tasks"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {

	log.Infoln("Starting Crynux Relay Wallet...")

	log.Infoln("Initializing config...")
	if err := config.InitConfig(""); err != nil {
		print("Error reading config file")
		print(err.Error())
		os.Exit(1)
	}

	conf := config.GetConfig()

	log.Infoln("Initializing log...")
	if err := config.InitLog(conf); err != nil {
		print("Error initializing log")
		print(err.Error())
		os.Exit(1)
	}

	log.Infoln("Initializing database...")
	if err := config.InitDB(conf); err != nil {
		log.Errorln(err.Error())
		os.Exit(1)
	}

	log.Infoln("Starting database migration...")
	startDBMigration()
	log.Infoln("DB migrations are done!")

	log.Infoln("Starting tasks...")
	ctx := context.Background()

	heartbeatCtx, heartbeatCancel := context.WithCancel(context.Background())
	go tasks.StartHeartbeat(heartbeatCtx)

	taskExit := make(chan struct{}, 2)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				taskExit <- struct{}{}
				return
			}
			taskExit <- struct{}{}
		}()
		tasks.StartSyncRelayAccounts(ctx)
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				taskExit <- struct{}{}
				return
			}
			taskExit <- struct{}{}
		}()
		tasks.StartProcessWithdrawalRequests(ctx)
	}()

	go func() {
		<-taskExit
		heartbeatCancel()
	}()
}

func startDBMigration() {

	migrate.InitMigration(config.GetDB())

	if err := migrate.Migrate(); err != nil {
		log.Errorln(err.Error())
		if err = migrate.Rollback(); err != nil {
			log.Errorln(err.Error())
		}
		os.Exit(1)
	}
}
