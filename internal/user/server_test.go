package user

import (
	"context"
	"testing"

	"github.com/adnanahmady/go-grpc-microservices/config"
	"github.com/adnanahmady/go-grpc-microservices/pkg/applog"
	"github.com/adnanahmady/go-grpc-microservices/pkg/proto"
	"github.com/adnanahmady/go-grpc-microservices/pkg/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetUser_Unit(t *testing.T) {
	cfg := config.GetConfig()
	lgr := applog.NewAppLogger(cfg, "test_user")
	server := NewServer()
	ctx := request.WithLogger(context.Background(), lgr)

	t.Run("given id when user exists then should return it", func(t *testing.T) {
		// Arrange
		req := &proto.GetUserRequest{Id: "1"}

		// Act
		resp, err := server.GetUser(ctx, req)
		require.NoError(t, err)
		require.NotEmpty(t, resp)

		// Assert
		assert.Equal(t, "1", resp.Id)
		assert.NotEmpty(t, resp.Name)
	})

	t.Run("given id when user doesnt exist then should error", func(t *testing.T) {
		// Arrange
		req := &proto.GetUserRequest{Id: "999"}

		// Act
		resp, err := server.GetUser(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)

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
