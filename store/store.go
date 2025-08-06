package store

import (
	"fmt"
	"log"
	"os"

	"github.com/Bronku/iroon/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func LoadStore(dsn string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			IgnoreRecordNotFoundError: true},
	)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return db, fmt.Errorf("error connecting to database %w", err)
	}
	err = db.AutoMigrate(&models.Order{}, &models.OrderItem{}, &models.Product{})
	if err != nil {
		return db, fmt.Errorf("error migrating database %w", err)
	}
	return db, nil
}
