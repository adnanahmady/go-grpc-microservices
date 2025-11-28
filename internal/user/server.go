package user

import (
	"context"
	"fmt"

	"github.com/adnanahmady/go-grpc-microservices/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	st := status.New(codes.NotFound, ErrUserNotFound.Error())
	detail := &proto.ErrorDetail{
		ErrorCode: "USER_NOT_FOUND",
		Message:   fmt.Sprintf("%s: %s", ErrUserNotFound, req.Id),
	}
	detailedSt, err := st.WithDetails(detail)
	if err != nil {
		return nil, st.Err()
	}
	return nil, detailedSt.Err()
}
