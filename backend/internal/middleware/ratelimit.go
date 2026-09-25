package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/redis/go-redis/v9"

	"uzbekona.dev/backend/internal/apperr"
)

// RateLimit — IP bo'yicha so'rovlar sonini cheklaydi.
// Redis bo'lsa hisoblagich barcha instance'lar orasida umumiy bo'ladi.
func RateLimit(storage fiber.Storage, name string, max int, window time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: window,
		Storage:    storage,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "rl:" + name + ":" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return apperr.New(fiber.StatusTooManyRequests, "rate_limited", "Juda ko‘p so‘rov. Birozdan keyin qayta urinib ko‘ring")
		},
	})
}

// redisStorage — fiber.Storage interfeysining go-redis asosidagi implementatsiyasi.
type redisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) fiber.Storage {
	return &redisStorage{client: client}
}

func (s *redisStorage) Get(key string) ([]byte, error) {
	v, err := s.client.Get(context.Background(), key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return v, err
}

func (s *redisStorage) Set(key string, val []byte, exp time.Duration) error {
	return s.client.Set(context.Background(), key, val, exp).Err()
}

func (s *redisStorage) Delete(key string) error {
	return s.client.Del(context.Background(), key).Err()
}

func (s *redisStorage) Reset() error {
	return nil
}

func (s *redisStorage) Close() error {
	return nil
}
