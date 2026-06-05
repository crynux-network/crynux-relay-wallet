package models

import (
	"context"
	"crynux_relay_wallet/config"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// TransactionStatus represents the status of a blockchain transaction.
type TransactionStatus uint8

const (
	TransactionStatusPending TransactionStatus = iota
	TransactionStatusSent
	TransactionStatusConfirmed
	TransactionStatusFailed
	TransactionStatusSending
	TransactionStatusCancelled
)

// BlockchainTransaction represents a blockchain transaction that needs to be sent
type BlockchainTransaction struct {
	gorm.Model
	Network            string            `json:"network" gorm:"index;not null"`
	Type               string            `json:"type" gorm:"index;not null"`
	Status             TransactionStatus `json:"status" gorm:"index;not null;default:0"`
	FromAddress        string            `json:"from_address" gorm:"not null"`
	ToAddress          string            `json:"to_address" gorm:"not null"`
	Value              string            `json:"value" gorm:"not null;default:'0'"`
	Data               sql.NullString    `json:"data" gorm:"null"`
	TxHash             sql.NullString    `json:"tx_hash" gorm:"null;uniqueIndex"`
	BlockNumber        sql.NullInt64     `json:"block_number" gorm:"null"`
	Nonce              sql.NullInt64     `json:"nonce" gorm:"null"`
	GasUsed            sql.NullInt64     `json:"gas_used" gorm:"null"`
	EffectiveGasPrice  sql.NullString    `json:"effective_gas_price" gorm:"null"`
	StatusMessage      sql.NullString    `json:"status_message" gorm:"null"`
	RetryCount         uint8             `json:"retry_count" gorm:"not null;default:0"`
	MaxRetries         uint8             `json:"max_retries" gorm:"not null;default:0"`
	RetryTransactionID sql.NullInt64     `json:"retry_transaction_id" gorm:"null;index"`
	NextRetryAt        sql.NullTime      `json:"next_retry_at" gorm:"null"`
	SentAt             sql.NullTime      `json:"sent_at" gorm:"null"`
	ConfirmedAt        sql.NullTime      `json:"confirmed_at" gorm:"null"`
	FailedAt           sql.NullTime      `json:"failed_at" gorm:"null"`
}

// TableName returns the table name for BlockchainTransaction
func (BlockchainTransaction) TableName() string {
	return "blockchain_transactions"
}

// Save saves the blockchain transaction to database
func (tx *BlockchainTransaction) Save(ctx context.Context, db *gorm.DB) error {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.WithContext(dbCtx).Save(&tx).Error; err != nil {
		return err
	}
	return nil
}

func (tx *BlockchainTransaction) Sync(ctx context.Context, db *gorm.DB) error {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.WithContext(dbCtx).Model(tx).First(tx).Error; err != nil {
		return err
	}
	return nil
}

// Update updates the blockchain transaction in database
func (tx *BlockchainTransaction) Update(ctx context.Context, db *gorm.DB, values map[string]interface{}) error {
	if tx.ID == 0 {
		return gorm.ErrRecordNotFound
	}
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.WithContext(dbCtx).Model(tx).Updates(values).Error; err != nil {
		return err
	}
	return nil
}

// GetPendingTransactions gets all pending transactions from database
func GetPendingTransactions(ctx context.Context, db *gorm.DB, offset, limit int) ([]BlockchainTransaction, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var transactions []BlockchainTransaction
	if err := db.WithContext(dbCtx).Where("status = ?", TransactionStatusPending).Order("id").Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// GetSentTransactions gets all sent transactions that need confirmation from database
func GetSentTransactions(ctx context.Context, db *gorm.DB, offset, limit int) ([]BlockchainTransaction, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var transactions []BlockchainTransaction
	if err := db.WithContext(dbCtx).Where("status = ?", TransactionStatusSent).Order("id").Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func GetSentTransactionCountByNetwork(ctx context.Context, db *gorm.DB, network string) (int64, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var count int64
	if err := db.WithContext(dbCtx).Model(&BlockchainTransaction{}).Where("network = ?", network).Where("status = ?", TransactionStatusSent).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetTransactionByHash gets a transaction by its hash
func GetTransactionByHash(ctx context.Context, db *gorm.DB, txHash string) (*BlockchainTransaction, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var transaction BlockchainTransaction
	if err := db.WithContext(dbCtx).Where("tx_hash = ?", txHash).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func GetTransactionByID(ctx context.Context, db *gorm.DB, id uint) (*BlockchainTransaction, error) {
	var transaction BlockchainTransaction
	if err := db.WithContext(ctx).First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (tx *BlockchainTransaction) MarkSent(ctx context.Context, db *gorm.DB, txHash string, nonce int64) error {
	updates := map[string]interface{}{
		"status":  TransactionStatusSent,
		"tx_hash": txHash,
		"nonce":   nonce,
		"sent_at": time.Now(),
	}

	return tx.Update(ctx, db, updates)
}

func (tx *BlockchainTransaction) ClaimForSending(ctx context.Context, db *gorm.DB) (bool, error) {
	if tx.ID == 0 {
		return false, gorm.ErrRecordNotFound
	}
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := db.WithContext(dbCtx).Model(&BlockchainTransaction{}).
		Where("id = ? AND status = ? AND tx_hash IS NULL", tx.ID, TransactionStatusPending).
		Updates(map[string]interface{}{
			"status": TransactionStatusSending,
		})
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	tx.Status = TransactionStatusSending
	return true, nil
}

func (tx *BlockchainTransaction) ReleaseSending(ctx context.Context, db *gorm.DB, errorMsg string) error {
	if tx.ID == 0 {
		return gorm.ErrRecordNotFound
	}
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	updates := map[string]interface{}{
		"status": TransactionStatusPending,
	}
	if errorMsg != "" {
		updates["status_message"] = errorMsg
	}
	result := db.WithContext(dbCtx).Model(&BlockchainTransaction{}).
		Where("id = ? AND status = ? AND tx_hash IS NULL", tx.ID, TransactionStatusSending).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		tx.Status = TransactionStatusPending
	}
	return nil
}

func (tx *BlockchainTransaction) CancelUnbroadcasted(ctx context.Context, db *gorm.DB, reason string) (bool, error) {
	if tx.ID == 0 {
		return false, gorm.ErrRecordNotFound
	}
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	updates := map[string]interface{}{
		"status":         TransactionStatusCancelled,
		"failed_at":      time.Now(),
		"status_message": reason,
	}
	result := db.WithContext(dbCtx).Model(&BlockchainTransaction{}).
		Where("id = ? AND status = ? AND tx_hash IS NULL", tx.ID, TransactionStatusPending).
		Updates(updates)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	tx.Status = TransactionStatusCancelled
	tx.FailedAt = sql.NullTime{Time: updates["failed_at"].(time.Time), Valid: true}
	tx.StatusMessage = sql.NullString{String: reason, Valid: true}
	return true, nil
}

func (tx *BlockchainTransaction) MarkConfirmed(ctx context.Context, db *gorm.DB, blockNumber, gasUsed int64, effectiveGasPrice string) error {
	updates := map[string]interface{}{
		"status":              TransactionStatusConfirmed,
		"confirmed_at":        time.Now(),
		"block_number":        blockNumber,
		"gas_used":            gasUsed,
		"effective_gas_price": effectiveGasPrice,
	}

	return tx.Update(ctx, db, updates)
}

func (tx *BlockchainTransaction) MarkFailed(ctx context.Context, db *gorm.DB, blockNumber, gasUsed int64, effectiveGasPrice string, errorMsg string) error {
	updates := map[string]interface{}{
		"status":              TransactionStatusFailed,
		"failed_at":           time.Now(),
		"block_number":        blockNumber,
		"gas_used":            gasUsed,
		"effective_gas_price": effectiveGasPrice,
		"status_message":      errorMsg,
	}
	// make tx_hash unique
	if tx.TxHash.Valid {
		updates["tx_hash"] = tx.TxHash.String + "_" + strconv.FormatInt(int64(tx.RetryCount), 10)
	}

	return tx.Update(ctx, db, updates)
}

func (tx *BlockchainTransaction) CreateRetryTransaction(ctx context.Context, db *gorm.DB) error {
	appConfig := config.GetConfig()
	blockchain, ok := appConfig.Blockchains[tx.Network]
	if !ok {
		return fmt.Errorf("network %s not found", tx.Network)
	}
	var retryTransactionID sql.NullInt64
	if tx.RetryTransactionID.Valid {
		retryTransactionID = tx.RetryTransactionID
	} else {
		retryTransactionID = sql.NullInt64{Int64: int64(tx.ID), Valid: true}
	}
	nextTransaction := &BlockchainTransaction{
		Network:            tx.Network,
		Type:               tx.Type,
		Status:             TransactionStatusPending,
		FromAddress:        tx.FromAddress,
		ToAddress:          tx.ToAddress,
		Value:              tx.Value,
		Data:               tx.Data,
		RetryCount:         tx.RetryCount + 1,
		MaxRetries:         tx.MaxRetries,
		NextRetryAt:        sql.NullTime{Time: time.Now().Add(time.Duration(blockchain.RetryInterval) * time.Second)},
		RetryTransactionID: retryTransactionID,
	}
	if err := nextTransaction.Save(ctx, db); err != nil {
		return err
	}
	return nil
}

func GetRetryTransactionsByID(ctx context.Context, db *gorm.DB, id uint) ([]BlockchainTransaction, error) {
	var transactions []BlockchainTransaction
	if err := db.WithContext(ctx).Model(&BlockchainTransaction{}).Where("retry_transaction_id = ?", id).Order("id").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
