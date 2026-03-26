package tasks

import (
	"context"
	"crynux_relay_wallet/blockchain"
	"crynux_relay_wallet/config"
	"crynux_relay_wallet/models"
	"crynux_relay_wallet/relay_api"
	"crynux_relay_wallet/utils"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
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
var ErrTaskFeeDepositPayloadInvalid = NewTaskFeeError("deposit log payload is invalid")
var ErrTaskFeeDepositTxMismatch = NewTaskFeeError("deposit transaction does not match relay account event log")

type depositPayload struct {
	TxHash  string `json:"tx_hash"`
	Network string `json:"network"`
}

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

	for _, eventLog := range logs {
		if eventLog.Type == relay_api.TaskFeeLogTypeWithdraw || eventLog.Type == relay_api.TaskFeeLogTypeWithdrawRefund {
			continue
		}

		amount, success := new(big.Int).SetString(eventLog.Amount, 10)
		if !success {
			return nil, ErrTaskFeeAmountInvalid
		}

		if _, ok := merged[eventLog.Address]; !ok {
			merged[eventLog.Address] = big.NewInt(0)
		}

		if eventLog.Type == relay_api.TaskFeeLogTypeTaskPayment {
			merged[eventLog.Address].Sub(merged[eventLog.Address], amount)
		} else {
			merged[eventLog.Address].Add(merged[eventLog.Address], amount)
		}
	}

	return merged, nil
}

func validateDepositLog(ctx context.Context, eventLog relay_api.TaskFeeLog) error {
	amount, success := new(big.Int).SetString(eventLog.Amount, 10)
	if !success {
		return ErrTaskFeeAmountInvalid
	}

	var payload depositPayload
	if strings.TrimSpace(eventLog.Payload) == "" {
		return ErrTaskFeeDepositPayloadInvalid
	}
	if err := json.Unmarshal([]byte(eventLog.Payload), &payload); err != nil {
		return ErrTaskFeeDepositPayloadInvalid
	}
	txHashBytes, err := hexutil.Decode(payload.TxHash)
	if payload.TxHash == "" || payload.Network == "" || err != nil || len(txHashBytes) != common.HashLength {
		return ErrTaskFeeDepositPayloadInvalid
	}

	client, err := blockchain.GetBlockchainClient(payload.Network)
	if err != nil {
		return err
	}

	txHash := common.HexToHash(payload.TxHash)
	receipt, err := client.RpcClient.TransactionReceipt(ctx, txHash)
	if errors.Is(err, ethereum.NotFound) {
		return fmt.Errorf("deposit transaction not found: %s", payload.TxHash)
	}
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return ErrTaskFeeDepositTxMismatch
	}

	tx, _, err := client.RpcClient.TransactionByHash(ctx, txHash)
	if errors.Is(err, ethereum.NotFound) {
		return fmt.Errorf("deposit transaction not found: %s", payload.TxHash)
	}
	if err != nil {
		return err
	}
	if tx.To() == nil || !strings.EqualFold(tx.To().Hex(), config.GetConfig().Relay.DepositAddress) {
		return ErrTaskFeeDepositTxMismatch
	}

	from, err := types.Sender(types.LatestSignerForChainID(client.ChainID), tx)
	if err != nil {
		return err
	}
	if !strings.EqualFold(from.Hex(), eventLog.Address) || tx.Value().Cmp(amount) != 0 {
		return ErrTaskFeeDepositTxMismatch
	}

	return nil
}

func checkTaskFeeLogs(ctx context.Context, db *gorm.DB, logs []relay_api.TaskFeeLog) error {
	appConfig := config.GetConfig()
	maxTaskFeeAmount := utils.EtherToWei(big.NewInt(int64(appConfig.Tasks.SyncTaskFeeLogs.MaxTaskFeeAmount)))

	addressLogCount := make(map[string]int)
	for _, eventLog := range logs {
		if eventLog.Type == relay_api.TaskFeeLogTypeWithdraw || eventLog.Type == relay_api.TaskFeeLogTypeWithdrawRefund {
			continue
		}

		amount, success := new(big.Int).SetString(eventLog.Amount, 10)
		if !success {
			return ErrTaskFeeAmountInvalid
		}
		if eventLog.Type == relay_api.TaskFeeLogTypeDeposit {
			if err := validateDepositLog(ctx, eventLog); err != nil {
				return err
			}
		} else if amount.Cmp(maxTaskFeeAmount) > 0 {
			return ErrTaskFeeAmountTooLarge
		}
		addressLogCount[eventLog.Address]++
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

	if len(addressLogCount) == 0 {
		return nil
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
