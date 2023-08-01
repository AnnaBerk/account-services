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

func (r *UserRepo) CreateUserWithAccount(ctx context.Context, user entity.User) (id int64, err error) {
	// Start transaction
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("UserRepo.CreateUser - Begin: %v", err)
	}

	// Create account
	stmtSQL := `INSERT INTO account (balance) VALUES(0) RETURNING id `
	var accountId int64
	err = tx.QueryRow(ctx, stmtSQL).Scan(&accountId)
	if err != nil {
		_ = tx.Rollback(ctx)
		return 0, fmt.Errorf("UserRepo.CreateUser - tx.QueryRow: %v", err)
	}

	// Create user with reference to account
	stmtSQL = `INSERT INTO "user" (name, password, account_id) VALUES($1, $2, $3) RETURNING id `
	err = tx.QueryRow(ctx, stmtSQL, user.Username, user.Password, accountId).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				_ = tx.Rollback(ctx)
				return 0, repoerrs.ErrAlreadyExists
			}
		}
		_ = tx.Rollback(ctx)
		return 0, fmt.Errorf("UserRepo.CreateUser - tx.QueryRow: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("UserRepo.CreateUser - Commit: %v", err)
	}

	return id, nil
}
