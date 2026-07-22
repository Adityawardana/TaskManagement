package shared

import (
	"errors"
	"net/http"
	"time"

	"taskmanagement/pkg/shared/domain"

	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		_ = c.JSON(appErr.Status, domain.ErrorResponse{
			Status:    appErr.Status,
			Code:      appErr.Code,
			Message:   appErr.Message,
			Timestamp: time.Now().UTC(),
		})
		return
	}

	_ = c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
		Status:    http.StatusInternalServerError,
		Code:      "INTERNAL_ERROR",
		Message:   "internal server error",
		Timestamp: time.Now().UTC(),
	})
}
