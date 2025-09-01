package models

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type WithdrawRecord struct {
	gorm.Model
	RemoteID                uint                  `json:"remote_id" gorm:"not null;uniqueIndex"`
	Address                 string                `json:"address" gorm:"not null;index"`
	BenefitAddress          string                `json:"benefit_address" gorm:"not null;index"`
	Amount                  BigInt                `json:"amount" gorm:"not null"`
	Network                 string                `json:"network" gorm:"not null;index"`
	Status                  WithdrawStatus        `json:"status" gorm:"not null;default:0;index"`
}

type WithdrawStatus int8

const (
	WithdrawStatusPending WithdrawStatus = iota
	WithdrawStatusSuccess
	WithdrawStatusFailed
	WithdrawStatusFinished
)

func (w *WithdrawRecord) GetLatestBlockchainTransaction(ctx context.Context, db *gorm.DB) (*BlockchainTransaction, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var transaction BlockchainTransaction
	if err := db.WithContext(dbCtx).Model(&BlockchainTransaction{}).Where("withdraw_id = ?", w.ID).Last(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, nil
}

func (w *WithdrawRecord) UpdateStatus(ctx context.Context, db *gorm.DB, status WithdrawStatus) error {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return db.WithContext(dbCtx).Model(w).Update("status", status).Error
}