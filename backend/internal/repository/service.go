package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"uzbekona.dev/backend/internal/model"
)

type ServiceRepo struct {
	db DBTX
}

var serviceColumns = `s.id, s.slug, s.title, s.summary, s.description, s.features, s.stack,
	s.preview_id, ` + mediaJSON("s.preview_id") + `, s.status, s.position, s.seo, s.created_at, s.updated_at`

func scanService(row pgx.Row) (*model.Service, error) {
	var s model.Service
	err := row.Scan(&s.ID, &s.Slug, &s.Title, &s.Summary, &s.Description, &s.Features, &s.Stack,
		&s.PreviewID, &s.Preview, &s.Status, &s.Position, &s.SEO, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ServiceRepo) list(ctx context.Context, publishedOnly bool) ([]model.Service, error) {
	sql := `SELECT ` + serviceColumns + ` FROM services s`
	if publishedOnly {
		sql += ` WHERE s.status = 'published'`
	}
	rows, err := r.db.Query(ctx, sql+` ORDER BY s.position, s.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Service{}
	for rows.Next() {
		s, err := scanService(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *s)
	}
	return out, rows.Err()
}

func (r *ServiceRepo) ListPublished(ctx context.Context) ([]model.Service, error) {
	return r.list(ctx, true)
}

func (r *ServiceRepo) ListAll(ctx context.Context) ([]model.Service, error) {
	return r.list(ctx, false)
}

func (r *ServiceRepo) GetByID(ctx context.Context, id int64) (*model.Service, error) {
	s, err := scanService(r.db.QueryRow(ctx, `SELECT `+serviceColumns+` FROM services s WHERE s.id = $1`, id))
	return s, notFound(err)
}

func (r *ServiceRepo) GetPublishedBySlug(ctx context.Context, slug string) (*model.Service, error) {
	s, err := scanService(r.db.QueryRow(ctx,
		`SELECT `+serviceColumns+` FROM services s WHERE s.slug = $1 AND s.status = 'published'`, slug))
	return s, notFound(err)
}

func (r *ServiceRepo) SlugExists(ctx context.Context, slug string, excludeID int64) (bool, error) {
	return slugExists(ctx, r.db, "services", slug, excludeID)
}

func (r *ServiceRepo) Create(ctx context.Context, slug string, in *model.ServiceInput) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `INSERT INTO services
		(slug, title, summary, description, features, stack, preview_id, status, seo, position)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,(SELECT COALESCE(max(position), 0) + 1 FROM services))
		RETURNING id`,
		slug, in.Title, in.Summary, in.Description, strs(in.Features), strs(in.Stack), in.PreviewID, in.Status, in.SEO,
	).Scan(&id)
	return id, err
}

func (r *ServiceRepo) Update(ctx context.Context, id int64, slug string, in *model.ServiceInput) error {
	tag, err := r.db.Exec(ctx, `UPDATE services SET slug = $2, title = $3, summary = $4, description = $5,
		features = $6, stack = $7, preview_id = $8, status = $9, seo = $10 WHERE id = $1`,
		id, slug, in.Title, in.Summary, in.Description, strs(in.Features), strs(in.Stack), in.PreviewID, in.Status, in.SEO)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

func (r *ServiceRepo) Reorder(ctx context.Context, ids []int64) error {
	return reorder(ctx, r.db, "services", ids)
}

func (r *ServiceRepo) Delete(ctx context.Context, id int64) error {
	return deleteByID(ctx, r.db, "services", id)
}
