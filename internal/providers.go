package internal

import (
	"github.com/adnanahmady/go-grpc-microservices/config"
	"github.com/adnanahmady/go-grpc-microservices/internal/inventory"
	"github.com/adnanahmady/go-grpc-microservices/internal/user"
	"github.com/adnanahmady/go-grpc-microservices/pkg/applog"
	"github.com/adnanahmady/go-grpc-microservices/pkg/proto"
	"github.com/adnanahmady/go-grpc-microservices/pkg/request"
	"github.com/google/wire"
)

type UserService struct {
	Config      *config.Config
	Logger      applog.Logger
	Server      *user.Server
	Middlewares *request.Middlewares
}

var UserServiceSet = wire.NewSet(
	config.GetConfig,

	applog.NewAppLogger,
	wire.Bind(new(applog.Logger), new(*applog.AppLogger)),

	user.NewServer,
	wire.Bind(new(proto.UserServiceServer), new(*user.Server)),

	request.NewMiddlewares,

	wire.Struct(new(UserService), "*"),
)

type InventoryService struct {
	Config      *config.Config
	Logger      applog.Logger
	Server      *inventory.Server
	Middlewares *request.Middlewares
}

var InventoryServiceSet = wire.NewSet(
	config.GetConfig,

	applog.NewAppLogger,
	wire.Bind(new(applog.Logger), new(*applog.AppLogger)),

	inventory.NewServer,
	wire.Bind(new(proto.InventoryServiceServer), new(*inventory.Server)),

	request.NewMiddlewares,

	wire.Struct(new(InventoryService), "*"),
)

