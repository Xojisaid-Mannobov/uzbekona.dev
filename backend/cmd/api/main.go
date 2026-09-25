// Uzbekona.dev API serveri.
//
// Ishga tushirish:
//
//	go run ./cmd/api            — server (AUTO_MIGRATE=true bo'lsa migratsiyalar avtomatik)
//	go run ./cmd/api migrate    — faqat migratsiyalarni qo'llash
//	go run ./cmd/api rollback   — oxirgi migratsiyani bekor qilish
//	go run ./cmd/api seed       — demo kontentni yozish (bir marta)
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"uzbekona.dev/backend/internal/cache"
	"uzbekona.dev/backend/internal/config"
	"uzbekona.dev/backend/internal/database"
	"uzbekona.dev/backend/internal/handler"
	"uzbekona.dev/backend/internal/middleware"
	"uzbekona.dev/backend/internal/repository"
	"uzbekona.dev/backend/internal/routes"
	"uzbekona.dev/backend/internal/service"
	"uzbekona.dev/backend/internal/validator"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server to‘xtadi", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	setupLogger(cfg)

	ctx := context.Background()
	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	// CLI buyruqlari
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			return database.Migrate(ctx, pool)
		case "rollback":
			return database.Rollback(ctx, pool)
		case "seed":
			if err := database.Migrate(ctx, pool); err != nil {
				return err
			}
			if err := database.SeedDemo(ctx, pool); err != nil {
				return err
			}
			return database.SeedNews(ctx, pool)
		default:
			return fmt.Errorf("noma’lum buyruq: %s (migrate | rollback | seed)", os.Args[1])
		}
	}

	if cfg.AutoMigrate {
		if err := database.Migrate(ctx, pool); err != nil {
			return err
		}
	}
	if cfg.SeedDemo {
		if err := database.SeedDemo(ctx, pool); err != nil {
			return err
		}
		if err := database.SeedNews(ctx, pool); err != nil {
			return err
		}
	}

	repos := repository.New(pool)
	services := service.New(cfg, repos)
	if err := services.Auth.EnsureAdmin(ctx); err != nil {
		return err
	}
	if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
		return err
	}

	// Statistika: 13 oydan eski ko'rishlar kuniga bir marta tozalanadi
	go func() {
		for {
			services.Analytics.Cleanup(ctx)
			time.Sleep(24 * time.Hour)
		}
	}()

	// Redis ixtiyoriy: bo'lsa kesh va rate limit hisoblagichlari shu yerda saqlanadi
	var (
		store   cache.Cache = cache.NewMemory()
		storage fiber.Storage
	)
	if cfg.RedisURL != "" {
		opts, err := redis.ParseURL(cfg.RedisURL)
		if err != nil {
			return fmt.Errorf("REDIS_URL noto‘g‘ri: %w", err)
		}
		client := redis.NewClient(opts)
		defer client.Close()
		pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		err = client.Ping(pingCtx).Err()
		cancel()
		if err != nil {
			slog.Warn("Redis’ga ulanib bo‘lmadi — xotiradagi kesh ishlatiladi", "error", err)
		} else {
			store = cache.NewRedis(client)
			storage = middleware.NewRedisStorage(client)
			slog.Info("Redis ulandi")
		}
	}

	app := routes.NewApp(routes.Deps{
		Config:   cfg,
		Handler:  handler.New(cfg, services, repos, validator.New()),
		Services: services,
		Cache:    store,
		Storage:  storage,
	})

	// Graceful shutdown: SIGINT/SIGTERM kelganda joriy so'rovlar tugatiladi
	errCh := make(chan error, 1)
	go func() {
		slog.Info("API ishga tushdi", "port", cfg.Port, "env", cfg.AppEnv)
		errCh <- app.Listen(":" + cfg.Port)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-errCh:
		return err
	case <-quit:
		slog.Info("to‘xtatilmoqda…")
		return app.ShutdownWithTimeout(15 * time.Second)
	}
}

func setupLogger(cfg *config.Config) {
	var h slog.Handler
	if cfg.IsProduction() {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	slog.SetDefault(slog.New(h))
}
