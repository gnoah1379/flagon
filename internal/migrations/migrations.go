package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"flagon/pkg/database"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	migrateDb "github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed postgres/*.sql
var postgresFs embed.FS

//go:embed sqlite/*.sql
var sqliteFs embed.FS

func Up(driverName string, instance *database.DB) error {
	srcDriver, err := getSourceDriver(driverName)
	if err != nil {
		return fmt.Errorf("error getting source driver: %w", err)
	}
	sqlDB, err := instance.SqlDB()
	if err != nil {
		return fmt.Errorf("error getting sql.DB: %w", err)
	}
	databaseDriver, err := getDatabaseDriver(driverName, sqlDB)
	if err != nil {
		return fmt.Errorf("error getting data source name: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", srcDriver, driverName, databaseDriver)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %w", err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error running migration: %w", err)
	}
	return nil
}

func getSourceDriver(driverName string) (source.Driver, error) {
	switch driverName {
	case "postgres":
		return iofs.New(postgresFs, "postgres")
	case "sqlite":
		return iofs.New(sqliteFs, "sqlite")
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driverName)
	}
}

func getDatabaseDriver(driverName string, instance *sql.DB) (migrateDb.Driver, error) {
	switch driverName {
	case "postgres":
		return postgres.WithInstance(instance, &postgres.Config{})
	case "sqlite":
		return sqlite.WithInstance(instance, &sqlite.Config{})
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driverName)
	}
}
