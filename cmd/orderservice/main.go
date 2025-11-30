package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/adnanahmady/go-grpc-microservices/internal/order"
	"github.com/adnanahmady/go-grpc-microservices/pkg/proto"
	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM,
	)
	defer stop()
	wg := sync.WaitGroup{}

	orderService, err := order.InitService("order")
	if err != nil {
		log.Fatalf("failed to initilize order service: %v", err)
	}

	cfg := orderService.Config
	lgr := orderService.Logger

	host := fmt.Sprintf("%s:%d", cfg.Order.Host, cfg.Order.Port)
	listener, err := net.Listen("tcp", host)
	if err != nil {
		lgr.Fatal("failed to create listener on %s", err, host)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(orderService.Middlewares.UnaryServerLoggingInterceptor()),
	)
	proto.RegisterOrderServiceServer(s, orderService.Server)

	wg.Add(1)
	go func() {
		defer wg.Done()
		lgr.Info("Order service is running on %s", host)
		if err := s.Serve(listener); err != nil {
			lgr.Error("failed to serve", err)
			os.Exit(1)
		}
		lgr.Info("grpc server is stopped")
	}()

	<-ctx.Done()
	stop()
	s.GracefulStop()
	wg.Wait()
	lgr.Info("Order service is gracefully stopped")
}
