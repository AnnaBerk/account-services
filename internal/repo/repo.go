package repo

import (
	"account-management/internal/entity"
	"account-management/internal/repo/pgdb"
	"account-management/pkg/psql"
	"context"
	"log"
)

type User interface {
	CreateUser(ctx context.Context, user entity.User) (id int64, err error)
}

type Repositories struct {
	User
}

func NewRepositories(pg *psql.Postgres) *Repositories {
	userRepo, err := pgdb.NewUserRepo(pg)
	if err != nil {
		// здесь обработка ошибки, в зависимости от вашей логики приложения
		log.Fatalf("Failed to create user repository: %s", err.Error())
	}
	return &Repositories{
		User: userRepo,
	}
}
