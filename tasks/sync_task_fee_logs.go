package tasks

import (
	"context"
	"crynux_relay_wallet/config"
	"crynux_relay_wallet/models"
	"crynux_relay_wallet/relay_api"
	"errors"
	"math/big"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func StartSyncTaskFeeLogs(ctx context.Context) error {

	interval := time.Duration(config.GetConfig().Tasks.SyncTaskFeeLogs.IntervalSeconds) * time.Second

	for {
		select {
		case <-ctx.Done():
			log.Infoln("Sync task fee logs task is stopping")
			return nil
		case <-time.After(interval):
			if err := syncTaskFeeLogs(); err != nil {
				log.Errorln(err)
				return err
			}
		}
	}
}

func syncTaskFeeLogs() error {
	db := config.GetDB()

	var system models.System
	err := db.First(&system).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	logs, err := relay_api.GetTaskFeeLogs(system.LatestTaskFeeLogID, int(config.GetConfig().Tasks.SyncTaskFeeLogs.BatchSize))
	if err != nil {
		return err
	}

	if len(logs) == 0 {
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, taskFeeLog := range logs {
			var account models.RelayAccount

			err := tx.Where(models.RelayAccount{Address: taskFeeLog.Address}).First(&account).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// create new account
					amount, success := new(big.Int).SetString(taskFeeLog.Amount, 10)
					if !success {
						return errors.New("cannot parse amount from task fee log")
					}
					account = models.RelayAccount{
						Address: taskFeeLog.Address,
						Balance: models.BigInt{Int: *amount},
					}
					if err := tx.Create(&account).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				// update account
				amount, success := new(big.Int).SetString(taskFeeLog.Amount, 10)
				if !success {
					return errors.New("cannot parse amount from task fee log")
				}
				newBalance := new(big.Int).Add(&account.Balance.Int, amount)
				account.Balance = models.BigInt{Int: *newBalance}
				if err := tx.Save(&account).Error; err != nil {
					return err
				}
			}
		}

		system.LatestTaskFeeLogID = logs[len(logs)-1].ID
		system.LatestTaskFeeLogTimestamp = logs[len(logs)-1].Timestamp
		return tx.Save(&system).Error
	})
}
