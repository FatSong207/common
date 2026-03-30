package gorm

import (
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// MustAutoMigrate 自動遷移資料表，失敗時 panic
func MustAutoMigrate(db *gorm.DB, models ...schema.Tabler) {
	if err := autoMigrate(db, models...); err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}
}

// autoMigrate 自動遷移資料表
func autoMigrate(db *gorm.DB, models ...schema.Tabler) error {
	arr := make([]any, 0)
	for _, model := range models {
		arr = append(arr, model)
	}
	return db.AutoMigrate(arr...)
}
