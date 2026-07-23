package database

import (
	"fmt"
	"sync/atomic"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/ZhuJincheng-git/stride-backend/internal/model"
)

var sqliteSeq uint64

func OpenSQLite() (*gorm.DB, error) {
	id := atomic.AddUint64(&sqliteSeq, 1)
	dsn := fmt.Sprintf("file:stride_test_%d?mode=memory&cache=shared&_foreign_keys=on", id)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return nil, err
	}
	if err := model.AutoMigrate(db); err != nil {
		return nil, err
	}
	return db, nil
}