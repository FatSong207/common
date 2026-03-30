package gorm

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MustNewGormDB 建立 GORM 資料庫連線，失敗時 panic
func MustNewGormDB(dsn string) *gorm.DB {
	db, err := newGormDB(dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return db
}

// newGormDB 建立 GORM 資料庫連線
func newGormDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	return db, nil
}
