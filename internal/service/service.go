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
}

type ServicesDependencies struct {
	Repos  *repo.Repositories
	Hasher hasher.PasswordHasher

	SignKey  string
	TokenTTL time.Duration
}

type Services struct {
	Auth Auth
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Auth: NewAuthService(deps.Repos.User, deps.Hasher, deps.SignKey, deps.TokenTTL),
	}
}
