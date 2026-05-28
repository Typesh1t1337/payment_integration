package a_order

import "errors"

var (
	InvalidQuantityError     = errors.New("invalid quantity")
	InvalidOrderStatusError  = errors.New("invalid order status")
	UnexpectedOrderError     = errors.New("unexpected order")
	OrderBeingProcessedError = errors.New("order being processed")
	ProductNotFoundError     = errors.New("product not found")
)
