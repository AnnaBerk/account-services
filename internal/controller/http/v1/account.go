package v1

import (
	"account-management/internal/entity"
	"account-management/internal/service"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
	"net/http"
)

type accountRoutes struct {
	accountService service.Account
	log            *slog.Logger
}

func newAccountRoutes(g *echo.Group, accountService service.Account) {
	r := &accountRoutes{
		accountService: accountService,
	}
	g.POST("/deposit", r.deposit)

}

type accountDepositInput struct {
	Id     int `json:"id" validate:"required"`
	Amount int `json:"amount" validate:"required,gt=0"`
}

func (r *accountRoutes) deposit(c echo.Context) error {
	var input accountDepositInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.accountService.Deposit(c.Request().Context(), entity.AccountDepositInput{
		Id:     input.Id,
		Amount: input.Amount,
	}, r.log)
	if err != nil {
		if err == service.ErrAccountNotFound {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}
