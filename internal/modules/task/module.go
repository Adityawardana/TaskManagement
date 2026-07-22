package task

import (
	taskresthandler "taskmanagement/internal/modules/task/delivery/resthandler"
	taskusecase "taskmanagement/internal/modules/task/usecase"

	"github.com/labstack/echo/v4"
)

const moduleName = "Task"

// Module model
type Module struct {
	restHandler *taskresthandler.RestHandler
}

// NewModule module constructor
func NewModule(usecase taskusecase.TaskUsecase) *Module {
	return &Module{
		restHandler: taskresthandler.NewRestHandler(usecase),
	}
}

// Name returns module name
func (m *Module) Name() string {
	return moduleName
}

// RegisterRoutes registers REST routes for this module on a group
func (m *Module) RegisterRoutes(g *echo.Group) {
	m.restHandler.RegisterRoutes(g)
}
