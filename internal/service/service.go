package service

import (
	"account-management/internal/lib/hasher"
	"account-management/internal/repo"
	"time"
)

type AuthCreateUserInput struct {
	Username string
	Password string
}

type ServicesDependencies struct {
	Repos  *repo.Repositories
	Hasher hasher.PasswordHasher

	SignKey  string
	TokenTTL time.Duration
}
