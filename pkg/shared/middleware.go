package shared

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"taskmanagement/pkg/helper"
	"taskmanagement/pkg/shared/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	CtxUserIDKey    = "user_id"
	CtxRequestIDKey = "request_id"
)

func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rid := uuid.NewString()
		c.Set(CtxRequestIDKey, rid)
		c.Response().Header().Set("X-Request-ID", rid)
		return next(c)
	}
}

func AuthMiddleware(jwtManager *helper.JWTManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				return domain.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "missing bearer token", nil)
			}
			token := strings.TrimPrefix(h, "Bearer ")
			userID, err := jwtManager.Parse(token)
			if err != nil {
				return domain.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "invalid token", err)
			}
			c.Set(CtxUserIDKey, userID)
			return next(c)
		}
	}
}

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		if err != nil {
			c.Error(err)
		}
		status := c.Response().Status
		if status == 0 {
			status = http.StatusOK
		}
		level := "INFO"
		if status >= 500 {
			level = "ERROR"
		} else if status >= 400 {
			level = "WARN"
		}
		payload := map[string]any{
			"request_id":  c.Get(CtxRequestIDKey),
			"timestamp":   time.Now().UTC().Format(time.RFC3339Nano),
			"level":       level,
			"method":      c.Request().Method,
			"path":        c.Path(),
			"status_code": status,
			"latency_ms":  float64(time.Since(start).Microseconds()) / 1000.0,
		}
		if uid, ok := c.Get(CtxUserIDKey).(string); ok && uid != "" {
			payload["user_id"] = uid
		}
		b, _ := json.Marshal(payload)
		c.Logger().Print(string(b))
		return nil
	}
}
