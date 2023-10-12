package database

import (
	"github.com/pkg/errors"
	"gitlab.com/idolauncher/go-template-kit/models"
	"gorm.io/gorm"
)

func Migration(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.User{},
	); err != nil {
		return errors.Wrap(err, "db.AutoMigrate")
	}

	return nil
}
