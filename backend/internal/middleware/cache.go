package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"uzbekona.dev/backend/internal/cache"
)

// PublicCachePrefix — public GET javoblari shu prefiks bilan saqlanadi.
const PublicCachePrefix = "public:"

// PublicCache — public GET javoblarini kesh'dan beradi (Redis yoki xotira).
// Admin har qanday o'zgarish qilganda InvalidateCache butun prefiksni tozalaydi.
func PublicCache(store cache.Cache, ttl time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Himoya qatlami: admin yo'llari hech qachon keshlanmaydi
		if c.Method() != fiber.MethodGet || ttl <= 0 || strings.Contains(c.Path(), "/admin") {
			return c.Next()
		}
		key := PublicCachePrefix + c.OriginalURL()
		if body, ok := store.Get(c.UserContext(), key); ok {
			c.Set("X-Cache", "HIT")
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
			return c.Send(body)
		}

		if err := c.Next(); err != nil {
			return err
		}
		if c.Response().StatusCode() == fiber.StatusOK {
			body := append([]byte(nil), c.Response().Body()...)
			store.Set(c.UserContext(), key, body, ttl)
		}
		c.Set("X-Cache", "MISS")
		return nil
	}
}

// InvalidateCache — admin muvaffaqiyatli o'zgartirish kiritgach public keshni tozalaydi.
func InvalidateCache(store cache.Cache) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if c.Method() == fiber.MethodGet || err != nil {
			return err
		}
		if status := c.Response().StatusCode(); status < 400 {
			store.DeletePrefix(c.UserContext(), PublicCachePrefix)
		}
		return nil
	}
}
