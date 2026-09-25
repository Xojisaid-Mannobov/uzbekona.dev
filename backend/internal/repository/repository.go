// Package repository — PostgreSQL bilan ishlash qatlami.
// Barcha so'rovlar parametrlangan ($1, $2 …), shuning uchun SQL injection'dan himoyalangan.
package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("topilmadi")

// DBTX — pool ham, tranzaksiya ham qanoatlantiradigan umumiy interfeys.
type DBTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Repositories struct {
	pool       *pgxpool.Pool
	Admins     *AdminRepo
	Projects   *ProjectRepo
	Services   *ServiceRepo
	Team       *TeamRepo
	Labs       *LabRepo
	Articles   *ArticleRepo
	Categories *CategoryRepo
	Media      *MediaRepo
	Contacts   *ContactRepo
	Settings   *SettingsRepo
}

func New(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		pool:       pool,
		Admins:     &AdminRepo{db: pool},
		Projects:   &ProjectRepo{db: pool},
		Services:   &ServiceRepo{db: pool},
		Team:       &TeamRepo{db: pool},
		Labs:       &LabRepo{db: pool},
		Articles:   &ArticleRepo{db: pool},
		Categories: &CategoryRepo{db: pool},
		Media:      &MediaRepo{db: pool},
		Contacts:   &ContactRepo{db: pool},
		Settings:   &SettingsRepo{db: pool},
	}
}

// InTx funksiyani bitta tranzaksiya ichida bajaradi (xato bo'lsa rollback).
func (r *Repositories) InTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	return pgx.BeginFunc(ctx, r.pool, fn)
}

// Ping — health-check uchun.
func (r *Repositories) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}

// mediaJSON — media jadvalidan bitta yozuvni JSON obyekt sifatida oladigan subquery.
// col — faqat kod ichidagi konstanta (foydalanuvchi kiritmasi emas).
func mediaJSON(col string) string {
	return fmt.Sprintf(`(SELECT jsonb_build_object(
		'id', m.id, 'url', m.url, 'kind', m.kind, 'mime', m.mime,
		'width', m.width, 'height', m.height, 'alt', m.alt, 'variants', m.variants
	) FROM media m WHERE m.id = %s)`, col)
}

// reorder — ids tartibi bo'yicha position ustunini yangilaydi.
func reorder(ctx context.Context, db DBTX, table string, ids []int64) error {
	sql := fmt.Sprintf(`UPDATE %s AS t SET position = x.pos
		FROM unnest($1::bigint[]) WITH ORDINALITY AS x(id, pos)
		WHERE t.id = x.id`, table)
	_, err := db.Exec(ctx, sql, ids)
	return err
}

// slugExists — slug boshqa yozuvda band emasligini tekshiradi.
func slugExists(ctx context.Context, db DBTX, table, slug string, excludeID int64) (bool, error) {
	var exists bool
	sql := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM %s WHERE slug = $1 AND id <> $2)`, table)
	err := db.QueryRow(ctx, sql, slug, excludeID).Scan(&exists)
	return exists, err
}

// deleteByID — yozuvni o'chiradi; topilmasa ErrNotFound.
func deleteByID(ctx context.Context, db DBTX, table string, id int64) error {
	tag, err := db.Exec(ctx, fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, table), id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func notFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

// IsUniqueViolation — unique constraint buzilganini aniqlaydi.
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// IsForeignKeyViolation — mavjud bo'lmagan media/kategoriya ID berilganda.
func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}

// strs — nil slice'ni bo'sh massivga aylantiradi (NOT NULL ustunlar uchun).
func strs(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

func orEmpty[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

// likePattern — ILIKE uchun; foydalanuvchi kiritgan % va _ belgilarini ekranlaydi.
func likePattern(q string) string {
	r := strings.NewReplacer(`\`, `\\`, "%", `\%`, "_", `\_`)
	return "%" + r.Replace(q) + "%"
}
