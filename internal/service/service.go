package service

import (
	"account-management/internal/lib/hasher"
	"account-management/internal/repo"
	"context"
	"golang.org/x/exp/slog"
	"time"
)

type AuthCreateUserInput struct {
	Username string
	Password string
}

type Auth interface {
	CreateUser(ctx context.Context, input AuthCreateUserInput, log *slog.Logger) (id int64, err error)
	ParseToken(accessToken string) (int, error)
}

type ServicesDependencies struct {
	Repos  *repo.Repositories
	Hasher hasher.PasswordHasher

	SignKey  string
	TokenTTL time.Duration
}

type Services struct {
	Auth    Auth
	Account Account
}

type Account interface {
	CreateAccount(ctx context.Context) (int, error)
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Auth: NewAuthService(deps.Repos.User, deps.Hasher, deps.SignKey, deps.TokenTTL),
		//Account: NewAccountService(deps.Repos.Account),
	}
}
