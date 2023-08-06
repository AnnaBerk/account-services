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
	CreateUserWithAccount(ctx context.Context, input entity.AuthCreateUserInput) (id int64, err error)
	ParseToken(accessToken string) (int, error)
}

type Account interface {
	Deposit(ctx context.Context, input entity.AccountDepositInput) error
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

func NewServices(deps ServicesDependencies, log *slog.Logger) *Services {
	return &Services{
		Auth:    NewAuthService(deps.Repos.User, deps.Hasher, deps.SignKey, deps.TokenTTL, log),
		Account: NewAccountService(deps.Repos.Account, log),
	}
}
