package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"uzbekona.dev/backend/internal/model"
)

type MediaRepo struct {
	db DBTX
}

const mediaColumns = `id, url, kind, mime, width, height, alt, variants, path, original_name, size, created_at`

func scanMedia(row pgx.Row) (*model.Media, error) {
	var m model.Media
	err := row.Scan(&m.ID, &m.URL, &m.Kind, &m.Mime, &m.Width, &m.Height, &m.Alt, &m.Variants,
		&m.Path, &m.OriginalName, &m.Size, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	m.Variants = orEmpty(m.Variants)
	return &m, nil
}

func (r *MediaRepo) List(ctx context.Context, kind string, params model.ListParams) ([]model.Media, int, error) {
	where := []string{"TRUE"}
	args := []any{}
	if kind != "" {
		args = append(args, kind)
		where = append(where, fmt.Sprintf("kind = $%d", len(args)))
	}
	if params.Query != "" {
		args = append(args, likePattern(params.Query))
		where = append(where, fmt.Sprintf("(original_name ILIKE $%[1]d OR alt ILIKE $%[1]d)", len(args)))
	}
	cond := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, `SELECT count(*) FROM media WHERE `+cond, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args = append(args, params.Limit, params.Offset())
	rows, err := r.db.Query(ctx, fmt.Sprintf(`SELECT %s FROM media WHERE %s ORDER BY created_at DESC, id DESC LIMIT $%d OFFSET $%d`,
		mediaColumns, cond, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := []model.Media{}
	for rows.Next() {
		m, err := scanMedia(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *m)
	}
	return out, total, rows.Err()
}

func (r *MediaRepo) GetByID(ctx context.Context, id int64) (*model.Media, error) {
	m, err := scanMedia(r.db.QueryRow(ctx, `SELECT `+mediaColumns+` FROM media WHERE id = $1`, id))
	return m, notFound(err)
}

func (r *MediaRepo) Create(ctx context.Context, m *model.Media) error {
	return r.db.QueryRow(ctx, `INSERT INTO media (kind, path, url, original_name, mime, size, width, height, variants, alt)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id, created_at`,
		m.Kind, m.Path, m.URL, m.OriginalName, m.Mime, m.Size, m.Width, m.Height, orEmpty(m.Variants), m.Alt,
	).Scan(&m.ID, &m.CreatedAt)
}

func (r *MediaRepo) UpdateAlt(ctx context.Context, id int64, alt string) error {
	tag, err := r.db.Exec(ctx, `UPDATE media SET alt = $2 WHERE id = $1`, id, alt)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

func (r *MediaRepo) Delete(ctx context.Context, id int64) error {
	return deleteByID(ctx, r.db, "media", id)
}

// Usage — media qayerlarda ishlatilayotganini sanaydi (o'chirishdan oldin ogohlantirish uchun).
func (r *MediaRepo) Usage(ctx context.Context, id int64) (int, error) {
	var n int
	err := r.db.QueryRow(ctx, `SELECT
		(SELECT count(*) FROM projects WHERE cover_id = $1) +
		(SELECT count(*) FROM project_images WHERE media_id = $1) +
		(SELECT count(*) FROM services WHERE preview_id = $1) +
		(SELECT count(*) FROM team_members WHERE photo_id = $1) +
		(SELECT count(*) FROM articles WHERE cover_id = $1) +
		(SELECT count(*) FROM labs WHERE cover_id = $1) +
		(SELECT count(*) FROM project_blocks WHERE data::text LIKE '%"id": ' || $1::text || ',%') +
		(SELECT count(*) FROM articles WHERE content::text LIKE '%"id": ' || $1::text || ',%')`, id).Scan(&n)
	return n, err
}
