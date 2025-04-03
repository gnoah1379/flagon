//go:build wireinject
// +build wireinject

//go:generate wire

package migrate

import (
	"flagon/internal/migrations"
	"flagon/pkg/database"

	"github.com/google/wire"
)

func New() (*CmdRunner, error) {
	wire.Build(
		wire.Struct(new(CmdRunner), "*"),
		database.Open,
		migrations.NewMigrations,
	)
	return &CmdRunner{}, nil
}
