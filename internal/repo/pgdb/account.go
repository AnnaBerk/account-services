package pgdb

import (
	"account-management/pkg/psql"
	"context"
)

type AccountRepo struct {
	db *psql.Postgres
}

func NewAccountRepo(pg *psql.Postgres) *UserRepo {
	return &UserRepo{db: pg}
}

func (r *AccountRepo) CreateAccount(ctx context.Context) (int, error) {

	//stmtSQL := `INSERT INTO account (name, password) VALUES($1, $2) RETURNING id `
	//sql, args, _ := r.Builder.
	//	Insert("accounts").
	//	Values(squirrel.Expr("DEFAULT")).
	//	Suffix("RETURNING id").
	//	ToSql()
	//
	//var id int
	//err := r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	//if err != nil {
	//	log.Debugf("err: %v", err)
	//	var pgErr *pgconn.PgError
	//	if ok := errors.As(err, &pgErr); ok {
	//		if pgErr.Code == "23505" {
	//			return 0, repoerrs.ErrAlreadyExists
	//		}
	//	}
	//	return 0, fmt.Errorf("AccountRepo.CreateAccount - r.Pool.QueryRow: %v", err)
	//}
	//
	return 0, nil
}
