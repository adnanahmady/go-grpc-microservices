package user

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

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

func setupTestServer(t *testing.T) proto.UserServiceClient {
	listener := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() { listener.Close() })

	ps, err := InitService("test_user")
	require.NoError(t, err)

	logInterceptor := ps.Middlewares.UnaryServerLoggingInterceptor()
	srv := grpc.NewServer(grpc.UnaryInterceptor(logInterceptor))
	t.Cleanup(func() { srv.Stop() })

	proto.RegisterUserServiceServer(srv, NewServer())
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	ctx := request.WithLogger(context.Background(), ps.Logger)
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
		require.Empty(t, resp)

		// Assert
		st, ok := status.FromError(err)
		require.Truef(t, ok, "error should be a gRPC status error")
		assert.Equal(t, codes.NotFound, st.Code())
		assert.Equal(t, st.Message(), ErrUserNotFound.Error())
		require.Len(t, st.Details(), 1, "status should have one error detail")
		detail, ok := st.Details()[0].(*proto.ErrorDetail)
		require.Truef(t, ok, "detail should be of type ErrorDetail")
		assert.Equal(t, "USER_NOT_FOUND", detail.ErrorCode)
	})
}
