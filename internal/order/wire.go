//go:build wireinject
// +build wireinject

package order

import (
	"github.com/google/wire"
)

func InitService(serviceName string) (*OrderService, error) {
	wire.Build(OrderServiceSet)
	return nil, nil
}
