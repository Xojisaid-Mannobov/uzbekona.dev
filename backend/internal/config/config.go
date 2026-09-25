// Package config ilova sozlamalarini environment o'zgaruvchilaridan o'qiydi.
package config

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv string
	Port   string

	DatabaseURL string
	RedisURL    string

	JWTSecret    string
	JWTTTL       time.Duration
	CookieName   string
	CookieSecure bool
	CookieDomain string

	CORSOrigins []string
	TrustProxy  bool

	UploadDir      string
	UploadURL      string
	MaxUploadMB    int
	MaxJSONBodyKB  int
	PublicCacheTTL time.Duration

	AutoMigrate bool
	SeedDemo    bool

	AdminName     string
	AdminEmail    string
	AdminPassword string

	// Yangi so'rovlar haqida Telegram bildirishnomasi (ixtiyoriy)
	TelegramBotToken string
	TelegramChatID   string
}

// IsProduction — production rejimida xavfsizlik talablari qattiqroq.
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

// Load .env faylni (bo'lsa) o'qiydi va Config ni to'ldiradi.
func Load() (*Config, error) {
	loadDotEnv(".env")

	cfg := &Config{
		AppEnv:         env("APP_ENV", "development"),
		Port:           env("PORT", "8080"),
		DatabaseURL:    env("DATABASE_URL", "postgres://uzbekona:uzbekona@localhost:5432/uzbekona?sslmode=disable"),
		RedisURL:       env("REDIS_URL", ""),
		JWTSecret:      env("JWT_SECRET", ""),
		JWTTTL:         envDuration("JWT_TTL", 12*time.Hour),
		CookieName:     env("COOKIE_NAME", "uzb_admin"),
		CookieSecure:   envBool("COOKIE_SECURE", false),
		CookieDomain:   env("COOKIE_DOMAIN", ""),
		CORSOrigins:    envList("CORS_ORIGINS", "http://localhost:5173"),
		TrustProxy:     envBool("TRUST_PROXY", false),
		UploadDir:      env("UPLOAD_DIR", "./storage/uploads"),
		UploadURL:      strings.TrimRight(env("UPLOAD_URL", "/uploads"), "/"),
		MaxUploadMB:    envInt("MAX_UPLOAD_MB", 100),
		MaxJSONBodyKB:  envInt("MAX_JSON_BODY_KB", 1024),
		PublicCacheTTL: envDuration("PUBLIC_CACHE_TTL", 5*time.Minute),
		AutoMigrate:    envBool("AUTO_MIGRATE", true),
		SeedDemo:       envBool("SEED_DEMO", false),
		AdminName:      env("ADMIN_NAME", "Administrator"),
		AdminEmail:     env("ADMIN_EMAIL", ""),
		AdminPassword:  env("ADMIN_PASSWORD", ""),

		TelegramBotToken: env("TELEGRAM_BOT_TOKEN", ""),
		TelegramChatID:   env("TELEGRAM_CHAT_ID", ""),
	}

	if cfg.JWTSecret == "" {
		if cfg.IsProduction() {
			return nil, errors.New("JWT_SECRET production rejimida majburiy")
		}
		// Development uchun vaqtinchalik kalit — productionda ishlatilmaydi
		cfg.JWTSecret = "dev-only-insecure-secret-change-me-please-32b"
	}
	if cfg.IsProduction() && len(cfg.JWTSecret) < 32 {
		return nil, errors.New("JWT_SECRET kamida 32 belgidan iborat bo'lishi kerak")
	}

	return cfg, nil
}

func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v)
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	v, err := strconv.ParseBool(env(key, strconv.FormatBool(fallback)))
	if err != nil {
		return fallback
	}
	return v
}

func envInt(key string, fallback int) int {
	v, err := strconv.Atoi(env(key, strconv.Itoa(fallback)))
	if err != nil {
		return fallback
	}
	return v
}

func envDuration(key string, fallback time.Duration) time.Duration {
	v, err := time.ParseDuration(env(key, fallback.String()))
	if err != nil {
		return fallback
	}
	return v
}

func envList(key, fallback string) []string {
	var out []string
	for _, part := range strings.Split(env(key, fallback), ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// loadDotEnv oddiy KEY=VALUE formatidagi faylni o'qiydi.
// Mavjud environment o'zgaruvchilari ustidan yozilmaydi.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}
