package order

import "errors"

var (
	ErrOrderingUserNotFound    = errors.New("ordering user not found")
	ErrOrderingProductNotFound = errors.New("ordering product not found")
)
