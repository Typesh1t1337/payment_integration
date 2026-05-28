package a_order

import "errors"

var (
	ErrInvalidQuantity     = errors.New("invalid quantity")
	ErrInvalidOrderStatus  = errors.New("invalid order status")
	ErrUnexpectedOrder     = errors.New("unexpected order")
	ErrOrderBeingProcessed = errors.New("order being processed")
	ErrProductNotFound     = errors.New("product not found")
)
