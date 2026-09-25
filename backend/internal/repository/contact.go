package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"uzbekona.dev/backend/internal/model"
)

type ContactRepo struct {
	db DBTX
}

const contactColumns = `id, name, contact, email, project_type, budget, message, status, note, ip, user_agent, created_at, updated_at`

func scanContact(row pgx.Row) (*model.Contact, error) {
	var c model.Contact
	err := row.Scan(&c.ID, &c.Name, &c.Contact, &c.Email, &c.ProjectType, &c.Budget, &c.Message,
		&c.Status, &c.Note, &c.IP, &c.UserAgent, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func collectContacts(rows pgx.Rows) ([]model.Contact, error) {
	defer rows.Close()
	out := []model.Contact{}
	for rows.Next() {
		c, err := scanContact(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return out, rows.Err()
}

func (r *ContactRepo) Create(ctx context.Context, c *model.Contact) error {
	return r.db.QueryRow(ctx, `INSERT INTO contacts (name, contact, email, project_type, budget, message, ip, user_agent)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, status, created_at`,
		c.Name, c.Contact, c.Email, c.ProjectType, c.Budget, c.Message, c.IP, c.UserAgent,
	).Scan(&c.ID, &c.Status, &c.CreatedAt)
}

func (r *ContactRepo) List(ctx context.Context, params model.ListParams) ([]model.Contact, int, error) {
	where := []string{"TRUE"}
	args := []any{}
	if params.Status != "" {
		args = append(args, params.Status)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}
	if params.Query != "" {
		args = append(args, likePattern(params.Query))
		where = append(where, fmt.Sprintf("(name ILIKE $%[1]d OR email ILIKE $%[1]d OR contact ILIKE $%[1]d OR message ILIKE $%[1]d)", len(args)))
	}
	cond := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, `SELECT count(*) FROM contacts WHERE `+cond, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args = append(args, params.Limit, params.Offset())
	rows, err := r.db.Query(ctx, fmt.Sprintf(`SELECT %s FROM contacts WHERE %s ORDER BY created_at DESC, id DESC LIMIT $%d OFFSET $%d`,
		contactColumns, cond, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	items, err := collectContacts(rows)
	return items, total, err
}

func (r *ContactRepo) Recent(ctx context.Context, limit int) ([]model.Contact, error) {
	rows, err := r.db.Query(ctx, `SELECT `+contactColumns+` FROM contacts WHERE status <> 'spam' ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	return collectContacts(rows)
}

func (r *ContactRepo) GetByID(ctx context.Context, id int64) (*model.Contact, error) {
	c, err := scanContact(r.db.QueryRow(ctx, `SELECT `+contactColumns+` FROM contacts WHERE id = $1`, id))
	return c, notFound(err)
}

func (r *ContactRepo) Update(ctx context.Context, id int64, in *model.ContactUpdateInput) error {
	tag, err := r.db.Exec(ctx, `UPDATE contacts SET status = $2, note = $3 WHERE id = $1`, id, in.Status, in.Note)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

func (r *ContactRepo) Delete(ctx context.Context, id int64) error {
	return deleteByID(ctx, r.db, "contacts", id)
}

// ByDay — oxirgi N kun ichida kunlik so'rovlar soni (bo'sh kunlar 0 bilan).
func (r *ContactRepo) ByDay(ctx context.Context, days int) ([]model.DayStat, error) {
	rows, err := r.db.Query(ctx, `SELECT to_char(d.day, 'YYYY-MM-DD'), count(c.id)
		FROM generate_series(current_date - ($1::int - 1), current_date, interval '1 day') AS d(day)
		LEFT JOIN contacts c ON c.created_at::date = d.day::date AND c.status <> 'spam'
		GROUP BY d.day ORDER BY d.day`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.DayStat{}
	for rows.Next() {
		var s model.DayStat
		if err := rows.Scan(&s.Day, &s.Count); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
