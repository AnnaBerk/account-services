package repo

import (
	"account-management/internal/entity"
	"account-management/internal/repo/pgdb"
	"account-management/pkg/psql"
	"context"
)

type User interface {
	CreateUser(ctx context.Context, user entity.User) (id int64, err error)
}

type Account interface {
	CreateAccount(ctx context.Context) (int, error)
}

type Repositories struct {
	User
	Account
}

func NewRepositories(pg *psql.Postgres) *Repositories {
	return &Repositories{
		User:    pgdb.NewUserRepo(pg),
		Account: pgdb.NewAccountRepo(pg),
	}
}
