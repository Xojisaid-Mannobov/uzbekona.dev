package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"uzbekona.dev/backend/internal/model"
)

type ArticleRepo struct {
	db DBTX
}

// articleColumns — ro'yxatlarda og'ir content ustuni o'rniga bo'sh massiv olinadi.
func articleColumns(withContent bool) string {
	content := `'[]'::jsonb`
	if withContent {
		content = `a.content`
	}
	return `a.id, a.slug, a.title, a.excerpt, a.cover_id, ` + mediaJSON("a.cover_id") + `,
		a.category_id, (SELECT jsonb_build_object('id', c.id, 'name', c.name, 'slug', c.slug, 'position', c.position)
			FROM article_categories c WHERE c.id = a.category_id),
		` + content + `, a.reading_time, a.author_name, a.status, a.featured, a.seo,
		a.published_at, a.created_at, a.updated_at`
}

func scanArticle(row pgx.Row) (*model.Article, error) {
	var a model.Article
	err := row.Scan(&a.ID, &a.Slug, &a.Title, &a.Excerpt, &a.CoverID, &a.Cover,
		&a.CategoryID, &a.Category, &a.Content, &a.ReadingTime, &a.AuthorName, &a.Status, &a.Featured, &a.SEO,
		&a.PublishedAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func collectArticles(rows pgx.Rows) ([]model.Article, error) {
	defer rows.Close()
	out := []model.Article{}
	for rows.Next() {
		a, err := scanArticle(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *a)
	}
	return out, rows.Err()
}

// ListPublished — Journal sahifasi: kategoriya bo'yicha filtr va pagination.
func (r *ArticleRepo) ListPublished(ctx context.Context, f model.ArticleFilter) ([]model.Article, int, error) {
	cond := `a.status = 'published' AND a.published_at <= now()`
	args := []any{}
	if f.Category != "" {
		args = append(args, f.Category)
		cond += fmt.Sprintf(` AND a.category_id = (SELECT id FROM article_categories WHERE slug = $%d)`, len(args))
	}

	var total int
	if err := r.db.QueryRow(ctx, `SELECT count(*) FROM articles a WHERE `+cond, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, f.Limit, (f.Page-1)*f.Limit)
	rows, err := r.db.Query(ctx, fmt.Sprintf(`SELECT %s FROM articles a WHERE %s
		ORDER BY a.published_at DESC, a.id DESC LIMIT $%d OFFSET $%d`,
		articleColumns(false), cond, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	items, err := collectArticles(rows)
	return items, total, err
}

func (r *ArticleRepo) ListAdmin(ctx context.Context, params model.ListParams) ([]model.Article, int, error) {
	where := []string{"TRUE"}
	args := []any{}
	if params.Status != "" {
		args = append(args, params.Status)
		where = append(where, fmt.Sprintf("a.status = $%d", len(args)))
	}
	if params.Query != "" {
		args = append(args, likePattern(params.Query))
		where = append(where, fmt.Sprintf("(a.title ILIKE $%[1]d OR a.excerpt ILIKE $%[1]d)", len(args)))
	}
	cond := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, `SELECT count(*) FROM articles a WHERE `+cond, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args = append(args, params.Limit, params.Offset())
	rows, err := r.db.Query(ctx, fmt.Sprintf(`SELECT %s FROM articles a WHERE %s
		ORDER BY a.created_at DESC, a.id DESC LIMIT $%d OFFSET $%d`,
		articleColumns(false), cond, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	items, err := collectArticles(rows)
	return items, total, err
}

func (r *ArticleRepo) GetByID(ctx context.Context, id int64) (*model.Article, error) {
	a, err := scanArticle(r.db.QueryRow(ctx, `SELECT `+articleColumns(true)+` FROM articles a WHERE a.id = $1`, id))
	return a, notFound(err)
}

func (r *ArticleRepo) GetPublishedBySlug(ctx context.Context, slug string) (*model.Article, error) {
	a, err := scanArticle(r.db.QueryRow(ctx, `SELECT `+articleColumns(true)+` FROM articles a
		WHERE a.slug = $1 AND a.status = 'published' AND a.published_at <= now()`, slug))
	return a, notFound(err)
}

// Related — shu kategoriyadagi (bo'lmasa — eng yangi) boshqa maqolalar.
func (r *ArticleRepo) Related(ctx context.Context, a *model.Article, limit int) ([]model.Article, error) {
	rows, err := r.db.Query(ctx, `SELECT `+articleColumns(false)+` FROM articles a
		WHERE a.status = 'published' AND a.published_at <= now() AND a.id <> $1
		ORDER BY (a.category_id IS NOT DISTINCT FROM $2) DESC, a.published_at DESC
		LIMIT $3`, a.ID, a.CategoryID, limit)
	if err != nil {
		return nil, err
	}
	return collectArticles(rows)
}

func (r *ArticleRepo) SlugExists(ctx context.Context, slug string, excludeID int64) (bool, error) {
	return slugExists(ctx, r.db, "articles", slug, excludeID)
}

func (r *ArticleRepo) Create(ctx context.Context, slug string, readingTime int, in *model.ArticleInput) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `INSERT INTO articles
		(slug, title, excerpt, cover_id, category_id, content, reading_time, author_name, status, featured, seo, published_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,
			CASE WHEN $9 = 'published' THEN COALESCE($12, now()) ELSE $12 END)
		RETURNING id`,
		slug, in.Title, in.Excerpt, in.CoverID, in.CategoryID, orEmpty(in.Content), readingTime,
		in.AuthorName, in.Status, in.Featured, in.SEO, in.PublishedAt,
	).Scan(&id)
	return id, err
}

func (r *ArticleRepo) Update(ctx context.Context, id int64, slug string, readingTime int, in *model.ArticleInput) error {
	tag, err := r.db.Exec(ctx, `UPDATE articles SET slug = $2, title = $3, excerpt = $4, cover_id = $5,
		category_id = $6, content = $7, reading_time = $8, author_name = $9, status = $10, featured = $11, seo = $12,
		published_at = CASE WHEN $10 = 'published' THEN COALESCE($13, published_at, now()) ELSE COALESCE($13, published_at) END
		WHERE id = $1`,
		id, slug, in.Title, in.Excerpt, in.CoverID, in.CategoryID, orEmpty(in.Content), readingTime,
		in.AuthorName, in.Status, in.Featured, in.SEO, in.PublishedAt)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

func (r *ArticleRepo) Delete(ctx context.Context, id int64) error {
	return deleteByID(ctx, r.db, "articles", id)
}

// ─── Kategoriyalar ──────────────────────────────────────────

type CategoryRepo struct {
	db DBTX
}

// List — har bir kategoriya bilan e'lon qilingan maqolalar soni.
func (r *CategoryRepo) List(ctx context.Context) ([]model.ArticleCategory, error) {
	rows, err := r.db.Query(ctx, `SELECT c.id, c.name, c.slug, c.position,
			(SELECT count(*) FROM articles a WHERE a.category_id = c.id AND a.status = 'published')
		FROM article_categories c ORDER BY c.position, c.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.ArticleCategory{}
	for rows.Next() {
		var c model.ArticleCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Position, &c.Count); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *CategoryRepo) SlugExists(ctx context.Context, slug string, excludeID int64) (bool, error) {
	return slugExists(ctx, r.db, "article_categories", slug, excludeID)
}

func (r *CategoryRepo) Create(ctx context.Context, slug string, in *model.CategoryInput) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `INSERT INTO article_categories (name, slug, position) VALUES ($1,$2,$3) RETURNING id`,
		in.Name, slug, in.Position).Scan(&id)
	return id, err
}

func (r *CategoryRepo) Update(ctx context.Context, id int64, slug string, in *model.CategoryInput) error {
	tag, err := r.db.Exec(ctx, `UPDATE article_categories SET name = $2, slug = $3, position = $4 WHERE id = $1`,
		id, in.Name, slug, in.Position)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

func (r *CategoryRepo) Delete(ctx context.Context, id int64) error {
	return deleteByID(ctx, r.db, "article_categories", id)
}
