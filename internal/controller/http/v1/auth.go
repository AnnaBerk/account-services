package v1

import (
	"account-management/internal/entity"
	"account-management/internal/service"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
	"net/http"
)

type authRoutes struct {
	authService service.Auth
	log         *slog.Logger
}

func newAuthRoutes(g *echo.Group, authService service.Auth, log *slog.Logger) {
	r := &authRoutes{
		authService: authService,
		log:         log,
	}

	g.POST("/sign-up", r.signUp)
}

type signUpInput struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
}

func (r *authRoutes) signUp(c echo.Context) error {
	var input signUpInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.authService.CreateUserWithAccount(c.Request().Context(), entity.AuthCreateUserInput{
		Username: input.Username,
		Password: input.Password,
	}, r.log)
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Id int64 `json:"id"`
	}

	return c.JSON(http.StatusCreated, response{
		Id: id,
	})
}
