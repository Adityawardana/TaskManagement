package auth

import (
	authresthandler "taskmanagement/internal/modules/auth/delivery/resthandler"
	authusecase "taskmanagement/internal/modules/auth/usecase"

	"github.com/labstack/echo/v4"
)

const moduleName = "Auth"

// Module model
type Module struct {
	restHandler *authresthandler.RestHandler
}

// NewModule module constructor
func NewModule(usecase authusecase.AuthUsecase) *Module {
	return &Module{
		restHandler: authresthandler.NewRestHandler(usecase),
	}
}

// Name returns module name
func (m *Module) Name() string {
	return moduleName
}

// RegisterRoutes registers REST routes for this module
func (m *Module) RegisterRoutes(e *echo.Echo) {
	m.restHandler.RegisterRoutes(e)
}
