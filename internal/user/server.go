package user

import (
	"context"
	"fmt"

	"github.com/adnanahmady/go-grpc-microservices/pkg/proto"
)

var users = map[string]*proto.User{
	"1": {
		Id:   "1",
		Name: "John Doe",
	},
	"2": {
		Id:   "2",
		Name: "Jane Doe",
	},
}

var _ proto.UserServiceServer = (*Server)(nil)

type Server struct {
	proto.UnimplementedUserServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetUser(
	ctx context.Context,
	req *proto.GetUserRequest,
) (*proto.User, error) {
	if user, ok := users[req.Id]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("user not found: %s", req.Id)
}