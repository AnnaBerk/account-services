package service

import (
	"account-management/internal/entity"
	"account-management/internal/lib/hasher"
	sl "account-management/internal/lib/slog"
	"account-management/internal/repo"
	"account-management/internal/repo/repoerrs"
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"golang.org/x/exp/slog"
	"time"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId int
}

type AuthService struct {
	userRepo       repo.User
	passwordHasher hasher.PasswordHasher
	signKey        string
	tokenTTL       time.Duration
	log            *slog.Logger
}

func NewAuthService(userRepo repo.User, passwordHasher hasher.PasswordHasher, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		signKey:        signKey,
		tokenTTL:       tokenTTL,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, input AuthCreateUserInput, log *slog.Logger) (id int64, err error) {
	user := entity.User{
		Username: input.Username,
		Password: s.passwordHasher.Hash(input.Password),
	}

	userId, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return 0, ErrUserAlreadyExists
		}
		log.Error("AuthService.CreateUser - c.userRepo.CreateUser: %v", sl.Err(err))
		return 0, ErrCannotCreateUser
	}
	return userId, nil
}
