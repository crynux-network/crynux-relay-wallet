package models

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type WithdrawRecord struct {
	gorm.Model
	RemoteID                uint           `json:"remote_id" gorm:"not null;uniqueIndex"`
	Address                 string         `json:"address" gorm:"not null;index"`
	BenefitAddress          string         `json:"benefit_address" gorm:"not null;index"`
	Amount                  BigInt         `json:"amount" gorm:"not null"`
	Network                 string         `json:"network" gorm:"not null;index"`
	Status                  WithdrawStatus `json:"status" gorm:"not null;default:0;index"`
	BlockchainTransactionID sql.NullInt64  `json:"blockchain_transaction_id" gorm:"index"`
	WithdrawalFee           BigInt         `json:"withdrawal_fee" gorm:"not null"`
}

type WithdrawStatus int8

const (
	WithdrawStatusPending WithdrawStatus = iota
	WithdrawStatusSuccess
	WithdrawStatusFailed
	WithdrawStatusFinished
)

func (w *WithdrawRecord) UpdateStatus(ctx context.Context, db *gorm.DB, status WithdrawStatus) error {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return db.WithContext(dbCtx).Model(w).Update("status", status).Error
}
