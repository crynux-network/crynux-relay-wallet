package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func M20260609(db *gorm.DB) *gormigrate.Gormigrate {
	type VestingRecord struct {
		Type string `gorm:"not null;size:32;default:other;index"`
	}

	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "M20260609",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Migrator().AddColumn(&VestingRecord{}, "Type"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropColumn(&VestingRecord{}, "Type"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
