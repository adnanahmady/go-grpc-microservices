package main

import (
	"log"
	"net"

	"github.com/adnanahmady/go-grpc-microservices/internal"
	"github.com/adnanahmady/go-grpc-microservices/pkg/proto"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	userService, err := internal.InitializeUserService()
	if err != nil {
		log.Fatalf("failed to initialize user service: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterUserServiceServer(s, userService.Server)

	log.Println("User service is running on port 5001")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}