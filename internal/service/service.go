package service

import (
	"account-management/internal/entity"
	"account-management/internal/lib/hasher"
	"account-management/internal/repo"
	"context"
	"golang.org/x/exp/slog"
	"time"
)

type Auth interface {
	CreateUserWithAccount(ctx context.Context, input entity.AuthCreateUserInput, log *slog.Logger) (id int64, err error)
	ParseToken(accessToken string) (int, error)
}

type Account interface {
}

type Services struct {
	Auth    Auth
	Account Account
}

type ServicesDependencies struct {
	Repos  *repo.Repositories
	Hasher hasher.PasswordHasher

	SignKey  string
	TokenTTL time.Duration
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Auth:    NewAuthService(deps.Repos.User, deps.Hasher, deps.SignKey, deps.TokenTTL),
		Account: NewAccountService(deps.Repos.Account),
	}
}
