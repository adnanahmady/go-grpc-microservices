package inventory

import (
	"context"
	"fmt"

	"github.com/adnanahmady/go-grpc-microservices/pkg/proto"
)

var products = map[string]*proto.Product{
	"1": {
		Id:       "1",
		Name:     "Product 1",
		Quantity: 10,
	},
	"2": {
		Id:       "2",
		Name:     "Product 2",
		Quantity: 5,
	},
	"3": {
		Id:       "3",
		Name:     "Product 3",
		Quantity: 1,
	},
}

var _ proto.InventoryServiceServer = (*Server)(nil)

type Server struct {
	proto.UnimplementedInventoryServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetProduct(
	ctx context.Context,
	req *proto.GetProductRequest,
) (*proto.Product, error) {
	if p, ok := products[req.Id]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("%w: %v", ErrProductNotFound, req.Id)
}
