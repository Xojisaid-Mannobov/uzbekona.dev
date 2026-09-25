// Package database PostgreSQL ulanishi, migratsiyalar va seed ma'lumotlari.
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect pgx connection pool yaratadi va bazaga ulanishni tekshiradi.
func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("DATABASE_URL noto'g'ri: %w", err)
	}
	cfg.MaxConns = 20
	cfg.MinConns = 2
	cfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Docker'da postgres konteyneri kechroq tayyor bo'lishi mumkin — bir necha marta urinamiz
	var pingErr error
	for attempt := 1; attempt <= 10; attempt++ {
		pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		pingErr = pool.Ping(pingCtx)
		cancel()
		if pingErr == nil {
			return pool, nil
		}
		time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
	}
	pool.Close()
	return nil, fmt.Errorf("PostgreSQL'ga ulanib bo'lmadi: %w", pingErr)
}
