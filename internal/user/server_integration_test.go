package user

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/adnanahmady/go-grpc-microservices/config"
	"github.com/adnanahmady/go-grpc-microservices/pkg/applog"
	"github.com/adnanahmady/go-grpc-microservices/pkg/proto"
	"github.com/adnanahmady/go-grpc-microservices/pkg/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func setupTestServer(t *testing.T) proto.UserServiceClient {
	listener := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() { listener.Close() })
	cfg := config.GetConfig()
	lgr := applog.NewAppLogger(cfg, "test_user")

	m := request.NewMiddlewares(lgr, cfg)
	srv := grpc.NewServer(grpc.UnaryInterceptor(m.UnaryServerLoggingInterceptor()))
	t.Cleanup(func() { srv.Stop() })

	proto.RegisterUserServiceServer(srv, NewServer())
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	ctx := request.WithLogger(context.Background(), lgr)
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	return proto.NewUserServiceClient(conn)
}

func TestUserService_Integration(t *testing.T) {
	client := setupTestServer(t)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	t.Run("given id when user exists then should return it", func(t *testing.T) {
		// Arrange
		req := &proto.GetUserRequest{Id: "2"}

		// Act
		resp, err := client.GetUser(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Assert
		assert.Equal(t, "2", resp.Id)
		assert.NotEmpty(t, resp.Name)
	})

	t.Run("given id when user doesnt exist then should return error", func(t *testing.T) {
		// Arrange
		req := &proto.GetUserRequest{Id: "999"}

		// Act
		resp, err := client.GetUser(ctx, req)

		// Assert
		assert.Nil(t, resp)
		assert.Error(t, err)
	})
}
