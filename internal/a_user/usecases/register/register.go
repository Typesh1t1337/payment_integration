package register

import (
	"context"
	"errors"
	"payment_integration/internal/a_user"
	"payment_integration/internal/a_user/model"
	"payment_integration/internal/a_user/service"
	"payment_integration/internal/domain"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user a_user.CreateUser) (*model.User, error)
}

type RegisterUseCase struct {
	userRepository UserRepository
	jwtService     service.JwtService
}

func NewRegisterUseCase(userRepository UserRepository, jwtService service.JwtService) *RegisterUseCase {
	return &RegisterUseCase{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

type RegisterRequest struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type RegisterResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (uc *RegisterUseCase) Execute(ctx context.Context, request RegisterRequest) (*RegisterResponse, error) {
	password, err := a_user.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	createdUser, err := uc.userRepository.Create(ctx, a_user.CreateUser{
		Name: request.Name,
		Email: request.Email,
		HashedPassword: password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return nil, a_user.ErrUserAlreadyExists
		}
		return nil, err
	}

	return &RegisterResponse{
		ID: createdUser.Id.String(),
		Name: createdUser.Name,
		Email: createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}, nil
}