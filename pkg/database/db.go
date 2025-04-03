package database

import (
	"context"
	"database/sql"
	"errors"
	"flagon/pkg/config"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
	driverName string
}

func (db *DB) DriverName() string {
	return db.driverName
}

func Open() (*DB, error) {
	dbCfg := config.GetConfig().Database
	var gormCfg = &gorm.Config{
		Logger: &slogLogger{},
	}
	var db *gorm.DB
	var err error
	switch dbCfg.Driver {
	case "sqlite":
		dns := fmt.Sprintf("file:%s?_fk=true", dbCfg.Database)
		db, err = gorm.Open(sqlite.Open(dns), gormCfg)
	case "postgres":
		dns := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			dbCfg.Host, dbCfg.Port, dbCfg.Username, dbCfg.Password, dbCfg.Database, dbCfg.SSLMode)
		db, err = gorm.Open(postgres.Open(dns), gormCfg)
	case "mysql":
		dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Database)
		db, err = gorm.Open(mysql.Open(dns), gormCfg)
	default:
		return nil, errors.New("unsupported database driver")
	}
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if dbCfg.MaxConns > 0 {
		sqlDB.SetMaxOpenConns(dbCfg.MaxConns)
	}
	if dbCfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(dbCfg.MaxIdleConns)
	}
	if dbCfg.MaxConnLifetime > 0 {
		sqlDB.SetConnMaxLifetime(dbCfg.MaxConnLifetime)
	}
	if dbCfg.MaxConnIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(dbCfg.MaxConnIdleTime)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = sqlDB.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return &DB{
		DB:         db,
		driverName: dbCfg.Driver,
	}, nil
}

func (db *DB) SqlDB() (*sql.DB, error) {
	return db.DB.DB()

}

func (db *DB) Close() error {
	sqlDB, err := db.SqlDB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
