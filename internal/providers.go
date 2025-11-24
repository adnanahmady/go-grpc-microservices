package internal

import (
	"github.com/adnanahmady/go-grpc-microservices/internal/user"
	"github.com/adnanahmady/go-grpc-microservices/pkg/proto"
	"github.com/google/wire"
)

type UserService struct {
	Server *user.Server
}

var UserServiceSet = wire.NewSet(
	user.NewServer,
	wire.Bind(new(proto.UserServiceServer), new(*user.Server)),

	wire.Struct(new(UserService), "*"),
)