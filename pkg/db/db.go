package db

import (
	"context"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/thejasn/tester/pkg/config"
	"github.com/thejasn/tester/pkg/log"
)

func createConnection(ctx context.Context, url string, maxIdleConnections, maxOpenConnections int, maxConnectionLifetime time.Duration) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(maxIdleConnections)
	sqlDB.SetMaxOpenConns(maxOpenConnections)
	sqlDB.SetConnMaxLifetime(maxConnectionLifetime)
	return db, nil
}

// LoadDatabase initializes the Database struct with connections
func LoadDatabase(ctx context.Context, conf config.AppConfig) *gorm.DB {
	db, err := createConnection(
		ctx,
		conf.Database.URL,
		conf.Database.MaxIdleConnections,
		conf.Database.MaxOpenConnections,
		conf.Database.MaxConnectionLifetime)
	if err != nil {
		log.GetLogger(ctx).Fatalf("error during establishing master data source connection %v", err.Error())
	}
	return db
}
