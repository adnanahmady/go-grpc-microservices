//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
)


func InitializeUserService() (*UserService, error) {
	wire.Build(UserServiceSet)
	return nil, nil
}
