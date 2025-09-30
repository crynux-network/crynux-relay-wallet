package migrations

import (
	"database/sql"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20250930(db *gorm.DB) *gormigrate.Gormigrate {
	type BlockchainTransaction struct {
		Nonce sql.NullInt64 `json:"nonce" gorm:"null"`
	}

	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20250930",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().AddColumn(&BlockchainTransaction{}, "Nonce"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropColumn(&BlockchainTransaction{}, "Nonce"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
