package uow

import "context"

type TxKey struct{}

func Do[T any](ctx context.Context, u UoW, fn func(ctx context.Context) (T, error)) (T, error) {
	result, err := u.Do(ctx, func(ctx context.Context) (any, error) {
		return fn(ctx)
	})
	if err != nil {
		var zero T
		return zero, err
	}

	return result.(T), nil
}
