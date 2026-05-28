package a_order

type OrderStatus string

const (
	OrderStatusCreated  OrderStatus = "created"
	OrderStatusPaid     OrderStatus = "paid"
	OrderStatusHandling OrderStatus = "handling"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusCreated, OrderStatusPaid, OrderStatusHandling:
		return true
	}
	return false
}

func NewOrderStatus(s string) (OrderStatus, error) {
	status := OrderStatus(s)
	if !status.IsValid() {
		return "", InvalidOrderStatusError
	}

	return status, nil
}

func (s *OrderStatus) Scan(src any) error {
	str, ok := src.(string)
	if !ok {
		return InvalidOrderStatusError
	}
	status, err := NewOrderStatus(str)
	if err != nil {
		return err
	}
	*s = status
	return nil
}
