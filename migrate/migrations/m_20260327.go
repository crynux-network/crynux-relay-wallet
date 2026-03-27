package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20260327(db *gorm.DB) *gormigrate.Gormigrate {
	type DepositRecord struct {
		gorm.Model
		Network             string `json:"network" gorm:"type:string;size:191;not null;uniqueIndex:idx_deposit_record_network_tx_hash"`
		TxHash              string `json:"tx_hash" gorm:"type:string;size:191;not null;uniqueIndex:idx_deposit_record_network_tx_hash"`
		DepositAddress      string `json:"deposit_address" gorm:"type:string;size:191;not null"`
		FromAddress         string `json:"from_address" gorm:"type:string;size:191;not null;index"`
		Amount              string `json:"amount" gorm:"type:string;size:255;not null"`
		RelayAccountEventID uint   `json:"relay_account_event_id" gorm:"not null;index"`
	}

	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20260327",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().CreateTable(&DepositRecord{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropTable("deposit_records"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
