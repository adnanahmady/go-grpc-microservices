package user

import (
	"github.com/adnanahmady/go-grpc-microservices/config"
	"github.com/adnanahmady/go-grpc-microservices/pkg/applog"
	"github.com/adnanahmady/go-grpc-microservices/pkg/proto"
	"github.com/adnanahmady/go-grpc-microservices/pkg/request"
	"github.com/google/wire"
)

type UserService struct {
	Config      *config.Config
	Logger      applog.Logger
	Server      *Server
	Middlewares *request.Middlewares
}

var UserServiceSet = wire.NewSet(
	config.GetConfig,

	applog.NewAppLogger,
	wire.Bind(new(applog.Logger), new(*applog.AppLogger)),

	NewServer,
	wire.Bind(new(proto.UserServiceServer), new(*Server)),

	request.NewMiddlewares,

	wire.Struct(new(UserService), "*"),
)
