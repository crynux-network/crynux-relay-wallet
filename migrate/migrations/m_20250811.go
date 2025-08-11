package migrations

import (
	"crynux_relay_wallet/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20250811(db *gorm.DB) *gormigrate.Gormigrate {

	type RelayAccounts struct {
		gorm.Model
		Address string        `json:"address" gorm:"uniqueIndex"`
		Balance models.BigInt `json:"balance" gorm:"type:string;size:255"`
	}

	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20230810",
			Migrate: func(tx *gorm.DB) error {

				if err := tx.Migrator().CreateTable(&RelayAccounts{}); err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropTable("relay_accounts"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
