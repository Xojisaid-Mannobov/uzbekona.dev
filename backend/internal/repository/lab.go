package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"uzbekona.dev/backend/internal/model"
)

type LabRepo struct {
	db DBTX
}

var labColumns = `l.id, l.slug, l.title, l.description, l.stage, l.url, l.repo_url,
	l.cover_id, ` + mediaJSON("l.cover_id") + `, l.stack, l.is_published, l.position, l.created_at, l.updated_at`

func scanLab(row pgx.Row) (*model.Lab, error) {
	var l model.Lab
	err := row.Scan(&l.ID, &l.Slug, &l.Title, &l.Description, &l.Stage, &l.URL, &l.RepoURL,
		&l.CoverID, &l.Cover, &l.Stack, &l.IsPublished, &l.Position, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *LabRepo) list(ctx context.Context, publishedOnly bool) ([]model.Lab, error) {
	sql := `SELECT ` + labColumns + ` FROM labs l`
	if publishedOnly {
		sql += ` WHERE l.is_published`
	}
	rows, err := r.db.Query(ctx, sql+` ORDER BY l.position, l.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Lab{}
	for rows.Next() {
		l, err := scanLab(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *l)
	}
	return out, rows.Err()
}

func (r *LabRepo) ListPublished(ctx context.Context) ([]model.Lab, error) {
	return r.list(ctx, true)
}

func (r *LabRepo) ListAll(ctx context.Context) ([]model.Lab, error) {
	return r.list(ctx, false)
}

func (r *LabRepo) GetByID(ctx context.Context, id int64) (*model.Lab, error) {
	l, err := scanLab(r.db.QueryRow(ctx, `SELECT `+labColumns+` FROM labs l WHERE l.id = $1`, id))
	return l, notFound(err)
}

func (r *LabRepo) SlugExists(ctx context.Context, slug string, excludeID int64) (bool, error) {
	return slugExists(ctx, r.db, "labs", slug, excludeID)
}

func (r *LabRepo) Create(ctx context.Context, slug string, in *model.LabInput) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `INSERT INTO labs
		(slug, title, description, stage, url, repo_url, cover_id, stack, is_published, position)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,(SELECT COALESCE(max(position), 0) + 1 FROM labs)) RETURNING id`,
		slug, in.Title, in.Description, in.Stage, in.URL, in.RepoURL, in.CoverID, strs(in.Stack), in.IsPublished,
	).Scan(&id)
	return id, err
}

func (r *LabRepo) Update(ctx context.Context, id int64, slug string, in *model.LabInput) error {
	tag, err := r.db.Exec(ctx, `UPDATE labs SET slug = $2, title = $3, description = $4, stage = $5,
		url = $6, repo_url = $7, cover_id = $8, stack = $9, is_published = $10 WHERE id = $1`,
		id, slug, in.Title, in.Description, in.Stage, in.URL, in.RepoURL, in.CoverID, strs(in.Stack), in.IsPublished)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

func (r *LabRepo) Reorder(ctx context.Context, ids []int64) error {
	return reorder(ctx, r.db, "labs", ids)
}

func (r *LabRepo) Delete(ctx context.Context, id int64) error {
	return deleteByID(ctx, r.db, "labs", id)
}
