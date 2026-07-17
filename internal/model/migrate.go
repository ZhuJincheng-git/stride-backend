package model

import "gorm.io/gorm"

// AllModels returns every model registered with GORM.
func AllModels() []any {
	return []any{
		&User{},
	}
}

// AutoMigrate runs `db.AutoMigrate` on every registered model.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(AllModels()...)
}