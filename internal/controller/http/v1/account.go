package v1

import (
	"account-management/internal/service"
	"github.com/labstack/echo/v4"
)

type accountRoutes struct {
	accountService service.Account
}

func newAccountRoutes(g *echo.Group, accountService service.Account) {
	r := &accountRoutes{
		accountService: accountService,
	}
	_ = r

}
