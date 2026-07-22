package resthandler

import (
	"net/http"

	authusecase "taskmanagement/internal/modules/auth/usecase"
	shareddomain "taskmanagement/pkg/shared/domain"

	"github.com/labstack/echo/v4"
)

type RestHandler struct {
	usecase authusecase.AuthUsecase
}

func NewRestHandler(usecase authusecase.AuthUsecase) *RestHandler {
	return &RestHandler{usecase: usecase}
}

func (h *RestHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/auth/register", h.Register)
	e.POST("/auth/login", h.Login)
}

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TeamName string `json:"team_name"`
}

func (h *RestHandler) Register(c echo.Context) error {
	var req registerReq
	if err := c.Bind(&req); err != nil {
		return shareddomain.NewAppError(400, "INVALID_REQUEST", "invalid request body", err)
	}
	user, err := h.usecase.Register(c.Request().Context(), req.Email, req.Password, req.TeamName)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]any{"id": user.ID, "email": user.Email, "created_at": user.CreatedAt})
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *RestHandler) Login(c echo.Context) error {
	var req loginReq
	if err := c.Bind(&req); err != nil {
		return shareddomain.NewAppError(400, "INVALID_REQUEST", "invalid request body", err)
	}
	token, err := h.usecase.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"access_token": token, "token_type": "Bearer"})
}
