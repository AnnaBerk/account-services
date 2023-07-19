package pgdb

import (
	"account-management/pkg/psql"
	"context"
	"fmt"
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
