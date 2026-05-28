package handler

import (
	"payment_integration/internal/a_user/usecases/login"
	"payment_integration/internal/a_user/usecases/refresh"
	"payment_integration/internal/a_user/usecases/register"
)

type Handler struct {
	loginUseCase    login.LoginUseCase
	registerUseCase register.RegisterUseCase
	refreshUseCase  refresh.RefreshUseCase
}