package gorm

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

type DB struct {
	db *gorm.DB
}

func NewDB(dsn string, opts ...Option) *DB {
	s := &DB{
		db: mustNewGormDB(dsn),
	}

	// 逐一套用傳入的 Option
	for _, opt := range opts {
		opt(s)
	}

	return s
}

type Option func(db *DB)

func WithTracing() Option {
	return func(db *DB) {
		if err := db.db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
			log.Fatalf("failed to use tracing plugin: %v", err)
		}
	}
}

// mustNewGormDB 建立 GORM 資料庫連線，失敗時 panic
func mustNewGormDB(dsn string) *gorm.DB {
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
