package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/ZhuJincheng-git/stride-backend/internal/config"
)

func Open(cfg *config.Config) (*gorm.DB, error) {
	logLevel := gormlogger.Warn
	if cfg.LogLevel == "debug" {
		logLevel = gormlogger.Info
	}

	gormLogger := gormlogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gormlogger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logLevel,
		IgnoreRecordNotFoundError: true,
		Colorful:                  !cfg.IsReleaseMode(),
	})

	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN()), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: false,
		NowFunc:                                  func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		return nil, fmt.Errorf("connect mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("acquire underlying *sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.DBMaxIdle)
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpen)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime())

	return db, nil
}
