package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	userHandler "payment_integration/internal/a_user/handler"
	"payment_integration/internal/a_user/repository"
	"payment_integration/internal/a_user/service"
	"payment_integration/internal/a_user/usecases/login"
	"payment_integration/internal/a_user/usecases/refresh"
	"payment_integration/internal/a_user/usecases/register"
	"payment_integration/internal/config"
	"payment_integration/internal/transport/http_transport"
	"payment_integration/internal/transport/http_transport/middleware"
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

	// uow := uow.NewSQLUoW(pool)
	jwtService, err := service.NewJwtService(cfg.Jwt.PrivateKey, cfg.Jwt.PublicKey, cfg.Jwt.AccessTTL, cfg.Jwt.RefreshTTL)
	if err != nil {
		log.Fatal(err)
	}
	userRepository := repository.NewPostgresUserRepository(pool)
	loginUseCase := login.NewLoginUseCase(userRepository, *jwtService)
	registerUseCase := register.NewRegisterUseCase(userRepository, *jwtService)
	refreshUseCase := refresh.NewRefreshUseCase(*jwtService, userRepository)
	userHandler := userHandler.NewHandler(*loginUseCase, *registerUseCase, *refreshUseCase, logger, validator.New(), cfg)

	mux := http.NewServeMux()
	userHandler.RegisterRoutes(mux, http_transport.GetV1Prefix("auth"))

	privateMux := http.NewServeMux()
	
	mux.Handle("/api/v1/", middleware.AuthMiddleware(jwtService)(privateMux))

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
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