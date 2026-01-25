package postgres

import (
	"time"

	"github.com/paincake00/todo-go/internal/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(uri string, maxOpenConn, maxIdleConn int, maxConnLifetime time.Duration) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(maxOpenConn)
	sqlDB.SetMaxIdleConns(maxIdleConn)
	sqlDB.SetConnMaxLifetime(maxConnLifetime)

	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
