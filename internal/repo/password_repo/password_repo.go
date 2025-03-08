package password_repo

import (
	"context"
	"database/sql"

	"github.com/vandi37/password-manager/internal/postgresql/database"
	"github.com/vandi37/password-manager/internal/postgresql/module"
	"github.com/vandi37/password-manager/internal/repo"
	"github.com/vandi37/vanerrors"
)

type PasswordRepo struct {
	db *database.DB
}

func New(db *database.DB) *PasswordRepo {
	return &PasswordRepo{db: db}
}

func (r *PasswordRepo) Create(ctx context.Context, password module.Password) error {
	stmt, err := r.db.PrepareContext(ctx, `insert into passwords (company, username, password, nonce, user_id) values ($1, $2, $3, $4, $5);`)
	if err != nil {
		return vanerrors.Wrap(repo.ErrorPreparing, err)
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, password.Company, password.Username, password.Password, password.Nonce, password.UserId)
	if err != nil {
		return vanerrors.Wrap(repo.ErrorExecuting, err)
	}

	return repo.ReturnByRes(res, repo.Equals(1))
}

func (r *PasswordRepo) UpdateUsername(ctx context.Context, passwordId int, username string) error {
	stmt, err := r.db.PrepareContext(ctx, `update passwords set username = $1 where id = $2;`)
	if err != nil {
		return vanerrors.Wrap(repo.ErrorPreparing, err)
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, username, passwordId)
	if err != nil {
		return vanerrors.Wrap(repo.ErrorExecuting, err)
	}
	return repo.ReturnByRes(res, repo.Equals(1))
}

func (r *PasswordRepo) Update(ctx context.Context, passwordId int, password []byte, nonce []byte) error {
	stmt, err := r.db.PrepareContext(ctx, `update passwords set password = $1, nonce = $2 where id = $3;`)
	if err != nil {
		return vanerrors.Wrap(repo.ErrorPreparing, err)
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, password, nonce, passwordId)
	if err != nil {
		return vanerrors.Wrap(repo.ErrorExecuting, err)
	}
	return repo.ReturnByRes(res, repo.Equals(1))
}

func (r *PasswordRepo) Remove(ctx context.Context, passwordId int) error {
	res, err := r.db.ExecContext(ctx, "delete from passwords where id = $1", passwordId)
	if err != nil {
		return vanerrors.Wrap(repo.ErrorExecuting, err)
	}

	return repo.ReturnByRes(res, repo.Equals(1))
}

func (r *PasswordRepo) GetByUserId(ctx context.Context, id int64) ([]module.Password, error) {
	rows, err := r.db.QueryContext(ctx, "select id, company, username, password, nonce, user_id from passwords where user_id = $1", id)
	if err != nil {
		return nil, vanerrors.Wrap(repo.ErrorExecuting, err)
	}

	defer rows.Close()

	return scanPasswordRows(rows)
}

func (r *PasswordRepo) GetByCompany(ctx context.Context, id int64, company string) ([]module.Password, error) {
	stmt, err := r.db.PrepareContext(ctx, `select id, company, username, password, nonce, user_id from passwords where user_id = $1 and company = $2;`)
	if err != nil {
		return nil, vanerrors.Wrap(repo.ErrorPreparing, err)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id, company)
	if err != nil {
		return nil, vanerrors.Wrap(repo.ErrorExecuting, err)
	}

	defer rows.Close()

	return scanPasswordRows(rows)
}

func scanPasswordRows(rows *sql.Rows) ([]module.Password, error) {
	var res []module.Password
	for rows.Next() {
		var password module.Password
		err := rows.Scan(&password.Id, &password.Company, &password.Username, &password.Password, &password.Nonce, &password.UserId)
		if err != nil {
			return nil, vanerrors.Wrap(repo.ErrorScanning, err)
		}
		res = append(res, password)
	}
	return res, nil
}
