package service

import (
	"account-management/internal/entity"
	sl "account-management/internal/lib/slog"
	"account-management/internal/repo"
	"context"
	"github.com/labstack/gommon/log"
	"golang.org/x/exp/slog"
)

type AccountService struct {
	accountRepo repo.Account
	log         *slog.Logger
}

func NewAccountService(accountRepo repo.Account, log *slog.Logger) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		log:         log,
	}
}

func (s *AccountService) Deposit(ctx context.Context, input entity.AccountDepositInput) error {
	userID := input.Id
	account, err := s.accountRepo.GetAccountById(ctx, userID)
	if err != nil {
		log.Error("AccountService.GetAccountById - GetAccountById: %w", sl.Err(err))
		return err
	}

	err = s.accountRepo.Deposit(ctx, account.Id, input.Amount)
	if err != nil {
		log.Error("AccountService.Deposit - GetAccountById: %w", sl.Err(err))
		return ErrDeposit
	}
	return nil
}
