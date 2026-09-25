// Package middleware — xavfsizlik, autentifikatsiya, kesh va loglash qatlamlari.
package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"uzbekona.dev/backend/internal/apperr"
	"uzbekona.dev/backend/internal/model"
	"uzbekona.dev/backend/internal/service"
)

const adminKey = "admin"

// RequestLogger — har bir so'rovga ID beradi va natijani strukturali log qiladi.
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		id := c.Get(fiber.HeaderXRequestID)
		if id == "" || len(id) > 64 {
			b := make([]byte, 8)
			_, _ = rand.Read(b)
			id = hex.EncodeToString(b)
		}
		c.Set(fiber.HeaderXRequestID, id)

		err := c.Next()

		status := c.Response().StatusCode()
		if err != nil {
			if fe, ok := err.(*fiber.Error); ok {
				status = fe.Code
			} else if ae, ok := err.(*apperr.Error); ok {
				status = ae.Status
			} else {
				status = fiber.StatusInternalServerError
			}
		}

		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}
		attrs := []any{
			"id", id,
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"duration_ms", time.Since(start).Milliseconds(),
			"ip", c.IP(),
		}
		if err != nil && status >= 500 {
			attrs = append(attrs, "error", err.Error())
		}
		slog.Log(c.Context(), level, "http", attrs...)
		return err
	}
}

// JSONBodyLimit — fayl yuklashdan tashqari so'rovlar tanasi hajmini cheklaydi.
func JSONBodyLimit(maxKB int) fiber.Handler {
	limit := maxKB * 1024
	return func(c *fiber.Ctx) error {
		if strings.HasPrefix(c.Get(fiber.HeaderContentType), fiber.MIMEMultipartForm) {
			return c.Next()
		}
		if len(c.Body()) > limit {
			return apperr.New(fiber.StatusRequestEntityTooLarge, "payload_too_large", "So‘rov hajmi juda katta")
		}
		return c.Next()
	}
}

// RequireAdmin — HttpOnly cookie'dagi JWT'ni tekshiradi va adminni kontekstga qo'yadi.
func RequireAdmin(auth *service.AuthService, cookieName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Admin javoblari brauzer yoki proxy keshida saqlanmasin
		c.Set(fiber.HeaderCacheControl, "no-store")
		token := c.Cookies(cookieName)
		if token == "" {
			return apperr.Unauthorized("")
		}
		admin, err := auth.Authenticate(c.UserContext(), token)
		if err != nil {
			return err
		}
		c.Locals(adminKey, admin)
		return c.Next()
	}
}

// CurrentAdmin — RequireAdmin'dan keyin joriy admin.
func CurrentAdmin(c *fiber.Ctx) *model.Admin {
	admin, _ := c.Locals(adminKey).(*model.Admin)
	return admin
}

// OriginGuard — cookie asosidagi sessiya uchun CSRF himoyasi:
// o'zgartiruvchi so'rovlar faqat ruxsat etilgan Origin'dan qabul qilinadi.
func OriginGuard(allowed []string) fiber.Handler {
	set := make(map[string]bool, len(allowed))
	for _, o := range allowed {
		set[strings.TrimRight(o, "/")] = true
	}
	return func(c *fiber.Ctx) error {
		switch c.Method() {
		case fiber.MethodGet, fiber.MethodHead, fiber.MethodOptions:
			return c.Next()
		}
		origin := c.Get(fiber.HeaderOrigin)
		if origin == "" {
			if ref, err := url.Parse(c.Get(fiber.HeaderReferer)); err == nil && ref.Host != "" {
				origin = ref.Scheme + "://" + ref.Host
			}
		}
		// Brauzer bo'lmagan klientlar (curl) Origin yubormaydi; SameSite=Strict cookie
		// esa begona saytdan kelgan so'rovga umuman qo'shilmaydi.
		if origin == "" || set[origin] || origin == c.BaseURL() {
			return c.Next()
		}
		return apperr.Forbidden("So‘rov manbasi ruxsat etilmagan")
	}
}
