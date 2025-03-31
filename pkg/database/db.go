package database

import (
	"context"
	"database/sql"
	"errors"
	"flagon/pkg/config"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type DB struct {
	*gorm.DB
}

func Open(cfg config.Database) (*DB, error) {
	var gormCfg = &gorm.Config{
		Logger: &slogLogger{},
	}
	var db *gorm.DB
	var err error
	switch cfg.Driver {
	case "sqlite":
		dns := fmt.Sprintf("file:%s?_fk=true", viper.GetString("db.hostname"))
		db, err = gorm.Open(sqlite.Open(dns), gormCfg)
	case "postgres":
		dns := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode)
		db, err = gorm.Open(postgres.Open(dns), gormCfg)
	case "mysql":
		dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
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
	if cfg.MaxConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxConnLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.MaxConnLifetime)
	}
	if cfg.MaxConnIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.MaxConnIdleTime)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = sqlDB.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
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
