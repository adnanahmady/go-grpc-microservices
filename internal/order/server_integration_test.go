package order

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/adnanahmady/go-grpc-microservices/config"
	"github.com/adnanahmady/go-grpc-microservices/pkg/applog"
	"github.com/adnanahmady/go-grpc-microservices/pkg/proto"
	"github.com/adnanahmady/go-grpc-microservices/pkg/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func setupOrderServer(t *testing.T) (context.Context, proto.OrderServiceClient) {
	listener := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() { listener.Close() })

	cfg := config.GetConfig()
	lgr := applog.NewAppLogger(cfg, "test_order")
	m := request.NewMiddlewares(lgr, cfg)
	srv := grpc.NewServer(grpc.UnaryInterceptor(m.UnaryServerLoggingInterceptor()))
	t.Cleanup(func() { srv.Stop() })

	proto.RegisterOrderServiceServer(srv, NewServer(&spyUserClient{}, &spyInventoryClient{}))
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("Server existed with error: %v", err)
		}
	}()

	ctx := request.WithLogger(context.Background(), lgr)
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	return ctx, proto.NewOrderServiceClient(conn)
}

func TestOrderService_Integration(t *testing.T) {
	ctx, client := setupOrderServer(t)

	t.Run("given order when user doesnt exist then should return error", func(t *testing.T) {
		// Arrange
		req := &proto.CreateOrderRequest{UserId: "100", ProductId: "1"}

		// Act
		resp, err := client.CreateOrder(ctx, req)
		require.Error(t, err)
		require.Empty(t, resp)

		// Assert
		st, ok := status.FromError(err)
		require.Truef(t, ok, "error should be a gRPC status error")
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, "user id is invalid", st.Message())
		detail, ok := st.Details()[0].(*proto.ErrorDetail)
		require.Truef(t, ok, "status should have one error detail")
		assert.Contains(t, "ordering user not found: 100", detail.Message)
	})

	t.Run("given order when product doesnt exist then should return error", func(t *testing.T) {
		// Arrange
		req := &proto.CreateOrderRequest{UserId: "1", ProductId: "100"}

		// Act
		resp, err := client.CreateOrder(ctx, req)
		require.Error(t, err)
		require.Empty(t, resp)

		// Assert
		st, ok := status.FromError(err)
		require.Truef(t, ok, "error should be a gRPC status error")
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Contains(t, "product id is invalid", st.Message())
		detail, ok := st.Details()[0].(*proto.ErrorDetail)
		require.Truef(t, ok, "status should have one error detail")
		assert.Contains(t, "ordering product not found: 100", detail.Message)
	})

	t.Run("given order when created then should return the order", func(t *testing.T) {
		// Arrange
		req := &proto.CreateOrderRequest{UserId: "3", ProductId: "2"}

		// Act
		resp, err := client.CreateOrder(ctx, req)
		require.NoError(t, err)
		require.NotEmpty(t, resp)

		// Assert
		assert.NotEmpty(t, resp.Id)
		assert.Equal(t, "3", resp.UserId)
		assert.Equal(t, "2", resp.ProductId)
		assert.Equal(t, "CREATED", resp.Status)
	})
}
