package order

import (
	"github.com/adnanahmady/go-grpc-microservices/config"
	"github.com/adnanahmady/go-grpc-microservices/pkg/applog"
	"github.com/adnanahmady/go-grpc-microservices/pkg/proto"
	"github.com/adnanahmady/go-grpc-microservices/pkg/request"
	"github.com/google/wire"
)

type OrderService struct {
	Config      *config.Config
	Logger      applog.Logger
	Server      *Server
	Middlewares *request.Middlewares
}

var OrderServiceSet = wire.NewSet(
	config.GetConfig,

	applog.NewAppLogger,
	wire.Bind(new(applog.Logger), new(*applog.AppLogger)),

	ConnectToUserService,
	ConnectToInventoryService,
	NewServer,
	wire.Bind(new(proto.OrderServiceServer), new(*Server)),

	request.NewMiddlewares,

	wire.Struct(new(OrderService), "*"),
)
