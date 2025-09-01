package tasks

import (
	"context"
	"crynux_relay_wallet/config"
	"crynux_relay_wallet/models"
	"crynux_relay_wallet/relay_api"
	"errors"
	"fmt"
	"math/big"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func StartSyncTaskFeeLogs(ctx context.Context) error {

	interval := time.Duration(config.GetConfig().Tasks.SyncTaskFeeLogs.IntervalSeconds) * time.Second

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Infoln("Sync task fee logs task is stopping")
			return nil
		case <-ticker.C:
			if err := syncTaskFeeLogs(ctx); err != nil {
				log.Errorf("Failed to sync task fee logs: %v", err)
				return err
			}
		}
	}
}

func mergeTaskFeeLogs(logs []relay_api.TaskFeeLog) (map[string]*big.Int, error) {
	merged := make(map[string]*big.Int)

	for _, log := range logs {
		amount, success := new(big.Int).SetString(log.TaskFee, 10)
		if !success {
			return nil, errors.New("cannot parse amount from task fee log")
		}

		if _, ok := merged[log.Address]; !ok {
			merged[log.Address] = amount
		} else {
			merged[log.Address].Add(merged[log.Address], amount)
		}
	}

	return merged, nil
}

func processTaskFeeLogs(ctx context.Context, db *gorm.DB, logs []relay_api.TaskFeeLog) error {
	merged, err := mergeTaskFeeLogs(logs)
	if err != nil {
		return err
	}

	var addresses []string
	for address := range merged {
		addresses = append(addresses, address)
	}

	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return db.WithContext(dbCtx).Transaction(func(tx *gorm.DB) error {
		var accounts []*models.RelayAccount
		if err := tx.Model(&models.RelayAccount{}).Where("address IN (?)", addresses).Find(&accounts).Error; err != nil {
			return err
		}

		existedAddresses := make(map[string]bool)
		for _, account := range accounts {
			existedAddresses[account.Address] = true
		}

		for _, account := range accounts {
			amount, ok := merged[account.Address]
			if !ok {
				continue
			}
			account.Balance.Add(&account.Balance.Int, amount)
		}

		var newAccounts []*models.RelayAccount
		for address, amount := range merged {
			if _, ok := existedAddresses[address]; !ok {
				newAccounts = append(newAccounts, &models.RelayAccount{Address: address, Balance: models.BigInt{Int: *amount}})
			}
		}

		if err := tx.CreateInBatches(newAccounts, 100).Error; err != nil {
			return err
		}

		var cases string
		for _, account := range accounts {
			cases += fmt.Sprintf(" WHEN address = '%s' THEN '%s'", account.Address, account.Balance.String())
		}
		if err := tx.Model(&models.RelayAccount{}).Where("address IN (?)", existedAddresses).
			Update("balance", gorm.Expr("CASE"+cases+" END")).Error; err != nil {
			return err
		}

		return nil
	})
}

func syncTaskFeeLogs(ctx context.Context) error {
	db := config.GetDB()

	var checkpoint models.TaskFeeCheckpoint
	err := db.First(&checkpoint).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	for {
		logs, err := relay_api.GetTaskFeeLogs(ctx, checkpoint.LatestTaskFeeLogID, int(config.GetConfig().Tasks.SyncTaskFeeLogs.BatchSize))
		if err != nil {
			return err
		}

		if len(logs) == 0 {
			break
		}
		
		if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := processTaskFeeLogs(ctx, tx, logs); err != nil {
				return err
			}

			checkpoint.LatestTaskFeeLogID = logs[len(logs)-1].ID
			checkpoint.LatestTaskFeeLogTimestamp = logs[len(logs)-1].CreatedAt
			return tx.Save(&checkpoint).Error
		}); err != nil {
			return err
		}
	}
	return nil
}
