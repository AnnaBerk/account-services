package pgdb

import (
	"account-management/internal/entity"
	"account-management/internal/repo/repoerrs"
	"account-management/pkg/psql"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type AccountRepo struct {
	db *psql.Postgres
}

func NewAccountRepo(pg *psql.Postgres) *UserRepo {
	return &UserRepo{db: pg}
}

func (r *AccountRepo) GetAccountById(ctx context.Context, id int) (account entity.Account, err error) {

	stmtSQL := `SELECT id, balance, created
				FROM account
				JOIN "user" ON account.id = "user".account_id
				WHERE "user".id = $1`

	err = r.db.Pool.QueryRow(ctx, stmtSQL, id).Scan(
		&account.Id,
		&account.Balance,
		&account.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Account{}, repoerrs.ErrNotFound
		}
		return entity.Account{}, fmt.Errorf("AccountRepo.GetAccountById - r.Pool.QueryRow: %v", err)
	}

	return account, nil
}
