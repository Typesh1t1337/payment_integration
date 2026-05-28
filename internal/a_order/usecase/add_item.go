package usecase

import (
	"context"
	"payment_integration/internal/a_order"
	"payment_integration/internal/a_order/model"
	"payment_integration/internal/uow"

	"github.com/google/uuid"
)

type OrderRepository interface {
	GetOrCreate(ctx context.Context, userId uuid.UUID) (*model.Order, error)
}

type OrderItemRepository interface {
	Create(ctx context.Context, value a_order.AddOrderItem) (*model.OrderItems, error)
}

type ProductRepository interface {
	Exists(ctx context.Context, productID uuid.UUID) bool
}

type AddItemUseCase struct {
	uow           uow.UoW
	repo          OrderRepository
	orderItemRepo OrderItemRepository
	productRepo   ProductRepository
}

func NewAddItemUseCase(uow uow.UoW, repo OrderRepository, orderItemRepo OrderItemRepository, productRepo ProductRepository) AddItemUseCase {
	return AddItemUseCase{uow: uow, repo: repo, orderItemRepo: orderItemRepo, productRepo: productRepo}
}

func (uc *AddItemUseCase) Execute(ctx context.Context, userID uuid.UUID, body a_order.AddItemRequest) error {
	_, err := uow.Do(ctx, uc.uow, func(ctx context.Context) (*struct{}, error) {
		isProductExists := uc.productRepo.Exists(ctx, body.ProductID)
		if !isProductExists {
			return nil, a_order.ProductNotFoundError
		}

		order, err := uc.repo.GetOrCreate(ctx, userID)
		if err != nil {
			return nil, err
		}

		createOrderItemDTO := a_order.AddOrderItem{
			ProductID: body.ProductID,
			OrderID:   order.ID,
		}

		err = createOrderItemDTO.SetQuantity(body.Quantity())

		if err != nil {
			return nil, err
		}

		_, err = uc.orderItemRepo.Create(ctx, createOrderItemDTO)
		if err != nil {
			return nil, err
		}

		return &struct{}{}, nil
	})

	if err != nil {
		return err
	}

	return nil
}
