package pgdb

import (
	"account-management/internal/entity"
	"account-management/internal/repo/repoerrs"
	"account-management/pkg/psql"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepo struct {
	db *psql.Postgres
}

func NewUserRepo(pg *psql.Postgres) *UserRepo {
	return &UserRepo{db: pg}
}

func (r *UserRepo) CreateUser(ctx context.Context, user entity.User) (id int64, err error) {

	stmtSQL := `INSERT INTO "user" (name, password) VALUES($1, $2) RETURNING id `

	err = r.db.Pool.QueryRow(ctx, stmtSQL, user.Username, user.Password).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, repoerrs.ErrAlreadyExists
			}
		}
		return 0, fmt.Errorf("UserRepo.CreateUser - r.Pool.QueryRow: %v", err)
	}

	return id, nil
}
