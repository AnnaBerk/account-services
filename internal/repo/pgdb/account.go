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

func NewAccountRepo(pg *psql.Postgres) *AccountRepo {
	return &AccountRepo{db: pg}
}

func (r *AccountRepo) GetAccountById(ctx context.Context, id int) (account entity.Account, err error) {

	stmtSQL := `SELECT id, account.balance, account.created
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
		return entity.Account{}, fmt.Errorf("AccountRepo.GetAccountById - r.Pool.QueryRow: %w", err)
	}

	return account, nil
}

func (r *AccountRepo) Deposit(ctx context.Context, id int, amount int) (err error) {
	stmtSQL := `UPDATE account 
				SET balance = balance + $1
				WHERE id = $2`

	conn, err := r.db.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("AccountRepo.Deposit - r.Pool.Acquire: %w", err)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, stmtSQL, amount, id)
	if err != nil {
		return fmt.Errorf("AccountRepo.Deposit - conn.Exec: %w", err)
	}

	return nil
}
