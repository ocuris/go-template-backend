package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ocuris/go-template-backend/internals/config"
	"github.com/ocuris/go-template-backend/internals/pkg/database/postgres"
	CustomMiddleware "github.com/ocuris/go-template-backend/internals/pkg/middleware"
	"github.com/ocuris/go-template-backend/internals/utils"
	"github.com/ocuris/go-template-backend/internals/utils/ctx"
	"github.com/ocuris/go-template-backend/internals/utils/logger"

	HealthzHandler "github.com/ocuris/go-template-backend/internals/modules/healthz/delivery/https"
	HealthzRepository "github.com/ocuris/go-template-backend/internals/modules/healthz/repository"
	HealthzUsecase "github.com/ocuris/go-template-backend/internals/modules/healthz/usecase"
)

var log = logger.NewLogger("workflow_engine")

// Init initializes and starts the application server.
func Init() {
	e := echo.New()

	// Load configuration
	cnf := config.NewImmutableConfigs()

	// Initialize PostgreSQL client
	postgresClient := postgres.NewPostgressClient(cnf)
	db, err := postgresClient.InitClient(context.Background())
	if err != nil {
		log.Panicf("Failed to initialize database: %s ", err.Error())
	}

	// use requestID middleware
	e.Use(CustomMiddleware.MiddlewareRequestID())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())

	// Middleware to inject dependencies into the request context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			customCtx := &ctx.CustomApplicationContext{
				Context:    c,
				PostgresDB: db,
			}
			return next(customCtx)
		}
	})
	e.Use(middleware.CORS())

	// Set validator globally
	validator := utils.DefaultValidator()
	e.Validator = validator

	// Initialize Health module
	healthzRepo := HealthzRepository.NewHealthzRepository(db)
	healthzUsecase := HealthzUsecase.NewHealthzUsecase(healthzRepo)
	HealthzHandler.NewHealthzHandler(e, healthzUsecase)

	// Start server in a separate goroutine
	serverAddr := fmt.Sprintf(":%d", cnf.GetPort())
	go func() {
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server shutdown unexpectedly: %v", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	}
	log.Info("Server exited properly.")
}
