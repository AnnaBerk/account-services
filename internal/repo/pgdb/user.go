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

func NewUserRepo(pg *psql.Postgres) (*UserRepo, error) {
	const op = "pgdb.NewUserRepo"

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS "user" (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created  date NOT NULL DEFAULT CURRENT_DATE
		);
	`
	_, err := pg.Pool.Exec(context.Background(), createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &UserRepo{db: pg}, nil
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
