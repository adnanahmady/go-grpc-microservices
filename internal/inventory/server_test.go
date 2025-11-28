package inventory

import (
	"context"
	"testing"

	"github.com/adnanahmady/go-grpc-microservices/config"
	"github.com/adnanahmady/go-grpc-microservices/pkg/applog"
	"github.com/adnanahmady/go-grpc-microservices/pkg/proto"
	"github.com/adnanahmady/go-grpc-microservices/pkg/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInventory_Unit(t *testing.T) {
	cfg := config.GetConfig()
	lgr := applog.NewAppLogger(cfg, "test_inventory")
	server := NewServer()
	ctx := request.WithLogger(context.Background(), lgr)

	t.Run("given id when product doesnt exist then should return error", func(t *testing.T) {
		// Arrange
		req := &proto.GetProductRequest{Id: "999"}

		// Act
		resp, err := server.GetProduct(ctx, req)
		require.Error(t, err)
		require.Empty(t, resp)

		// Assert
		assert.ErrorIs(t, err, ErrProductNotFound)
	})

	t.Run("given id when product exists then should return the product", func(t *testing.T) {
		// Arrange
		req := &proto.GetProductRequest{Id: "1"}

		// Act
		resp, err := server.GetProduct(ctx, req)
		require.NoError(t, err)
		require.NotEmpty(t, resp)

		// Assert
		assert.Equal(t, "1", resp.Id)
		assert.NotEmpty(t, resp.Name)
		assert.NotEmpty(t, resp.Quantity)
	})
}
