package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"uzbekona.dev/backend/internal/model"
)

type AdminRepo struct {
	db DBTX
}

const adminColumns = `id, name, email, password_hash, last_login_at, created_at`

func scanAdmin(row pgx.Row) (*model.Admin, error) {
	var a model.Admin
	if err := row.Scan(&a.ID, &a.Name, &a.Email, &a.PasswordHash, &a.LastLoginAt, &a.CreatedAt); err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AdminRepo) GetByEmail(ctx context.Context, email string) (*model.Admin, error) {
	a, err := scanAdmin(r.db.QueryRow(ctx, `SELECT `+adminColumns+` FROM admins WHERE lower(email) = lower($1)`, email))
	return a, notFound(err)
}

func (r *AdminRepo) GetByID(ctx context.Context, id int64) (*model.Admin, error) {
	a, err := scanAdmin(r.db.QueryRow(ctx, `SELECT `+adminColumns+` FROM admins WHERE id = $1`, id))
	return a, notFound(err)
}

func (r *AdminRepo) List(ctx context.Context) ([]model.Admin, error) {
	rows, err := r.db.Query(ctx, `SELECT `+adminColumns+` FROM admins ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Admin{}
	for rows.Next() {
		a, err := scanAdmin(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *a)
	}
	return out, rows.Err()
}

func (r *AdminRepo) Count(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRow(ctx, `SELECT count(*) FROM admins`).Scan(&n)
	return n, err
}

func (r *AdminRepo) Create(ctx context.Context, name, email, hash string) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `INSERT INTO admins (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`,
		name, email, hash).Scan(&id)
	return id, err
}

func (r *AdminRepo) UpdatePassword(ctx context.Context, id int64, hash string) error {
	_, err := r.db.Exec(ctx, `UPDATE admins SET password_hash = $2 WHERE id = $1`, id, hash)
	return err
}

func (r *AdminRepo) TouchLogin(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `UPDATE admins SET last_login_at = now() WHERE id = $1`, id)
	return err
}

func (r *AdminRepo) Delete(ctx context.Context, id int64) error {
	return deleteByID(ctx, r.db, "admins", id)
}
