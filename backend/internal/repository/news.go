package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"uzbekona.dev/backend/internal/model"
)

type NewsRepo struct {
	db DBTX
}

// newsColumns — ro'yxatlarda og'ir content ustuni o'rniga bo'sh massiv olinadi.
func newsColumns(withContent bool) string {
	content := `'[]'::jsonb`
	if withContent {
		content = `n.content`
	}
	return `n.id, n.slug, n.title, n.excerpt, n.tag, n.cover_id, ` + mediaJSON("n.cover_id") + `,
		` + content + `, n.status, n.pinned, n.views, n.seo, n.published_at, n.created_at, n.updated_at`
}

func scanNews(row pgx.Row) (*model.News, error) {
	var n model.News
	err := row.Scan(&n.ID, &n.Slug, &n.Title, &n.Excerpt, &n.Tag, &n.CoverID, &n.Cover,
		&n.Content, &n.Status, &n.Pinned, &n.Views, &n.SEO, &n.PublishedAt, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func collectNews(rows pgx.Rows) ([]model.News, error) {
	defer rows.Close()
	out := []model.News{}
	for rows.Next() {
		n, err := scanNews(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *n)
	}
	return out, rows.Err()
}

const newsPublished = `n.status = 'published' AND n.published_at <= now()`

// ListPublished — public ro'yxat: mahkamlangan (pinned) yangilik birinchi, qolgani sanasi bo'yicha.
func (r *NewsRepo) ListPublished(ctx context.Context, page, limit int) ([]model.News, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `SELECT count(*) FROM news n WHERE `+newsPublished).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx, `SELECT `+newsColumns(false)+` FROM news n WHERE `+newsPublished+`
		ORDER BY n.pinned DESC, n.published_at DESC, n.id DESC LIMIT $1 OFFSET $2`, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, err
	}
	items, err := collectNews(rows)
	return items, total, err
}

func (r *NewsRepo) ListAdmin(ctx context.Context, params model.ListParams) ([]model.News, int, error) {
	where := []string{"TRUE"}
	args := []any{}
	if params.Status != "" {
		args = append(args, params.Status)
		where = append(where, fmt.Sprintf("n.status = $%d", len(args)))
	}
	if params.Query != "" {
		args = append(args, likePattern(params.Query))
		where = append(where, fmt.Sprintf("(n.title ILIKE $%[1]d OR n.excerpt ILIKE $%[1]d OR n.tag ILIKE $%[1]d)", len(args)))
	}
	cond := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, `SELECT count(*) FROM news n WHERE `+cond, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args = append(args, params.Limit, params.Offset())
	rows, err := r.db.Query(ctx, fmt.Sprintf(`SELECT %s FROM news n WHERE %s
		ORDER BY COALESCE(n.published_at, n.created_at) DESC, n.id DESC LIMIT $%d OFFSET $%d`,
		newsColumns(false), cond, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	items, err := collectNews(rows)
	return items, total, err
}

func (r *NewsRepo) GetByID(ctx context.Context, id int64) (*model.News, error) {
	n, err := scanNews(r.db.QueryRow(ctx, `SELECT `+newsColumns(true)+` FROM news n WHERE n.id = $1`, id))
	return n, notFound(err)
}

func (r *NewsRepo) GetPublishedBySlug(ctx context.Context, slug string) (*model.News, error) {
	n, err := scanNews(r.db.QueryRow(ctx, `SELECT `+newsColumns(true)+` FROM news n
		WHERE n.slug = $1 AND `+newsPublished, slug))
	return n, notFound(err)
}

// Latest — batafsil sahifa ostidagi "boshqa yangiliklar".
func (r *NewsRepo) Latest(ctx context.Context, excludeID int64, limit int) ([]model.News, error) {
	rows, err := r.db.Query(ctx, `SELECT `+newsColumns(false)+` FROM news n
		WHERE `+newsPublished+` AND n.id <> $1
		ORDER BY n.published_at DESC LIMIT $2`, excludeID, limit)
	if err != nil {
		return nil, err
	}
	return collectNews(rows)
}

// TopViewed — statistika sahifasi uchun eng ko'p o'qilgan yangiliklar.
func (r *NewsRepo) TopViewed(ctx context.Context, limit int) ([]model.NewsViews, error) {
	rows, err := r.db.Query(ctx, `SELECT id, slug, title, views FROM news
		WHERE status = 'published' ORDER BY views DESC, published_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.NewsViews{}
	for rows.Next() {
		var n model.NewsViews
		if err := rows.Scan(&n.ID, &n.Slug, &n.Title, &n.Views); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

// IncrementViews — e'lon qilingan yangilik ko'rishlar sonini oshiradi (tracker chaqiradi).
func (r *NewsRepo) IncrementViews(ctx context.Context, slug string) error {
	_, err := r.db.Exec(ctx, `UPDATE news SET views = views + 1 WHERE slug = $1 AND status = 'published'`, slug)
	return err
}

func (r *NewsRepo) SlugExists(ctx context.Context, slug string, excludeID int64) (bool, error) {
	return slugExists(ctx, r.db, "news", slug, excludeID)
}

func (r *NewsRepo) Create(ctx context.Context, slug string, in *model.NewsInput) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `INSERT INTO news
		(slug, title, excerpt, tag, cover_id, content, status, pinned, seo, published_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,
			CASE WHEN $7 = 'published' THEN COALESCE($10, now()) ELSE $10 END)
		RETURNING id`,
		slug, in.Title, in.Excerpt, in.Tag, in.CoverID, orEmpty(in.Content), in.Status, in.Pinned, in.SEO, in.PublishedAt,
	).Scan(&id)
	return id, err
}

func (r *NewsRepo) Update(ctx context.Context, id int64, slug string, in *model.NewsInput) error {
	tag, err := r.db.Exec(ctx, `UPDATE news SET slug = $2, title = $3, excerpt = $4, tag = $5, cover_id = $6,
		content = $7, status = $8, pinned = $9, seo = $10,
		published_at = CASE WHEN $8 = 'published' THEN COALESCE($11, published_at, now()) ELSE COALESCE($11, published_at) END
		WHERE id = $1`,
		id, slug, in.Title, in.Excerpt, in.Tag, in.CoverID, orEmpty(in.Content), in.Status, in.Pinned, in.SEO, in.PublishedAt)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

func (r *NewsRepo) Delete(ctx context.Context, id int64) error {
	return deleteByID(ctx, r.db, "news", id)
}

// UnpinOthers — bir vaqtda faqat bitta yangilik mahkamlangan bo'ladi.
func (r *NewsRepo) UnpinOthers(ctx context.Context, keepID int64) error {
	_, err := r.db.Exec(ctx, `UPDATE news SET pinned = false WHERE pinned AND id <> $1`, keepID)
	return err
}
