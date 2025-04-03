package repository

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewTokenRepository,
	NewUserRepository,
)
