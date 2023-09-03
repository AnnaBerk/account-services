package service

import (
	"account-management/internal/entity"
	sl "account-management/internal/lib/slog"
	"account-management/internal/repo/repoerrs"
	"context"
	"errors"
	"github.com/labstack/gommon/log"
)

func (s *AuthService) CreateUserWithAccount(ctx context.Context, input entity.AuthCreateUserInput) (id int64, err error) {
	user := entity.User{
		Username: input.Username,
		Password: s.passwordHasher.Hash(input.Password),
	}

	userId, err := s.userRepo.CreateUserWithAccount(ctx, user)
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return 0, ErrUserAlreadyExists
		}
		log.Error("AuthService.CreateUser - c.userRepo.CreateUser: %w", sl.Err(err))
		return 0, ErrCannotCreateUser
	}
	return userId, nil
}
