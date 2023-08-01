package v1

import (
	"account-management/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/slog"
)

func NewRouter(handler *echo.Echo, services *service.Services, log *slog.Logger) {

	handler.Use(middleware.Recover())

	auth := handler.Group("/auth")
	{
		newAuthRoutes(auth, services.Auth, log)
	}

	authMiddleware := &AuthMiddleware{services.Auth}

	v1 := handler.Group("/api/v1", authMiddleware.UserIdentity)
	{
		newAccountRoutes(v1.Group("/accounts"), services.Account)
	}

}
