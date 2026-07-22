package resthandler

import (
	"net/http"
	"strings"

	domain "taskmanagement/internal/modules/task/domain"
	taskusecase "taskmanagement/internal/modules/task/usecase"
	"taskmanagement/pkg/helper"
	"taskmanagement/pkg/shared"
	shareddomain "taskmanagement/pkg/shared/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RestHandler struct {
	usecase taskusecase.TaskUsecase
}

func NewRestHandler(usecase taskusecase.TaskUsecase) *RestHandler {
	return &RestHandler{usecase: usecase}
}

func (h *RestHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/tasks", h.CreateTask)
	g.GET("/tasks", h.ListTasks)
	g.GET("/tasks/:id", h.GetTask)
	g.PUT("/tasks/:id", h.UpdateTask)
	g.DELETE("/tasks/:id", h.DeleteTask)
	g.POST("/tasks/:id/assign", h.AssignTask)
}

type createTaskReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (h *RestHandler) CreateTask(c echo.Context) error {
	userID := c.Get(shared.CtxUserIDKey).(string)

	var req createTaskReq
	if err := c.Bind(&req); err != nil {
		return shareddomain.NewAppError(400, "INVALID_REQUEST", "invalid request body", err)
	}
	task, err := h.usecase.CreateTask(c.Request().Context(), userID, req.Title, req.Description, req.Status)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, task)
}

func (h *RestHandler) ListTasks(c echo.Context) error {
	userID := c.Get(shared.CtxUserIDKey).(string)
	q := domain.TaskQuery{
		Status: strings.TrimSpace(c.QueryParam("status")),
		Search: strings.TrimSpace(c.QueryParam("search")),
		Limit:  helper.AtoiDefault(c.QueryParam("limit"), 10),
		Page:   helper.AtoiDefault(c.QueryParam("page"), 1),
	}
	tasks, err := h.usecase.ListTasks(c.Request().Context(), userID, q)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"items": tasks, "limit": q.Limit, "page": q.Page})
}

func (h *RestHandler) GetTask(c echo.Context) error {
	userID := c.Get(shared.CtxUserIDKey).(string)
	task, err := h.usecase.GetTask(c.Request().Context(), userID, c.Param("id"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, task)
}

type updateTaskReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (h *RestHandler) UpdateTask(c echo.Context) error {
	userID := c.Get(shared.CtxUserIDKey).(string)
	var req updateTaskReq
	if err := c.Bind(&req); err != nil {
		return shareddomain.NewAppError(400, "INVALID_REQUEST", "invalid request body", err)
	}
	task, err := h.usecase.UpdateTask(c.Request().Context(), userID, c.Param("id"), req.Title, req.Description, req.Status)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, task)
}

func (h *RestHandler) DeleteTask(c echo.Context) error {
	userID := c.Get(shared.CtxUserIDKey).(string)
	if err := h.usecase.DeleteTask(c.Request().Context(), userID, c.Param("id")); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

type assignReq struct {
	AssigneeID string `json:"assignee_id"`
}

func (h *RestHandler) AssignTask(c echo.Context) error {
	userID := c.Get(shared.CtxUserIDKey).(string)
	var req assignReq
	if err := c.Bind(&req); err != nil {
		return shareddomain.NewAppError(400, "INVALID_REQUEST", "invalid request body", err)
	}
	if _, err := uuid.Parse(req.AssigneeID); err != nil {
		return shareddomain.NewAppError(400, "INVALID_REQUEST", "assignee_id must be valid UUID", err)
	}

	result, err := h.usecase.AssignTask(c.Request().Context(), userID, c.Param("id"), req.AssigneeID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}
