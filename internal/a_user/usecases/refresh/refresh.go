package refresh

import (
	"context"
	"errors"
	"payment_integration/internal/a_user"
	"payment_integration/internal/a_user/model"
	"payment_integration/internal/a_user/service"
	"payment_integration/internal/domain"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type RefreshUseCase struct {
	jwtService service.JwtService
	userRepository UserRepository
}

func NewRefreshUseCase(jwtService service.JwtService, userRepository UserRepository) *RefreshUseCase {
	return &RefreshUseCase{
		jwtService: jwtService,
		userRepository: userRepository,
	}
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

func (uc *RefreshUseCase) Execute(ctx context.Context, request *RefreshRequest) (*RefreshResponse, error) {
	claims, err := uc.jwtService.ParseToken(request.RefreshToken)
	if err != nil {
		return nil, err
	}
	if claims.Exp < time.Now().Unix() {
		return nil, a_user.ErrInvalidToken
	}
	userID, err := uuid.Parse(claims.Sub)
	if err != nil {
		return nil, err
	}
	_, err = uc.userRepository.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, a_user.ErrInvalidToken
		}
		return nil, err
	}
	accessToken, err := uc.jwtService.GenerateAccessToken(claims.Sub)
	if err != nil {
		return nil, err
	}
	return &RefreshResponse{
		AccessToken: accessToken,
	}, nil
}