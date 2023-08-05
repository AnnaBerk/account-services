package service

import "fmt"

var (
	ErrCannotParseToken = fmt.Errorf("cannot parse token")

	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrCannotCreateUser  = fmt.Errorf("cannot create user")

	ErrAccountNotFound = fmt.Errorf("account not found")
	ErrDeposit         = fmt.Errorf("cannot deposit")
)
