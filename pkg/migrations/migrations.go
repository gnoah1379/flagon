package migrations

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"flagon/pkg/database"
	"fmt"
	"log/slog"

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

type Migrator struct {
	*migrate.Migrate
}

func NewMigrations(db *database.DB) (*Migrator, error) {
	driverName := db.DriverName()
	srcDriver, err := getSourceDriver(driverName)
	if err != nil {
		return nil, fmt.Errorf("error getting source driver: %w", err)
	}
	sqlDB, err := db.SqlDB()
	if err != nil {
		return nil, fmt.Errorf("error getting sql.DB: %w", err)
	}
	databaseDriver, err := getDatabaseDriver(driverName, sqlDB)
	if err != nil {
		return nil, fmt.Errorf("error getting database driver: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", srcDriver, driverName, databaseDriver)
	if err != nil {
		return nil, fmt.Errorf("error creating migrate instance: %w", err)
	}
	m.Log = &logger{}
	return &Migrator{
		Migrate: m,
	}, nil
}

func (m *Migrator) Up() error {

	err := m.Migrate.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error running migrations: %w", err)
	}
	if err != nil && errors.Is(err, migrate.ErrNoChange) {
		slog.Info("Migrations completed successfully. No changes detected")
		return nil
	}
	slog.Info("Migrations completed successfully")
	return nil
}

type logger struct {
}

func (l *logger) Printf(format string, v ...interface{}) {
	slog.Info(fmt.Sprintf(format, v...))
}

func (l *logger) Verbose() bool {
	return slog.Default().Enabled(context.Background(), slog.LevelInfo)
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
