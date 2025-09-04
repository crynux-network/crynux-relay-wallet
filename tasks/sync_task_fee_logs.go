package tasks

import (
	"context"
	"crynux_relay_wallet/config"
	"crynux_relay_wallet/models"
	"crynux_relay_wallet/relay_api"
	"crynux_relay_wallet/utils"
	"errors"
	"fmt"
	"math/big"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TaskFeeError struct {
	Message string
}

func (e *TaskFeeError) Error() string {
	return e.Message
}

func NewTaskFeeError(message string) *TaskFeeError {
	return &TaskFeeError{Message: message}
}

func IsTaskFeeError(err error) bool {
	var taskFeeError *TaskFeeError
	return errors.As(err, &taskFeeError)
}

var ErrTaskFeeAmountTooLarge = NewTaskFeeError("task fee amount is greater than max task fee amount threshold")
var ErrTaskFeeAmountInvalid = NewTaskFeeError("cannot parse amount from task fee log")
var ErrTaskFeeAddressCountTooLarge = NewTaskFeeError("task fee logs count of a single address is greater than max address count threshold")
var ErrTaskFeeNewAddressCountTooLarge = NewTaskFeeError("task fee logs count of new addresses is greater than max new address count threshold")

func StartSyncTaskFeeLogs(ctx context.Context) error {

	intervalSeconds := config.GetConfig().Tasks.SyncTaskFeeLogs.IntervalSeconds
	interval := time.Duration(intervalSeconds) * time.Second

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Infoln("Sync task fee logs task is stopping")
			return nil
		case <-ticker.C:
			if err := syncTaskFeeLogs(ctx, intervalSeconds); err != nil {
				log.Errorf("Failed to sync task fee logs: %v", err)
				if IsTaskFeeError(err) {
					return err
				}
			}
		}
	}
}

func mergeTaskFeeLogs(logs []relay_api.TaskFeeLog) (map[string]*big.Int, error) {
	merged := make(map[string]*big.Int)

	for _, log := range logs {
		amount, success := new(big.Int).SetString(log.TaskFee, 10)
		if !success {
			return nil, ErrTaskFeeAmountInvalid
		}

		if _, ok := merged[log.Address]; !ok {
			merged[log.Address] = amount
		} else {
			merged[log.Address].Add(merged[log.Address], amount)
		}
	}

	return merged, nil
}

func checkTaskFeeLogs(ctx context.Context, db *gorm.DB, logs []relay_api.TaskFeeLog) error {
	appConfig := config.GetConfig()
	maxTaskFeeAmount := utils.EtherToWei(big.NewInt(int64(appConfig.Tasks.SyncTaskFeeLogs.MaxTaskFeeAmount)))

	addressLogCount := make(map[string]int)
	for _, log := range logs {
		amount, success := new(big.Int).SetString(log.TaskFee, 10)
		if !success {
			return ErrTaskFeeAmountInvalid
		}
		if amount.Cmp(maxTaskFeeAmount) > 0 {
			return ErrTaskFeeAmountTooLarge
		}
		addressLogCount[log.Address]++
	}

	var maxCount int
	for _, count := range addressLogCount {
		if count > maxCount {
			maxCount = count
		}
	}
	if uint(maxCount) > appConfig.Tasks.SyncTaskFeeLogs.MaxAddressLogsCountInBatch {
		return ErrTaskFeeAddressCountTooLarge
	}

	addresses := make([]string, 0, len(addressLogCount))
	for address := range addressLogCount {
		addresses = append(addresses, address)
	}

	var accounts []*models.RelayAccount
	if err := db.WithContext(ctx).Model(&models.RelayAccount{}).Where("address IN (?)", addresses).Find(&accounts).Error; err != nil {
		return err
	}

	existedAddresses := make(map[string]bool)
	for _, account := range accounts {
		existedAddresses[account.Address] = true
	}

	var newAddressCount int
	for address := range addressLogCount {
		if _, ok := existedAddresses[address]; !ok {
			newAddressCount++
		}
	}
	if uint(newAddressCount) > appConfig.Tasks.SyncTaskFeeLogs.MaxNewAddressCountInBatch {
		return ErrTaskFeeNewAddressCountTooLarge
	}

	return nil
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

		existedAddresses := make([]string, len(accounts))
		existedAddressMap := make(map[string]bool)
		for i, account := range accounts {
			existedAddresses[i] = account.Address
			existedAddressMap[account.Address] = true
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
			if _, ok := existedAddressMap[address]; !ok {
				newAccounts = append(newAccounts, &models.RelayAccount{Address: address, Balance: models.BigInt{Int: *amount}})
			}
		}

		if len(newAccounts) > 0 {
			if err := tx.CreateInBatches(newAccounts, 100).Error; err != nil {
				return err
			}
		}

		if len(accounts) > 0 {
			var cases string
			for _, account := range accounts {
				cases += fmt.Sprintf(" WHEN address = '%s' THEN '%s'", account.Address, account.Balance.String())
			}
			if err := tx.Model(&models.RelayAccount{}).Where("address IN (?)", existedAddresses).
				Update("balance", gorm.Expr("CASE"+cases+" END")).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func syncTaskFeeLogs(ctx context.Context, intervalSeconds uint) error {
	db := config.GetDB()

	var checkpoint models.TaskFeeCheckpoint
	err := db.WithContext(ctx).First(&checkpoint).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	batchSize := int(config.GetConfig().Tasks.SyncTaskFeeLogs.BatchSize)
	for {
		logs, err := relay_api.GetTaskFeeLogs(ctx, checkpoint.LatestTaskFeeLogID, batchSize)
		if err != nil {
			return err
		}

		if len(logs) == 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(intervalSeconds) * time.Second):
				continue
			}
		}

		if err := checkTaskFeeLogs(ctx, db, logs); err != nil {
			return err
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
}
