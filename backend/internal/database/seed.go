package database

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed seed/demo.sql
var demoSQL string

//go:embed seed/news.sql
var newsSQL string

// SeedDemo boshlang'ich kontentni (xizmatlar, kategoriyalar, sozlamalar, demo loyihalar)
// faqat bir marta — baza bo'sh bo'lganda — yozadi.
func SeedDemo(ctx context.Context, pool *pgxpool.Pool) error {
	var seeded bool
	err := pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM settings WHERE key = 'seeded')`).Scan(&seeded)
	if err != nil {
		return err
	}
	if seeded {
		return nil
	}

	err = pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, demoSQL); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `INSERT INTO settings (key, value) VALUES ('seeded', 'true')`)
		return err
	})
	if err != nil {
		return err
	}
	slog.Info("demo kontent yozildi")
	return nil
}

// SeedNews — demo yangiliklar. Yangiliklar keyinroq qo'shilgani uchun alohida belgi bilan:
// oldin to'ldirilgan bazaga ham bir marta yoziladi.
func SeedNews(ctx context.Context, pool *pgxpool.Pool) error {
	var seeded bool
	err := pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM settings WHERE key = 'seeded_news')`).Scan(&seeded)
	if err != nil || seeded {
		return err
	}
	err = pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, newsSQL); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `INSERT INTO settings (key, value) VALUES ('seeded_news', 'true')`)
		return err
	})
	if err == nil {
		slog.Info("demo yangiliklar yozildi")
	}
	return err
}
