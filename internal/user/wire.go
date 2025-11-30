//go:build wireinject
// +build wireinject

package user

import (
	"github.com/google/wire"
)

func InitService(serviceName string) (*UserService, error) {
	wire.Build(UserServiceSet)
	return nil, nil
}
