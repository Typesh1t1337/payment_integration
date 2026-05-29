package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"payment_integration/internal/a_order/handler"
	orderRepository "payment_integration/internal/a_order/repo"
	"payment_integration/internal/a_order/usecase"
	productRepository "payment_integration/internal/a_product/repo"
	userHandler "payment_integration/internal/a_user/handler"
	"payment_integration/internal/a_user/repository"
	"payment_integration/internal/a_user/service"
	"payment_integration/internal/a_user/usecases/login"
	"payment_integration/internal/a_user/usecases/refresh"
	"payment_integration/internal/a_user/usecases/register"
	"payment_integration/internal/config"
	"payment_integration/internal/transport/http_transport"
	"payment_integration/internal/transport/http_transport/middleware"
	"payment_integration/internal/uow"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("config/.env")

	cfg := config.MustLoad()
	logger := config.NewLogger(cfg.LogLevel)
	pool, err := config.NewDB(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	uow := uow.NewSQLUoW(pool)
	jwtService, err := service.NewJwtService(cfg.Jwt.PrivateKey, cfg.Jwt.PublicKey, cfg.Jwt.AccessTTL, cfg.Jwt.RefreshTTL)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewPostgresUserRepository(pool)
	orderRepo := orderRepository.NewOrderRepository(pool)
	orderItemRepo := orderRepository.NewOrderItemRepository(pool)
	productRepo := productRepository.NewProductRepository(pool)

	// usecases
	loginUseCase := login.NewLoginUseCase(userRepo, *jwtService)
	registerUseCase := register.NewRegisterUseCase(userRepo, *jwtService)
	refreshUseCase := refresh.NewRefreshUseCase(*jwtService, userRepo)

	addItemUseCase := usecase.NewAddItemUseCase(uow, orderRepo, orderItemRepo, productRepo)

	// handlers
	userAuthHandler := userHandler.NewHandler(*loginUseCase, *registerUseCase, *refreshUseCase, logger, validator.New(), cfg)
	orderHandler := handler.NewHandler(addItemUseCase, logger)

	mux := http.NewServeMux()

	userAuthHandler.RegisterRoutes(mux, http_transport.GetV1Prefix("auth"))

	privateMux := http.NewServeMux()
	orderHandler.RegisterRoutes(mux, http_transport.GetV1Prefix("orders"))

	mux.Handle("/api/v1/", middleware.AuthMiddleware(jwtService)(privateMux))
	rootHandler := middleware.Chain(mux, middleware.CloseBodyMiddleware())

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: rootHandler,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Starting server", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	logger.Info("Shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Failed to shutdown server", "error", err)
		os.Exit(1)
	}
	logger.Info("Server stopped")
}
