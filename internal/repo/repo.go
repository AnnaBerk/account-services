package repo

import (
	"account-management/internal/entity"
	"account-management/internal/repo/pgdb"
	"account-management/pkg/psql"
	"context"
)

type User interface {
	CreateUserWithAccount(ctx context.Context, user entity.User) (id int64, err error)
}

type Account interface {
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
