package service

import (
	"account-management/internal/entity"
	"account-management/internal/lib/hasher"
	sl "account-management/internal/lib/slog"
	"account-management/internal/repo"
	"account-management/internal/repo/repoerrs"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/gommon/log"
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

func NewAuthService(userRepo repo.User, passwordHasher hasher.PasswordHasher, signKey string, tokenTTL time.Duration, log *slog.Logger) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		signKey:        signKey,
		tokenTTL:       tokenTTL,
		log:            log,
	}
}

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

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})

	if err != nil {
		return 0, ErrCannotParseToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, ErrCannotParseToken
	}

	return claims.UserId, nil
}
