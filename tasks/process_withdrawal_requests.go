package tasks

import (
	"context"
	"crynux_relay_wallet/blockchain"
	"crynux_relay_wallet/config"
	"crynux_relay_wallet/models"
	"crynux_relay_wallet/relay_api"
	"database/sql"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrWithdrawalRequestStatusInvalid = errors.New("invalid withdrawal request status")
var ErrWithdrawalRequestAmountInvalid = errors.New("invalid withdrawal request amount")
var ErrWithdrawalRequestAddressNotExists = errors.New("withdrawal request address not exists")
var ErrWithdrawalRequestAmountTooLarge = errors.New("withdrawal request amount is too large")
var ErrWithdrawalRequestTaskFeeNotEnough = errors.New("withdrawal request task fee not enough")
var ErrWithdrawalRequestBalanceNotEnough = errors.New("withdrawal request balance not enough")
var ErrWithdrawalRequestBeneficialAddressInvalid = errors.New("withdrawal request beneficial address is invalid")

func StartSyncWithdrawalRequests(ctx context.Context) error {
	intervalSeconds := config.GetConfig().Tasks.SyncWithdrawalRequests.IntervalSeconds
	interval := time.Duration(intervalSeconds) * time.Second

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Infoln("Sync withdrawal requests task is stopping")
			return nil
		case <-ticker.C:
			if err := syncWithdrawalRequests(ctx, intervalSeconds); err != nil {
				log.Errorf("Failed to sync withdrawal requests: %v", err)
				return err
			}
		}
	}
}

func StartProcessWithdrawalRequests(ctx context.Context) error {
	return processWithdrawalRecords(ctx)
}

func checkWithdrawalRequests(ctx context.Context, db *gorm.DB, requests []relay_api.WithdrawalRequest) error {
	amountMap := make(map[string]*big.Int)
	networkAmountMap := make(map[string]*big.Int)
	for _, request := range requests {
		if request.Status != relay_api.WithdrawStatusPending {
			return ErrWithdrawalRequestStatusInvalid
		}
		amount, ok := big.NewInt(0).SetString(request.Amount, 10)
		if !ok {
			return ErrWithdrawalRequestAmountInvalid
		}
		if _, ok := amountMap[request.Address]; ok {
			amountMap[request.Address].Add(amountMap[request.Address], amount)
		} else {
			amountMap[request.Address] = big.NewInt(0).Set(amount)
		}
		if _, ok := networkAmountMap[request.Network]; ok {
			networkAmountMap[request.Network].Add(networkAmountMap[request.Network], amount)
		} else {
			networkAmountMap[request.Network] = big.NewInt(0).Set(amount)
		}
	}

	for network, amount := range networkAmountMap {
		bc, err := blockchain.GetBlockchainClient(network)
		if err != nil {
			return err
		}
		balance, err := bc.BalanceAt(ctx, common.HexToAddress(bc.Address))
		if err != nil {
			return err
		}
		if balance.Cmp(amount) < 0 {
			return ErrWithdrawalRequestBalanceNotEnough
		}
	}

	addresses := make([]string, 0, len(amountMap))
	for address := range amountMap {
		addresses = append(addresses, address)
	}

	var accounts []*models.RelayAccount
	if err := db.Model(&models.RelayAccount{}).Where("address IN (?)", addresses).Find(&accounts).Error; err != nil {
		return err
	}

	if len(accounts) != len(addresses) {
		return ErrWithdrawalRequestAddressNotExists
	}

	for _, account := range accounts {
		amount := amountMap[account.Address]
		if amount.Cmp(&account.Balance.Int) > 0 {
			return ErrWithdrawalRequestAmountTooLarge
		}
	}

	for _, request := range requests {
		beneficialAddress, err := blockchain.GetBenefitAddress(ctx, common.HexToAddress(request.Address), request.Network)
		if err != nil {
			return err
		}
		if beneficialAddress.Hex() != request.BenefitAddress {
			return ErrWithdrawalRequestBeneficialAddressInvalid
		}
	}

	return nil
}

func syncWithdrawalRequests(ctx context.Context, intervalSeconds uint) error {
	db := config.GetDB()

	var checkpoint models.WithdrawalRequestCheckpoint
	err := db.First(&checkpoint).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	for {
		requests, err := relay_api.GetWithdrawalRequests(ctx, checkpoint.LatestWithdrawalRequestID, int(config.GetConfig().Tasks.SyncWithdrawalRequests.BatchSize))
		if err != nil {
			return err
		}

		if len(requests) == 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(intervalSeconds) * time.Second):
				continue
			}
		}

		if err := checkWithdrawalRequests(ctx, db, requests); err != nil {
			return err
		}

		var records []*models.WithdrawRecord
		for _, request := range requests {
			amount, ok := big.NewInt(0).SetString(request.Amount, 10)
			if !ok {
				return ErrWithdrawalRequestAmountInvalid
			}
			records = append(records, &models.WithdrawRecord{
				RemoteID:       request.ID,
				Address:        request.Address,
				BenefitAddress: request.BenefitAddress,
				Amount:         models.BigInt{Int: *amount},
				Network:        request.Network,
				Status:         models.WithdrawStatusPending,
			})
		}

		if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&records).Error; err != nil {
				return err
			}

			checkpoint.LatestWithdrawalRequestID = requests[len(requests)-1].ID
			checkpoint.LatestWithdrawalRequestTimestamp = requests[len(requests)-1].CreatedAt
			return tx.Save(&checkpoint).Error
		}); err != nil {
			return err
		}
	}
}

func getUnfinishedWithdrawalRecords(ctx context.Context, db *gorm.DB, startID uint, limit int) ([]*models.WithdrawRecord, error) {
	var records []*models.WithdrawRecord
	err := db.WithContext(ctx).Where("status != ?", models.WithdrawStatusFinished).Where("id > ?", startID).Order("id ASC").Limit(limit).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func processWithdrawalRecord(ctx context.Context, db *gorm.DB, record *models.WithdrawRecord) (err error) {
	for record.Status == models.WithdrawStatusPending {

		var blockchainTransaction *models.BlockchainTransaction
		if !record.BlockchainTransactionID.Valid {
			err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
				blockchainTransaction, err = blockchain.QueueSendETH(ctx, tx, common.HexToAddress(record.BenefitAddress), big.NewInt(0).Set(&record.Amount.Int), record.Network)
				if err != nil {
					return err
				}

				record.BlockchainTransactionID = sql.NullInt64{Int64: int64(blockchainTransaction.ID), Valid: true}
				return tx.Save(record).Error
			})
			if err != nil {
				return err
			}
		} else {
			blockchainTransaction, err = models.GetTransactionByID(ctx, db, uint(record.BlockchainTransactionID.Int64))
			if err != nil {
				return err
			}
			if blockchainTransaction.Status == models.TransactionStatusFailed {
				retryTransactions, err := models.GetRetryTransactionsByID(ctx, db, uint(record.BlockchainTransactionID.Int64))
				if err != nil {
					return err
				}
				if len(retryTransactions) > 0 {
					blockchainTransaction = &retryTransactions[len(retryTransactions)-1]
				}
			}
		}

		for blockchainTransaction.Status != models.TransactionStatusConfirmed && blockchainTransaction.Status != models.TransactionStatusFailed {
			err = blockchainTransaction.Sync(ctx, db)
			if err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(1 * time.Second):
				continue
			}
		}

		if blockchainTransaction.Status == models.TransactionStatusConfirmed {
			err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
				var account models.RelayAccount
				err = tx.Model(&models.RelayAccount{}).Where("address = ?", record.Address).First(&account).Error
				if err != nil {
					return err
				}
				if account.Balance.Cmp(&record.Amount.Int) < 0 {
					return ErrWithdrawalRequestTaskFeeNotEnough
				}
				account.Balance.Sub(&account.Balance.Int, &record.Amount.Int)
				err = tx.Save(&account).Error
				if err != nil {
					return err
				}
				err = record.UpdateStatus(ctx, tx, models.WithdrawStatusSuccess)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
		} else if blockchainTransaction.RetryCount >= blockchainTransaction.MaxRetries {
			err = record.UpdateStatus(ctx, db, models.WithdrawStatusFailed)
			if err != nil {
				return err
			}
		}
	}

	if record.Status == models.WithdrawStatusSuccess {
		err = relay_api.FulfillWithdrawalRequest(ctx, record.RemoteID)
		if err != nil {
			return err
		}
	} else {
		err = relay_api.RejectWithdrawalRequest(ctx, record.RemoteID)
		if err != nil {
			return err
		}
	}

	err = record.UpdateStatus(ctx, db, models.WithdrawStatusFinished)
	if err != nil {
		return err
	}
	return nil
}

func processWithdrawalRecords(ctx context.Context) error {
	appConfig := config.GetConfig()
	db := config.GetDB()

	var startID uint
	limit := appConfig.Tasks.ProcessWithdrawalRequests.BatchSize

	for {
		records, err := getUnfinishedWithdrawalRecords(ctx, db, startID, int(limit))
		if err != nil {
			return err
		}

		if len(records) > 0 {
			startID = records[len(records)-1].ID
			for _, record := range records {
				go func(ctx context.Context, record *models.WithdrawRecord) {
					deadline := record.CreatedAt.Add(time.Duration(config.GetConfig().Tasks.ProcessWithdrawalRequests.Timeout) * time.Second)
					ctx, cancel := context.WithDeadline(ctx, deadline)
					defer cancel()

					for {
						log.Infof("ProcessWithdrawalRecords: process withdrawal record %d", record.ID)
						c := make(chan error)
						go func() {
							c <- processWithdrawalRecord(ctx, db, record)
						}()

						select {
						case <-ctx.Done():
							log.Infof("ProcessWithdrawalRecords: process withdrawal record %d timeout", record.ID)
							if err = relay_api.RejectWithdrawalRequest(ctx, record.RemoteID); err != nil {
								log.Errorf("ProcessWithdrawalRecords: reject timeout withdrawal record %d error %v", record.ID, err)
								time.Sleep(5 * time.Second)
								continue
							}
							if err := record.UpdateStatus(ctx, db, models.WithdrawStatusFinished); err != nil {
								log.Errorf("ProcessWithdrawalRecords: update timeout withdrawal record %d status error %v", record.ID, err)
								time.Sleep(5 * time.Second)
								continue
							}
						case err := <-c:
							if err != nil {
								log.Errorf("ProcessWithdrawalRecords: process withdrawal record %d error %v", record.ID, err)
								time.Sleep(5 * time.Second)
							} else {
								log.Infof("ProcessWithdrawalRecords: process withdrawal record %d successfully", record.ID)
								return
							}
						}
					}
				}(ctx, record)
			}
		}
		time.Sleep(time.Duration(appConfig.Tasks.ProcessWithdrawalRequests.IntervalSeconds) * time.Second)
	}

}
