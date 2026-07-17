package model

import "gorm.io/gorm"

func AllModels() []any {
	return []any{
		&User{},
	}
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(AllModels()...)
}