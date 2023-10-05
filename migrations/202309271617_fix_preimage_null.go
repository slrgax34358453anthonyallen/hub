package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Update payments with preimage as an empty string to use NULL instead
var _202309271617_fix_preimage_null = &gormigrate.Migration {
	ID: "202309271617_fix_preimage_null",
	Migrate: func(tx *gorm.DB) error {
		return tx.Table("payments").Where("preimage = ?", "").Update("preimage", nil).Error;
	},
	Rollback: func(tx *gorm.DB) error {
		return nil;
	},
}