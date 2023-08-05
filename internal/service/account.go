package service

import (
	"account-management/internal/entity"
	sl "account-management/internal/lib/slog"
	"account-management/internal/repo"
	"context"
	"golang.org/x/exp/slog"
)

type AccountService struct {
	accountRepo repo.Account
	log         *slog.Logger
}

func NewAccountService(accountRepo repo.Account) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

func (s *AccountService) Deposit(ctx context.Context, input entity.AccountDepositInput, log *slog.Logger) error {
	userID := input.Id
	account, err := s.accountRepo.GetAccountById(ctx, userID)
	if err != nil {
		log.Error("AccountService.GetAccountById - GetAccountById: %w", sl.Err(err))
		return ErrAccountNotFound
	}

	err = s.accountRepo.Deposit(ctx, account.Id, input.Amount)
	if err != nil {
		log.Error("AccountService.Deposit - GetAccountById: %w", sl.Err(err))
		return ErrDeposit
	}
	return nil
}
