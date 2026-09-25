package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"uzbekona.dev/backend/internal/model"
)

// ApplicationRepo — "Jamoaga qo'shilish" arizalari (nomzodlar).
type ApplicationRepo struct {
	db DBTX
}

const applicationColumns = `id, full_name, email, phone, telegram, position, experience, portfolio_url, resume_url,
	about, status, note, ip, user_agent, created_at, updated_at`

func scanApplication(row pgx.Row) (*model.Application, error) {
	var a model.Application
	err := row.Scan(&a.ID, &a.FullName, &a.Email, &a.Phone, &a.Telegram, &a.Position, &a.Experience,
		&a.PortfolioURL, &a.ResumeURL, &a.About, &a.Status, &a.Note, &a.IP, &a.UserAgent, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ApplicationRepo) Create(ctx context.Context, a *model.Application) error {
	return r.db.QueryRow(ctx, `INSERT INTO job_applications
		(full_name, email, phone, telegram, position, experience, portfolio_url, resume_url, about, ip, user_agent)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id, status, created_at`,
		a.FullName, a.Email, a.Phone, a.Telegram, a.Position, a.Experience, a.PortfolioURL, a.ResumeURL, a.About, a.IP, a.UserAgent,
	).Scan(&a.ID, &a.Status, &a.CreatedAt)
}

func (r *ApplicationRepo) List(ctx context.Context, params model.ListParams) ([]model.Application, int, error) {
	where := []string{"TRUE"}
	args := []any{}
	if params.Status != "" {
		args = append(args, params.Status)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}
	if params.Query != "" {
		args = append(args, likePattern(params.Query))
		where = append(where, fmt.Sprintf("(full_name ILIKE $%[1]d OR email ILIKE $%[1]d OR position ILIKE $%[1]d OR about ILIKE $%[1]d)", len(args)))
	}
	cond := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, `SELECT count(*) FROM job_applications WHERE `+cond, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args = append(args, params.Limit, params.Offset())
	rows, err := r.db.Query(ctx, fmt.Sprintf(`SELECT %s FROM job_applications WHERE %s
		ORDER BY created_at DESC, id DESC LIMIT $%d OFFSET $%d`, applicationColumns, cond, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := []model.Application{}
	for rows.Next() {
		a, err := scanApplication(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *a)
	}
	return out, total, rows.Err()
}

func (r *ApplicationRepo) GetByID(ctx context.Context, id int64) (*model.Application, error) {
	a, err := scanApplication(r.db.QueryRow(ctx, `SELECT `+applicationColumns+` FROM job_applications WHERE id = $1`, id))
	return a, notFound(err)
}

func (r *ApplicationRepo) Update(ctx context.Context, id int64, in *model.ApplicationUpdateInput) error {
	tag, err := r.db.Exec(ctx, `UPDATE job_applications SET status = $2, note = $3 WHERE id = $1`, id, in.Status, in.Note)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

func (r *ApplicationRepo) Delete(ctx context.Context, id int64) error {
	return deleteByID(ctx, r.db, "job_applications", id)
}
