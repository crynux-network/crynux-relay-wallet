package migrations

import (
	"database/sql"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20250902(db *gorm.DB) *gormigrate.Gormigrate {
	type TransactionStatus uint8

	type BlockchainTransaction struct {
		gorm.Model
		Network            string            `json:"network" gorm:"type:string;size:191;index;not null"`
		Type               string            `json:"type" gorm:"type:string;size:191;index;not null"`
		Status             TransactionStatus `json:"status" gorm:"index;not null;default:0"`
		FromAddress        string            `json:"from_address" gorm:"type:string;size:191;not null"`
		ToAddress          string            `json:"to_address" gorm:"type:string;size:191;not null"`
		Value              string            `json:"value" gorm:"type:string;size:191;not null;default:'0'"`
		Data               sql.NullString    `json:"data" gorm:"null"`
		TxHash             sql.NullString    `json:"tx_hash" gorm:"type:string;size:191;null;uniqueIndex"`
		BlockNumber        sql.NullInt64     `json:"block_number" gorm:"null"`
		GasUsed            sql.NullInt64     `json:"gas_used" gorm:"null"`
		EffectiveGasPrice  sql.NullString    `json:"effective_gas_price" gorm:"type:string;size:191;null"`
		StatusMessage      sql.NullString    `json:"status_message" gorm:"null"`
		RetryCount         uint8             `json:"retry_count" gorm:"not null;default:0"`
		MaxRetries         uint8             `json:"max_retries" gorm:"not null;default:0"`
		RetryTransactionID sql.NullInt64     `json:"retry_transaction_id" gorm:"null;index"`
		NextRetryAt        sql.NullTime      `json:"next_retry_at" gorm:"null"`
		SentAt             sql.NullTime      `json:"sent_at" gorm:"null"`
		ConfirmedAt        sql.NullTime      `json:"confirmed_at" gorm:"null"`
		FailedAt           sql.NullTime      `json:"failed_at" gorm:"null"`
	}

	type TaskFeeCheckpoint struct {
		ID                        uint   `json:"id" gorm:"primarykey"`
		LatestTaskFeeLogID        uint   `json:"latest_task_fee_log_id"`
		LatestTaskFeeLogTimestamp uint64 `json:"latest_task_fee_log_timestamp"`
	}

	type WithdrawalRequestCheckpoint struct {
		ID                               uint   `json:"id" gorm:"primarykey"`
		LatestWithdrawalRequestID        uint   `json:"latest_withdrawal_request_id"`
		LatestWithdrawalRequestTimestamp uint64 `json:"latest_withdrawal_request_timestamp"`
	}

	type WithdrawStatus int8

	type WithdrawRecord struct {
		gorm.Model
		RemoteID                uint           `json:"remote_id" gorm:"not null;uniqueIndex"`
		Address                 string         `json:"address" gorm:"type:string;size:191;not null;index"`
		BenefitAddress          string         `json:"benefit_address" gorm:"type:string;size:191;not null;index"`
		Amount                  string         `json:"amount" gorm:"type:string;size:191;not null"`
		Network                 string         `json:"network" gorm:"type:string;size:191;not null;index"`
		Status                  WithdrawStatus `json:"status" gorm:"not null;default:0;index"`
		BlockchainTransactionID sql.NullInt64  `json:"blockchain_transaction_id" gorm:"index"`
	}

	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20250902",
			Migrate: func(tx *gorm.DB) error {

				if err := tx.Migrator().CreateTable(&BlockchainTransaction{}); err != nil {
					return err
				}
				if err := tx.Migrator().CreateTable(&TaskFeeCheckpoint{}); err != nil {
					return err
				}
				if err := tx.Migrator().CreateTable(&WithdrawalRequestCheckpoint{}); err != nil {
					return err
				}
				if err := tx.Migrator().CreateTable(&WithdrawRecord{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropTable("blockchain_transactions"); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable("task_fee_checkpoints"); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable("withdrawal_request_checkpoints"); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable("withdraw_records"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
