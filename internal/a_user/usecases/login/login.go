package login

import (
	"context"
	"payment_integration/internal/a_user"
	"payment_integration/internal/a_user/model"
	"payment_integration/internal/a_user/service"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type LoginUseCase struct{
	userRepository UserRepository
	jwtService     service.JwtService
}

func NewLoginUseCase(userRepository UserRepository, jwtService service.JwtService) *LoginUseCase {
	return &LoginUseCase{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (uc *LoginUseCase) Execute(ctx context.Context, request LoginRequest) (*LoginResponse, error) {
	user, err := uc.userRepository.GetByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if !a_user.CheckPassword(request.Password, user.Password) {
		return nil, a_user.ErrInvalidPassword
	}
	accessToken, err := uc.jwtService.GenerateAccessToken(user.Id.String())
	if err != nil {
		return nil, err
	}
	refreshToken, err := uc.jwtService.GenerateRefreshToken(user.Id.String())
	if err != nil {
		return nil, err
	}
	return &LoginResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}, nil
}