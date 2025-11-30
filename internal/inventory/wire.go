//go:build wireinject
// +build wireinject

package inventory

import (
	"github.com/google/wire"
)

func InitService(serviceName string) (*InventoryService, error) {
	wire.Build(InventoryServiceSet)
	return nil, nil
}
