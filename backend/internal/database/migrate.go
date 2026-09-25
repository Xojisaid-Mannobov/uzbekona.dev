package database

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"uzbekona.dev/backend/migrations"
)

// Bir vaqtda bir nechta instance migratsiya qilmasligi uchun advisory lock kaliti
const migrationLockID = 727_001

type migration struct {
	version string
	name    string
}

// Migrate hali qo'llanmagan barcha *.up.sql fayllarni tartib bilan bajaradi.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	return withLock(ctx, pool, func(conn *pgxpool.Conn) error {
		applied, err := appliedVersions(ctx, conn)
		if err != nil {
			return err
		}

		for _, m := range listMigrations(".up.sql") {
			if applied[m.version] {
				continue
			}
			sql, err := fs.ReadFile(migrations.FS, m.name)
			if err != nil {
				return err
			}
			err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
				if _, err := tx.Exec(ctx, string(sql)); err != nil {
					return err
				}
				_, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, m.version)
				return err
			})
			if err != nil {
				return fmt.Errorf("migratsiya %s: %w", m.name, err)
			}
			slog.Info("migratsiya qo'llandi", "file", m.name)
		}
		return nil
	})
}

// Rollback oxirgi qo'llangan migratsiyani bekor qiladi.
func Rollback(ctx context.Context, pool *pgxpool.Pool) error {
	return withLock(ctx, pool, func(conn *pgxpool.Conn) error {
		if _, err := appliedVersions(ctx, conn); err != nil {
			return err
		}
		var version string
		err := conn.QueryRow(ctx, `SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1`).Scan(&version)
		if err == pgx.ErrNoRows {
			slog.Info("bekor qilinadigan migratsiya yo'q")
			return nil
		}
		if err != nil {
			return err
		}

		for _, m := range listMigrations(".down.sql") {
			if m.version != version {
				continue
			}
			sql, err := fs.ReadFile(migrations.FS, m.name)
			if err != nil {
				return err
			}
			return pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
				if _, err := tx.Exec(ctx, string(sql)); err != nil {
					return err
				}
				_, err := tx.Exec(ctx, `DELETE FROM schema_migrations WHERE version = $1`, version)
				slog.Info("migratsiya bekor qilindi", "file", m.name)
				return err
			})
		}
		return fmt.Errorf("%s uchun .down.sql topilmadi", version)
	})
}

func withLock(ctx context.Context, pool *pgxpool.Pool, fn func(conn *pgxpool.Conn) error) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if _, err := conn.Exec(ctx, `SELECT pg_advisory_lock($1)`, migrationLockID); err != nil {
		return err
	}
	defer conn.Exec(context.Background(), `SELECT pg_advisory_unlock($1)`, migrationLockID) //nolint:errcheck

	return fn(conn)
}

func appliedVersions(ctx context.Context, conn *pgxpool.Conn) (map[string]bool, error) {
	_, err := conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version    TEXT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	versions, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, err
	}

	applied := make(map[string]bool, len(versions))
	for _, v := range versions {
		applied[v] = true
	}
	return applied, nil
}

// listMigrations fayllarni versiya bo'yicha tartiblangan holda qaytaradi.
// Fayl nomi formati: 0001_nomi.up.sql
func listMigrations(suffix string) []migration {
	entries, _ := fs.ReadDir(migrations.FS, ".")
	var out []migration
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, suffix) {
			continue
		}
		version, _, _ := strings.Cut(name, "_")
		out = append(out, migration{version: version, name: name})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].version < out[j].version })
	return out
}
