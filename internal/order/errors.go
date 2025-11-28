package order

import "errors"

var (
	ErrOrderingUserNotFound    = errors.New("ordering user not found")
	ErrOrderingInvalidUser     = errors.New("user id is invalid")
	ErrOrderingProductNotFound = errors.New("ordering product not found")
	ErrOrderingInvalidProduct  = errors.New("product id is invalid")
	ErrProductIsSoldOut        = errors.New("product is sold out")
)
