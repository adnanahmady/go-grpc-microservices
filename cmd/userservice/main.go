package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"

	"github.com/adnanahmady/go-grpc-microservices/internal"
	"github.com/adnanahmady/go-grpc-microservices/pkg/proto"
	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	wg := sync.WaitGroup{}

	userService, err := internal.InitializeUserService()
	if err != nil {
		log.Fatalf("failed to initialize user service: %v", err)
	}

	cfg := userService.Config
	lgr := userService.Logger

	host := fmt.Sprintf("%s:%d", cfg.User.Host, cfg.User.Port)
	listener, err := net.Listen("tcp", host)
	if err != nil {
		lgr.Fatal("failed to listen", err)
		return
	}

	s := grpc.NewServer()
	proto.RegisterUserServiceServer(s, userService.Server)

	wg.Add(1)
	go func() {
		defer wg.Done()
		lgr.Info("User service is running on %s", host)
		if err := s.Serve(listener); err != nil {
			lgr.Error("failed to serve", err)
			return
		}
		lgr.Info("grpc server is stopped")
	}()

	<-ctx.Done()
	stop()
	s.GracefulStop()
	wg.Wait()
	lgr.Info("User service is gracefully stopped")
}
