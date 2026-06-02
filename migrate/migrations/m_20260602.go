package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20260602(db *gorm.DB) *gormigrate.Gormigrate {
	type WithdrawRecord struct {
		WithdrawalFee string `json:"withdrawal_fee" gorm:"not null;type:string;size:191;default:'0'"`
	}

	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20260602",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().AddColumn(&WithdrawRecord{}, "WithdrawalFee"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropColumn(&WithdrawRecord{}, "WithdrawalFee"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
