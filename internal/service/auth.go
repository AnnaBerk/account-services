package service

import (
	"account-management/internal/lib/hasher"
	"account-management/internal/repo"
	"github.com/golang-jwt/jwt"
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
}

func NewAuthService(userRepo repo.User, passwordHasher hasher.PasswordHasher, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		signKey:        signKey,
		tokenTTL:       tokenTTL,
	}
}

//func (s *AuthService) CreateUser(ctx context.Context, input AuthCreateUserInput) (int, error) {
//	user := entity.User{
//		Username: input.Username,
//		Password: s.passwordHasher.Hash(input.Password),
//	}
//
//	userId, err := s.userRepo.CreateUser(ctx, user)
//	if err != nil {
//		if errors.Is(err, repoerrs.ErrAlreadyExists) {
//			return 0, ErrUserAlreadyExists
//		}
//		log.Errorf("AuthService.CreateUser - c.userRepo.CreateUser: %v", err)
//		return 0, ErrCannotCreateUser
//	}
//	return userId, nil
//}
