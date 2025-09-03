package models

import (
	"context"
	"crynux_relay_wallet/config"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// TransactionStatus represents the status of a blockchain transaction
type TransactionStatus uint8

const (
	TransactionStatusPending   TransactionStatus = iota // 待发送
	TransactionStatusSent                               // 已发送
	TransactionStatusConfirmed                          // 已确认
	TransactionStatusFailed                             // 失败
)

// BlockchainTransaction represents a blockchain transaction that needs to be sent
type BlockchainTransaction struct {
	gorm.Model
	Network           string            `json:"network" gorm:"index;not null"`
	Type              string            `json:"type" gorm:"index;not null"`
	Status            TransactionStatus `json:"status" gorm:"index;not null;default:0"`
	FromAddress       string            `json:"from_address" gorm:"not null"`
	ToAddress         string            `json:"to_address" gorm:"not null"`
	Value             string            `json:"value" gorm:"not null;default:'0'"`
	Data              sql.NullString    `json:"data" gorm:"null"`
	TxHash            sql.NullString    `json:"tx_hash" gorm:"null;uniqueIndex"`
	BlockNumber       sql.NullInt64     `json:"block_number" gorm:"null"`
	GasUsed           sql.NullInt64     `json:"gas_used" gorm:"null"`
	EffectiveGasPrice sql.NullString    `json:"effective_gas_price" gorm:"null"`
	StatusMessage     sql.NullString    `json:"status_message" gorm:"null"`
	RetryCount        uint8             `json:"retry_count" gorm:"not null;default:0"`
	MaxRetries        uint8             `json:"max_retries" gorm:"not null;default:0"`
	RetryTransactionID sql.NullInt64     `json:"retry_transaction_id" gorm:"null;index"`
	NextRetryAt       sql.NullTime      `json:"next_retry_at" gorm:"null"`
	SentAt            sql.NullTime      `json:"sent_at" gorm:"null"`
	ConfirmedAt       sql.NullTime      `json:"confirmed_at" gorm:"null"`
	FailedAt          sql.NullTime      `json:"failed_at" gorm:"null"`
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

func (tx *BlockchainTransaction) MarkSent(ctx context.Context, db *gorm.DB, txHash string) error {
	updates := map[string]interface{}{
		"status": TransactionStatusSent,
		"tx_hash": txHash,
		"sent_at": time.Now(),
	}

	return tx.Update(ctx, db, updates)
}

func (tx *BlockchainTransaction) MarkConfirmed(ctx context.Context, db *gorm.DB, blockNumber, gasUsed int64, effectiveGasPrice string) error {
	updates := map[string]interface{}{
		"status": TransactionStatusConfirmed,
		"confirmed_at": time.Now(),
		"block_number": blockNumber,
		"gas_used": gasUsed,
		"effective_gas_price": effectiveGasPrice,
	}

	return tx.Update(ctx, db, updates)
}

func (tx *BlockchainTransaction) MarkFailed(ctx context.Context, db *gorm.DB, blockNumber, gasUsed int64, effectiveGasPrice string, errorMsg string) error {
	updates := map[string]interface{}{
		"status": TransactionStatusFailed,
		"failed_at": time.Now(),
		"block_number": blockNumber,
		"gas_used": gasUsed,
		"effective_gas_price": effectiveGasPrice,
		"error_message": errorMsg,
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
		Network: tx.Network,
		Type: tx.Type,
		Status: TransactionStatusPending,
		FromAddress: tx.FromAddress,
		ToAddress: tx.ToAddress,
		Value: tx.Value,
		Data: tx.Data,
		RetryCount: tx.RetryCount + 1,
		MaxRetries: tx.MaxRetries,
		NextRetryAt: sql.NullTime{Time: time.Now().Add(time.Duration(blockchain.RetryInterval) * time.Second)},
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