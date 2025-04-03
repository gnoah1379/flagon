//go:build wireinject
// +build wireinject

//go:generate wire
package server

import (
	v1 "flagon/pkg/api/v1"
	"flagon/pkg/cache"
	"flagon/pkg/database"
	"flagon/pkg/repository"
	"flagon/pkg/server"
	"flagon/pkg/service"
	"github.com/google/wire"
)

func New() (*CmdRunner, error) {
	wire.Build(
		wire.Struct(new(CmdRunner), "*"),
		server.NewHttpServer,
		v1.WireSet,
		repository.WireSet,
		service.WireSet,
		database.Open,
		cache.New,
	)
	return &CmdRunner{}, nil
}
