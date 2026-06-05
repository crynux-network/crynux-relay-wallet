package tasks

import (
	"context"
	"crynux_relay_wallet/blockchain"
	"crynux_relay_wallet/config"
	"crynux_relay_wallet/models"
	"crynux_relay_wallet/relay_api"
	"crynux_relay_wallet/utils"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WithdrawalRequestError struct {
	Message string
}

func (e *WithdrawalRequestError) Error() string {
	return e.Message
}

func NewWithdrawalRequestError(message string) *WithdrawalRequestError {
	return &WithdrawalRequestError{Message: message}
}

func IsWithdrawalRequestError(err error) bool {
	var withdrawalRequestError *WithdrawalRequestError
	return errors.As(err, &withdrawalRequestError)
}

var ErrWithdrawalRequestStatusInvalid = NewWithdrawalRequestError("invalid withdrawal request status")
var ErrWithdrawalRequestAmountInvalid = NewWithdrawalRequestError("invalid withdrawal request amount")
var ErrWithdrawalRequestAddressNotExists = NewWithdrawalRequestError("withdrawal request address not exists")
var ErrWithdrawalRequestAmountTooLarge = NewWithdrawalRequestError("withdrawal request amount is too large")
var ErrWithdrawalRequestTaskFeeNotEnough = NewWithdrawalRequestError("withdrawal request task fee not enough")
var ErrWithdrawalRequestRelayHotWalletNativeBalanceInsufficient = NewWithdrawalRequestError("relay hot wallet native token balance is insufficient")
var ErrWithdrawalRequestRelayHotWalletERC20BalanceInsufficient = NewWithdrawalRequestError("relay hot wallet ERC20 token balance is insufficient")
var ErrWithdrawalRequestRelayHotWalletNativeGasFeeInsufficient = NewWithdrawalRequestError("relay hot wallet native gas fee is insufficient")
var ErrWithdrawalRequestBeneficialAddressInvalid = NewWithdrawalRequestError("withdrawal request beneficial address is invalid")
var ErrWithdrawalRequestAmountTooSmall = NewWithdrawalRequestError("withdrawal request amount is too small")
var ErrWithdrawalRequestTransactionUnconfirmedTimeout = NewWithdrawalRequestError("withdrawal request transaction remains unconfirmed after timeout")

func parseWithdrawalAmount(amountText string) (*big.Int, error) {
	amount, ok := big.NewInt(0).SetString(amountText, 10)
	if !ok || amount.Sign() < 0 {
		return nil, ErrWithdrawalRequestAmountInvalid
	}
	return amount, nil
}

func withdrawalTotalAmount(amount, withdrawalFee *big.Int) *big.Int {
	return big.NewInt(0).Add(amount, withdrawalFee)
}

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
				if IsWithdrawalRequestError(err) {
					return err
				}
			}
		}
	}
}

func StartProcessWithdrawalRequests(ctx context.Context) error {
	intervalSeconds := config.GetConfig().Tasks.ProcessWithdrawalRequests.IntervalSeconds
	interval := time.Duration(intervalSeconds) * time.Second

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Infoln("Process withdrawal requests task is stopping")
			return nil
		case <-ticker.C:
			if err := processWithdrawalRecords(ctx); err != nil {
				log.Errorf("Failed to process withdrawal requests: %v", err)
				if IsWithdrawalRequestError(err) {
					return err
				}
			}
		}
	}
}

func checkWithdrawalRequests(ctx context.Context, db *gorm.DB, requests []relay_api.WithdrawalRequest) error {
	appConfig := config.GetConfig()
	minWithdrawalAmount := utils.EtherToWei(big.NewInt(0).SetUint64(appConfig.Tasks.SyncWithdrawalRequests.MinWithdrawalAmount))
	for _, request := range requests {
		amount, err := parseWithdrawalAmount(request.Amount)
		if err != nil {
			return err
		}
		withdrawalFee, err := parseWithdrawalAmount(request.WithdrawalFee)
		if err != nil {
			return err
		}
		if amount.Cmp(minWithdrawalAmount) < 0 {
			return ErrWithdrawalRequestAmountTooSmall
		}
		if withdrawalTotalAmount(amount, withdrawalFee).Sign() <= 0 {
			return ErrWithdrawalRequestAmountInvalid
		}
	}

	amountMap := make(map[string]*big.Int)
	networkAmountMap := make(map[string]*big.Int)
	for _, request := range requests {
		if request.Status != relay_api.WithdrawStatusPending {
			return ErrWithdrawalRequestStatusInvalid
		}
		amount, err := parseWithdrawalAmount(request.Amount)
		if err != nil {
			return err
		}
		withdrawalFee, err := parseWithdrawalAmount(request.WithdrawalFee)
		if err != nil {
			return err
		}
		totalAmount := withdrawalTotalAmount(amount, withdrawalFee)
		if _, ok := amountMap[request.Address]; ok {
			amountMap[request.Address].Add(amountMap[request.Address], totalAmount)
		} else {
			amountMap[request.Address] = big.NewInt(0).Set(totalAmount)
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
		blockchainConfig := appConfig.Blockchains[network]
		switch blockchainConfig.TokenType {
		case config.TokenTypeNative:
			balance, err := bc.BalanceAt(ctx, common.HexToAddress(bc.Address))
			if err != nil {
				return err
			}
			if balance.Cmp(amount) < 0 {
				return ErrWithdrawalRequestRelayHotWalletNativeBalanceInsufficient
			}
		case config.TokenTypeERC20:
			tokenBalance, err := bc.TokenBalanceAt(ctx, common.HexToAddress(blockchainConfig.TokenAddress), common.HexToAddress(bc.Address))
			if err != nil {
				return err
			}
			if tokenBalance.Cmp(amount) < 0 {
				return ErrWithdrawalRequestRelayHotWalletERC20BalanceInsufficient
			}
			gasBalance, err := bc.BalanceAt(ctx, common.HexToAddress(bc.Address))
			if err != nil {
				return err
			}
			if gasBalance.Sign() <= 0 {
				return ErrWithdrawalRequestRelayHotWalletNativeGasFeeInsufficient
			}
		default:
			return blockchain.ErrBlockchainNotFound
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
		ba, err := blockchain.GetBenefitAddress(ctx, common.HexToAddress(request.Address), request.Network)
		if err != nil {
			return err
		}
		if ba.Hex() != request.BenefitAddress {
			return ErrWithdrawalRequestBeneficialAddressInvalid
		}
	}

	return nil
}

func syncWithdrawalRequests(ctx context.Context, intervalSeconds uint) error {
	db := config.GetDB()

	var checkpoint models.WithdrawalRequestCheckpoint
	err := db.WithContext(ctx).First(&checkpoint).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	for {
		var taskFeeCheckpoint models.TaskFeeCheckpoint
		err = db.WithContext(ctx).First(&taskFeeCheckpoint).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		requests, err := relay_api.GetWithdrawalRequests(ctx, checkpoint.LatestWithdrawalRequestID, int(config.GetConfig().Tasks.SyncWithdrawalRequests.BatchSize))
		if err != nil {
			return err
		}

		end := 0
		for _, request := range requests {
			if request.RelayAccountEventID > taskFeeCheckpoint.LatestTaskFeeLogID {
				break
			}
			end++
		}
		requests = requests[:end]

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
			amount, err := parseWithdrawalAmount(request.Amount)
			if err != nil {
				return err
			}
			withdrawalFee, err := parseWithdrawalAmount(request.WithdrawalFee)
			if err != nil {
				return err
			}
			records = append(records, &models.WithdrawRecord{
				RemoteID:       request.ID,
				Address:        request.Address,
				BenefitAddress: request.BenefitAddress,
				Amount:         models.BigInt{Int: *amount},
				WithdrawalFee:  models.BigInt{Int: *withdrawalFee},
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
	var blockchainTransaction *models.BlockchainTransaction
	for record.Status == models.WithdrawStatusPending {

		if !record.BlockchainTransactionID.Valid {
			err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
				var toAddress common.Address
				if record.BenefitAddress != "" {
					toAddress = common.HexToAddress(record.BenefitAddress)
				} else {
					toAddress = common.HexToAddress(record.Address)
				}
				blockchainConfig := config.GetConfig().Blockchains[record.Network]
				switch blockchainConfig.TokenType {
				case config.TokenTypeNative:
					blockchainTransaction, err = blockchain.QueueSendETH(ctx, tx, toAddress, big.NewInt(0).Set(&record.Amount.Int), record.Network)
				case config.TokenTypeERC20:
					blockchainTransaction, err = blockchain.QueueSendERC20(ctx, tx, common.HexToAddress(blockchainConfig.TokenAddress), toAddress, big.NewInt(0).Set(&record.Amount.Int), record.Network)
				default:
					err = blockchain.ErrBlockchainNotFound
				}
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
				totalAmount := withdrawalTotalAmount(&record.Amount.Int, &record.WithdrawalFee.Int)
				var account models.RelayAccount
				err = tx.Model(&models.RelayAccount{}).Where("address = ?", record.Address).First(&account).Error
				if err != nil {
					return err
				}
				if account.Balance.Cmp(totalAmount) < 0 {
					return ErrWithdrawalRequestTaskFeeNotEnough
				}
				account.Balance.Sub(&account.Balance.Int, totalAmount)
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
		err = relay_api.FulfillWithdrawalRequest(ctx, record.RemoteID, blockchainTransaction.TxHash.String)
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

func rejectTimeoutWithdrawalRequest(ctx context.Context, db *gorm.DB, record *models.WithdrawRecord) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	log.Infof("ProcessWithdrawalRecords: process withdrawal record %d timeout", record.ID)
	if err := relay_api.RejectWithdrawalRequest(ctx, record.RemoteID); err != nil {
		log.Errorf("ProcessWithdrawalRecords: reject timeout withdrawal record %d error %v", record.ID, err)
		return err
	}
	if err := record.UpdateStatus(ctx, db, models.WithdrawStatusFinished); err != nil {
		log.Errorf("ProcessWithdrawalRecords: update timeout withdrawal record %d status error %v", record.ID, err)
		return err
	}
	return nil
}

func handleTimeoutWithdrawalRequest(ctx context.Context, db *gorm.DB, record *models.WithdrawRecord) error {
	if err := db.WithContext(ctx).First(record, record.ID).Error; err != nil {
		return err
	}
	if record.BlockchainTransactionID.Valid {
		log.Errorf("ProcessWithdrawalRecords: withdrawal record %d has unconfirmed blockchain transaction %d after timeout", record.ID, record.BlockchainTransactionID.Int64)
		return ErrWithdrawalRequestTransactionUnconfirmedTimeout
	}
	if err := rejectTimeoutWithdrawalRequest(context.Background(), db, record); err != nil {
		return fmt.Errorf("ProcessWithdrawalRecords: cannot reject timeout withdrawal request due to %w", err)
	}
	log.Infof("ProcessWithdrawalRecords: rejected timeout withdrawal record %d before blockchain transaction creation", record.ID)
	return nil
}

func processWithdrawalRecordWithRetry(ctx context.Context, db *gorm.DB, record *models.WithdrawRecord) error {
	deadline := record.CreatedAt.Add(time.Duration(config.GetConfig().Tasks.ProcessWithdrawalRequests.Timeout) * time.Second)

	for {
		if time.Now().After(deadline) {
			return handleTimeoutWithdrawalRequest(ctx, db, record)
		}

		recordCtx, cancel := context.WithDeadline(ctx, deadline)
		log.Infof("ProcessWithdrawalRecords: process withdrawal record %d", record.ID)
		err := processWithdrawalRecord(recordCtx, db, record)
		cancel()

		if err == nil {
			log.Infof("ProcessWithdrawalRecords: process withdrawal record %d successfully", record.ID)
			return nil
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return handleTimeoutWithdrawalRequest(ctx, db, record)
		}
		log.Errorf("ProcessWithdrawalRecords: process withdrawal record %d error %v", record.ID, err)
		if IsWithdrawalRequestError(err) {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
}

func processWithdrawalRecords(ctx context.Context) error {
	appConfig := config.GetConfig()
	db := config.GetDB()

	var startID uint
	limit := appConfig.Tasks.ProcessWithdrawalRequests.BatchSize
	interval := time.Duration(appConfig.Tasks.ProcessWithdrawalRequests.IntervalSeconds) * time.Second

	for {
		records, err := getUnfinishedWithdrawalRecords(ctx, db, startID, int(limit))
		if err != nil {
			return err
		}

		if len(records) == 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(interval):
				continue
			}
		}

		for _, record := range records {
			if err := processWithdrawalRecordWithRetry(ctx, db, record); err != nil {
				return err
			}
			startID = record.ID
		}
	}
}
