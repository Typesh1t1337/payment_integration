package a_order

import (
	"encoding/json"

	"github.com/google/uuid"
)

type AddItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	quantity  int
}

func (r *AddItemRequest) Quantity() int {
	return r.quantity
}

func (r *AddItemRequest) UnmarshalJSON(data []byte) error {
	type dto struct {
		ProductID uuid.UUID `json:"product_id"`
		Quantity  int       `json:"quantity"`
	}

	var raw dto

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Quantity <= 0 {
		return InvalidQuantityError
	}

	r.ProductID = raw.ProductID
	r.quantity = raw.Quantity

	return nil
}

type AddOrderItem struct {
	ProductID uuid.UUID
	quantity  int
	OrderID   uuid.UUID
}

func (c *AddOrderItem) Quantity() int {
	return c.quantity
}

func (c *AddOrderItem) SetQuantity(quantity int) error {
	if quantity < 0 {
		return InvalidQuantityError
	}

	c.quantity = quantity
	return nil
}
