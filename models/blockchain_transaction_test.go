package models

import (
	"context"
	"database/sql"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newBlockchainTransactionTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&BlockchainTransaction{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	return db
}

func TestBlockchainTransactionClaimReleaseAndCancel(t *testing.T) {
	ctx := context.Background()
	db := newBlockchainTransactionTestDB(t)
	transaction := &BlockchainTransaction{
		Network:     "testnet",
		Type:        "SendETH",
		Status:      TransactionStatusPending,
		FromAddress: "0x1111111111111111111111111111111111111111",
		ToAddress:   "0x2222222222222222222222222222222222222222",
		Value:       "1",
	}
	if err := transaction.Save(ctx, db); err != nil {
		t.Fatalf("save transaction: %v", err)
	}

	claimed, err := transaction.ClaimForSending(ctx, db)
	if err != nil {
		t.Fatalf("claim transaction: %v", err)
	}
	if !claimed {
		t.Fatalf("expected transaction to be claimed")
	}

	cancelled, err := transaction.CancelUnbroadcasted(ctx, db, "test cancellation")
	if err != nil {
		t.Fatalf("cancel claimed transaction: %v", err)
	}
	if cancelled {
		t.Fatalf("claimed transaction must not be cancellable")
	}

	if err := transaction.ReleaseSending(ctx, db, "temporary send error"); err != nil {
		t.Fatalf("release transaction: %v", err)
	}

	cancelled, err = transaction.CancelUnbroadcasted(ctx, db, "test cancellation")
	if err != nil {
		t.Fatalf("cancel pending transaction: %v", err)
	}
	if !cancelled {
		t.Fatalf("expected pending transaction to be cancelled")
	}

	var saved BlockchainTransaction
	if err := db.First(&saved, transaction.ID).Error; err != nil {
		t.Fatalf("load transaction: %v", err)
	}
	if saved.Status != TransactionStatusCancelled {
		t.Fatalf("expected cancelled status, got %d", saved.Status)
	}
	if saved.TxHash.Valid {
		t.Fatalf("cancelled unbroadcasted transaction must not have tx hash")
	}
}

func TestBlockchainTransactionClaimAndCancelRequireNoTxHash(t *testing.T) {
	ctx := context.Background()
	db := newBlockchainTransactionTestDB(t)
	transaction := &BlockchainTransaction{
		Network:     "testnet",
		Type:        "SendETH",
		Status:      TransactionStatusPending,
		FromAddress: "0x1111111111111111111111111111111111111111",
		ToAddress:   "0x2222222222222222222222222222222222222222",
		Value:       "1",
		TxHash:      sql.NullString{String: "0x1234", Valid: true},
	}
	if err := transaction.Save(ctx, db); err != nil {
		t.Fatalf("save transaction: %v", err)
	}

	claimed, err := transaction.ClaimForSending(ctx, db)
	if err != nil {
		t.Fatalf("claim transaction: %v", err)
	}
	if claimed {
		t.Fatalf("transaction with tx hash must not be claimed")
	}

	cancelled, err := transaction.CancelUnbroadcasted(ctx, db, "test cancellation")
	if err != nil {
		t.Fatalf("cancel transaction: %v", err)
	}
	if cancelled {
		t.Fatalf("transaction with tx hash must not be cancelled")
	}
}
