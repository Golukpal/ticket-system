package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Golukpal/ticket-system/internal/auth"
	"github.com/Golukpal/ticket-system/internal/config"
	"github.com/Golukpal/ticket-system/internal/database"
	"github.com/Golukpal/ticket-system/internal/handler"
	"github.com/Golukpal/ticket-system/internal/repository"
	"github.com/Golukpal/ticket-system/internal/routes"
	"github.com/Golukpal/ticket-system/internal/service"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	dbPool, err := database.NewPostgresPool(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer dbPool.Close()

	jwtService, err := auth.NewJWTService(cfg.JWTSecret)
	if err != nil {
		log.Fatalf("JWT setup failed: %v", err)
	}

	userRepository := repository.NewUserRepository(dbPool)

	authService := service.NewAuthService(
		userRepository,
	)

	authHandler := handler.NewAuthHandler(
		authService,
		jwtService,
	)

	router := routes.Setup(authHandler)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("server started on port %s", cfg.Port)

		if err := server.ListenAndServe(); err != nil {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)

	signal.Notify(
		shutdown,
		os.Interrupt,
		syscall.SIGTERM,
	)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}

	case sig := <-shutdown:
		log.Printf("shutdown signal received: %s", sig)

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("graceful shutdown failed: %v", err)
		}
	}
}