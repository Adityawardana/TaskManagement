package taskmanagement

import (
	"context"
	"log"
	"net/http"
	"time"

	"taskmanagement/configs"
	"taskmanagement/internal/modules/auth"
	authrepository "taskmanagement/internal/modules/auth/repository"
	authusecase "taskmanagement/internal/modules/auth/usecase"
	"taskmanagement/internal/modules/task"
	taskrepository "taskmanagement/internal/modules/task/repository"
	taskusecase "taskmanagement/internal/modules/task/usecase"
	"taskmanagement/pkg/helper"
	"taskmanagement/pkg/shared"
	sharedrepository "taskmanagement/pkg/shared/repository"

	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

// Service model
type Service struct {
	name    string
	cfg     configs.Environment
	modules struct {
		auth *auth.Module
		task *task.Module
	}
}

// NewService constructor
func NewService(serviceName string) *Service {
	cfg := configs.LoadConfigs()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := sharedrepository.Open(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Initialize per-module repositories
	authRepo := authrepository.NewAuthRepository(db)
	taskRepo := taskrepository.NewTaskRepository(db)

	jwtTTL, err := time.ParseDuration(cfg.JWTExpiresIn)
	if err != nil {
		log.Fatalf("invalid JWT_EXPIRES_IN: %v", err)
	}

	jwtManager := helper.NewJWTManager(cfg.JWTSecret, jwtTTL)
	passwordHasher := helper.NewPasswordHasher()
	notifier := helper.NewMockNotifier()

	// Initialize usecases
	authUC := authusecase.NewAuthUsecase(authRepo, passwordHasher, jwtManager)
	taskUC := taskusecase.NewTaskUsecase(taskRepo, notifier)

	s := &Service{
		name: serviceName,
		cfg:  cfg,
	}

	// Initialize modules
	s.modules.auth = auth.NewModule(authUC)
	s.modules.task = task.NewModule(taskUC)

	return s
}

// Run starts the service
func (s *Service) Run() {
	cfg := configs.GetEnv()

	jwtTTL, _ := time.ParseDuration(cfg.JWTExpiresIn)
	jwtManager := helper.NewJWTManager(cfg.JWTSecret, jwtTTL)

	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = shared.HTTPErrorHandler
	e.Use(echoMiddleware.Recover())
	e.Use(shared.RequestIDMiddleware)
	e.Use(shared.LoggingMiddleware)

	// Register auth module routes (public)
	s.modules.auth.RegisterRoutes(e)

	// Register task module routes (authenticated)
	authenticated := e.Group("", shared.AuthMiddleware(jwtManager))
	s.modules.task.RegisterRoutes(authenticated)

	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           e,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("server running on :%s", cfg.AppPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
