package order

import (
	"context"
	"fmt"

	"github.com/adnanahmady/go-grpc-microservices/config"
	"github.com/adnanahmady/go-grpc-microservices/pkg/proto"
	"github.com/adnanahmady/go-grpc-microservices/pkg/request"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectToUserService(cfg *config.Config) (proto.UserServiceClient, error) {
	userAddr := fmt.Sprintf("%s:%d", cfg.Gateway.User.Host, cfg.Gateway.User.Port)
	opt := grpc.WithTransportCredentials(insecure.NewCredentials())
	userConn, err := grpc.NewClient(userAddr, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}
	return proto.NewUserServiceClient(userConn), nil
}

func ConnectToInventoryService(cfg *config.Config) (proto.InventoryServiceClient, error) {
	invAddr := fmt.Sprintf("%s:%d", cfg.Gateway.Inventory.Host, cfg.Gateway.Inventory.Port)
	opt := grpc.WithTransportCredentials(insecure.NewCredentials())
	invConn, err := grpc.NewClient(invAddr, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to inventory service: %w", err)
	}
	return proto.NewInventoryServiceClient(invConn), nil
}

var _ proto.OrderServiceServer = (*Server)(nil)

type Server struct {
	proto.UnimplementedOrderServiceServer
	userClient      proto.UserServiceClient
	inventoryClient proto.InventoryServiceClient
}

func NewServer(
	userClient proto.UserServiceClient,
	invClient proto.InventoryServiceClient,
) *Server {
	return &Server{
		userClient:      userClient,
		inventoryClient: invClient,
	}
}

func (s *Server) CreateOrder(
	ctx context.Context,
	req *proto.CreateOrderRequest,
) (*proto.Order, error) {
	lgr := request.GetLogger(ctx)
	_, err := s.userClient.GetUser(ctx, &proto.GetUserRequest{Id: req.UserId})
	if err != nil {
		lgr.Error("failet to find user", err)
		return nil, fmt.Errorf("%w: %s", ErrOrderingUserNotFound, req.UserId)
	}

	pr := &proto.GetProductRequest{Id: req.ProductId}
	product, err := s.inventoryClient.GetProduct(ctx, pr)
	if err != nil {
		lgr.Error("failed to find product", err)
		return nil, fmt.Errorf("%w: %s", ErrOrderingProductNotFound, req.ProductId)
	}
	if product.Quantity < 1 {
		return nil, fmt.Errorf("%w: %s", ErrProductIsSoldOut, req.ProductId)
	}

	newOrder := &proto.Order{
		Id:        uuid.NewString(),
		UserId:    req.UserId,
		ProductId: req.ProductId,
		Status:    "CREATED",
	}

	lgr.Info("Order created successfully: %v", newOrder.Id)
	return newOrder, nil
}
