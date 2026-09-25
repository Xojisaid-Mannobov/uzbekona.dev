package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"uzbekona.dev/backend/internal/model"
)

type ProjectRepo struct {
	db DBTX
}

var projectColumns = `p.id, p.slug, p.title, p.tagline, p.short_description, p.full_description,
	p.cover_id, ` + mediaJSON("p.cover_id") + `, p.year, p.client, p.industry,
	p.platforms, p.services, p.stack, p.metrics, p.live_url, p.accent, p.status,
	p.featured, p.position, p.seo, p.published_at, p.created_at, p.updated_at`

func scanProject(row pgx.Row) (*model.Project, error) {
	var p model.Project
	err := row.Scan(
		&p.ID, &p.Slug, &p.Title, &p.Tagline, &p.ShortDescription, &p.FullDescription,
		&p.CoverID, &p.Cover, &p.Year, &p.Client, &p.Industry,
		&p.Platforms, &p.Services, &p.Stack, &p.Metrics, &p.LiveURL, &p.Accent, &p.Status,
		&p.Featured, &p.Position, &p.SEO, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	p.Metrics = orEmpty(p.Metrics)
	return &p, nil
}

func collectProjects(rows pgx.Rows) ([]model.Project, error) {
	defer rows.Close()
	out := []model.Project{}
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *p)
	}
	return out, rows.Err()
}

// ListPublished — public sayt uchun e'lon qilingan loyihalar.
func (r *ProjectRepo) ListPublished(ctx context.Context, featuredOnly bool, limit int) ([]model.Project, error) {
	sql := `SELECT ` + projectColumns + ` FROM projects p WHERE p.status = 'published'`
	if featuredOnly {
		sql += ` AND p.featured`
	}
	sql += ` ORDER BY p.position, p.id LIMIT $1`
	rows, err := r.db.Query(ctx, sql, limit)
	if err != nil {
		return nil, err
	}
	return collectProjects(rows)
}

// ListAdmin — admin panel ro'yxati: qidiruv, status filtri va pagination.
func (r *ProjectRepo) ListAdmin(ctx context.Context, params model.ListParams) ([]model.Project, int, error) {
	where := []string{"TRUE"}
	args := []any{}
	if params.Status != "" {
		args = append(args, params.Status)
		where = append(where, fmt.Sprintf("p.status = $%d", len(args)))
	}
	if params.Query != "" {
		args = append(args, likePattern(params.Query))
		where = append(where, fmt.Sprintf("(p.title ILIKE $%[1]d OR p.client ILIKE $%[1]d OR p.slug ILIKE $%[1]d)", len(args)))
	}
	cond := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, `SELECT count(*) FROM projects p WHERE `+cond, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, params.Limit, params.Offset())
	rows, err := r.db.Query(ctx, fmt.Sprintf(`SELECT %s FROM projects p WHERE %s
		ORDER BY p.position, p.id LIMIT $%d OFFSET $%d`, projectColumns, cond, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	items, err := collectProjects(rows)
	return items, total, err
}

func (r *ProjectRepo) Recent(ctx context.Context, limit int) ([]model.Project, error) {
	rows, err := r.db.Query(ctx, `SELECT `+projectColumns+` FROM projects p ORDER BY p.updated_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	return collectProjects(rows)
}

func (r *ProjectRepo) GetByID(ctx context.Context, id int64) (*model.Project, error) {
	p, err := scanProject(r.db.QueryRow(ctx, `SELECT `+projectColumns+` FROM projects p WHERE p.id = $1`, id))
	return p, notFound(err)
}

func (r *ProjectRepo) GetPublishedBySlug(ctx context.Context, slug string) (*model.Project, error) {
	p, err := scanProject(r.db.QueryRow(ctx,
		`SELECT `+projectColumns+` FROM projects p WHERE p.slug = $1 AND p.status = 'published'`, slug))
	return p, notFound(err)
}

// Next — tartib bo'yicha keyingi e'lon qilingan loyiha (oxirgisidan keyin — birinchisi).
func (r *ProjectRepo) Next(ctx context.Context, current *model.Project) (*model.Project, error) {
	p, err := scanProject(r.db.QueryRow(ctx, `SELECT `+projectColumns+` FROM projects p
		WHERE p.status = 'published' AND p.id <> $1
		ORDER BY ((p.position, p.id) <= ($2, $1)), p.position, p.id
		LIMIT 1`, current.ID, current.Position))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

func (r *ProjectRepo) Blocks(ctx context.Context, projectID int64) ([]model.Block, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, type, data FROM project_blocks WHERE project_id = $1 ORDER BY position, id`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Block{}
	for rows.Next() {
		var b model.Block
		if err := rows.Scan(&b.ID, &b.Type, &b.Data); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *ProjectRepo) Gallery(ctx context.Context, projectID int64) ([]model.GalleryItem, error) {
	rows, err := r.db.Query(ctx, `SELECT pi.media_id, pi.caption, `+mediaJSON("pi.media_id")+`
		FROM project_images pi WHERE pi.project_id = $1 ORDER BY pi.position, pi.id`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.GalleryItem{}
	for rows.Next() {
		var g model.GalleryItem
		if err := rows.Scan(&g.MediaID, &g.Caption, &g.Media); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

func (r *ProjectRepo) SlugExists(ctx context.Context, slug string, excludeID int64) (bool, error) {
	return slugExists(ctx, r.db, "projects", slug, excludeID)
}

// Create — yangi loyiha; position oxiriga qo'yiladi.
func (r *ProjectRepo) Create(ctx context.Context, db DBTX, slug string, in *model.ProjectInput) (int64, error) {
	var id int64
	err := db.QueryRow(ctx, `INSERT INTO projects (
			slug, title, tagline, short_description, full_description, cover_id, year, client, industry,
			platforms, services, stack, metrics, live_url, accent, status, featured, seo, published_at, position
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,
			CASE WHEN $16 = 'published' THEN now() END,
			(SELECT COALESCE(max(position), 0) + 1 FROM projects))
		RETURNING id`,
		slug, in.Title, in.Tagline, in.ShortDescription, in.FullDescription, in.CoverID, in.Year,
		in.Client, in.Industry, strs(in.Platforms), strs(in.Services), strs(in.Stack),
		orEmpty(in.Metrics), in.LiveURL, in.Accent, in.Status, in.Featured, in.SEO,
	).Scan(&id)
	return id, err
}

func (r *ProjectRepo) Update(ctx context.Context, db DBTX, id int64, slug string, in *model.ProjectInput) error {
	tag, err := db.Exec(ctx, `UPDATE projects SET
			slug = $2, title = $3, tagline = $4, short_description = $5, full_description = $6,
			cover_id = $7, year = $8, client = $9, industry = $10, platforms = $11, services = $12,
			stack = $13, metrics = $14, live_url = $15, accent = $16, status = $17, featured = $18, seo = $19,
			published_at = CASE WHEN $17 = 'published' THEN COALESCE(published_at, now()) ELSE published_at END
		WHERE id = $1`,
		id, slug, in.Title, in.Tagline, in.ShortDescription, in.FullDescription, in.CoverID, in.Year,
		in.Client, in.Industry, strs(in.Platforms), strs(in.Services), strs(in.Stack),
		orEmpty(in.Metrics), in.LiveURL, in.Accent, in.Status, in.Featured, in.SEO,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ReplaceBlocks — content builder bloklarini to'liq almashtiradi (tartib saqlanadi).
func (r *ProjectRepo) ReplaceBlocks(ctx context.Context, db DBTX, projectID int64, blocks []model.Block) error {
	if _, err := db.Exec(ctx, `DELETE FROM project_blocks WHERE project_id = $1`, projectID); err != nil {
		return err
	}
	if len(blocks) == 0 {
		return nil
	}
	rows := make([][]any, len(blocks))
	for i, b := range blocks {
		rows[i] = []any{projectID, b.Type, i + 1, b.Data}
	}
	_, err := copyRows(ctx, db, "project_blocks", []string{"project_id", "type", "position", "data"}, rows)
	return err
}

func (r *ProjectRepo) ReplaceGallery(ctx context.Context, db DBTX, projectID int64, items []model.GalleryItem) error {
	if _, err := db.Exec(ctx, `DELETE FROM project_images WHERE project_id = $1`, projectID); err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	rows := make([][]any, len(items))
	for i, g := range items {
		rows[i] = []any{projectID, g.MediaID, g.Caption, i + 1}
	}
	_, err := copyRows(ctx, db, "project_images", []string{"project_id", "media_id", "caption", "position"}, rows)
	return err
}

func (r *ProjectRepo) SetStatus(ctx context.Context, id int64, status string) error {
	tag, err := r.db.Exec(ctx, `UPDATE projects SET status = $2,
		published_at = CASE WHEN $2 = 'published' THEN COALESCE(published_at, now()) ELSE published_at END
		WHERE id = $1`, id, status)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

func (r *ProjectRepo) SetFeatured(ctx context.Context, id int64, featured bool) error {
	tag, err := r.db.Exec(ctx, `UPDATE projects SET featured = $2 WHERE id = $1`, id, featured)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

func (r *ProjectRepo) Reorder(ctx context.Context, ids []int64) error {
	return reorder(ctx, r.db, "projects", ids)
}

func (r *ProjectRepo) Delete(ctx context.Context, id int64) error {
	return deleteByID(ctx, r.db, "projects", id)
}

// copyRows — bir nechta qatorni bitta INSERT bilan yozadi.
func copyRows(ctx context.Context, db DBTX, table string, columns []string, rows [][]any) (int64, error) {
	var (
		placeholders []string
		args         []any
	)
	for _, row := range rows {
		ph := make([]string, len(row))
		for j, v := range row {
			args = append(args, v)
			ph[j] = fmt.Sprintf("$%d", len(args))
		}
		placeholders = append(placeholders, "("+strings.Join(ph, ",")+")")
	}
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", table, strings.Join(columns, ","), strings.Join(placeholders, ","))
	tag, err := db.Exec(ctx, sql, args...)
	return tag.RowsAffected(), err
}
