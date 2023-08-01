package service

import (
	"account-management/internal/repo"
)

type AccountService struct {
	accountRepo repo.Account
}

func NewAccountService(accountRepo repo.Account) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}
