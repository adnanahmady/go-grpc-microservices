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

	"github.com/adnanahmady/go-grpc-microservices/internal/inventory"
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

	inventoryService, err := inventory.InitService("inventory")
	if err != nil {
		log.Fatalf("failed to initialize inventory service: %v", err)
	}

	cfg := inventoryService.Config
	lgr := inventoryService.Logger

	host := fmt.Sprintf("%s:%d", cfg.Inventory.Host, cfg.Inventory.Port)
	listener, err := net.Listen("tcp", host)
	if err != nil {
		lgr.Fatal("failed to create listener on %s", err, host)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(inventoryService.Middlewares.UnaryServerLoggingInterceptor()),
	)
	proto.RegisterInventoryServiceServer(s, inventoryService.Server)

	wg.Add(1)
	go func() {
		defer wg.Done()
		lgr.Info("Inventory service is running on %s", host)
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
	lgr.Info("Inventory service is gracefully stopped")
}
