package order

import (
	"context"
	"fmt"
	"testing"

	"github.com/adnanahmady/go-grpc-microservices/config"
	"github.com/adnanahmady/go-grpc-microservices/pkg/applog"
	"github.com/adnanahmady/go-grpc-microservices/pkg/proto"
	"github.com/adnanahmady/go-grpc-microservices/pkg/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var _ proto.UserServiceClient = (*spyUserClient)(nil)

type spyUserClient struct {
}

// GetUser implements proto.UserServiceClient.
func (s *spyUserClient) GetUser(
	ctx context.Context,
	req *proto.GetUserRequest,
	opts ...grpc.CallOption,
) (*proto.User, error) {
	if req.Id == "100" {
		return nil, fmt.Errorf("user not found")
	}
	return &proto.User{Id: req.Id, Name: "John"}, nil
}

var _ proto.InventoryServiceClient = (*spyInventoryClient)(nil)

type spyInventoryClient struct {
}

// GetProduct implements proto.InventoryServiceClient.
func (s *spyInventoryClient) GetProduct(
	ctx context.Context,
	req *proto.GetProductRequest,
	opts ...grpc.CallOption,
) (*proto.Product, error) {
	if req.Id == "100" {
		return nil, fmt.Errorf("product not found")
	}
	return &proto.Product{
		Id:       req.Id,
		Name:     "Product 10",
		Quantity: 10,
	}, nil
}

func TestCreateOrder_Unit(t *testing.T) {
	cfg := config.GetConfig()
	lgr := applog.NewAppLogger(cfg, "test_order")
	userClient := &spyUserClient{}
	invClient := &spyInventoryClient{}
	server := NewServer(userClient, invClient)
	ctx := context.Background()
	ctx = request.WithLogger(ctx, lgr)

	t.Run("given order when user doesnt exist then should return error", func(t *testing.T) {
		// Arrange
		req := &proto.CreateOrderRequest{UserId: "100", ProductId: "1"}

		// Act
		resp, err := server.CreateOrder(ctx, req)
		require.Error(t, err)
		require.Empty(t, resp)

		// Assert
		assert.ErrorIs(t, err, ErrOrderingUserNotFound)
	})

	t.Run("given order when product doesnt exist then should return error", func(t *testing.T) {
		// Arrange
		req := &proto.CreateOrderRequest{UserId: "1", ProductId: "100"}

		// Act
		resp, err := server.CreateOrder(ctx, req)
		require.Error(t, err)
		require.Empty(t, resp)

		// Assert
		assert.ErrorIs(t, err, ErrOrderingProductNotFound)
	})

	t.Run("given order when created then should return the order", func(t *testing.T) {
		// Arrange
		req := &proto.CreateOrderRequest{UserId: "3", ProductId: "2"}

		// Act
		resp, err := server.CreateOrder(ctx, req)
		require.NoError(t, err)
		require.NotEmpty(t, resp)

		// Assert
		assert.NotEmpty(t, resp.Id)
		assert.Equal(t, "3", resp.UserId)
		assert.Equal(t, "2", resp.ProductId)
		assert.Equal(t, "CREATED", resp.Status)
	})
}
