package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"uzbekona.dev/backend/internal/middleware"
	"uzbekona.dev/backend/internal/model"
)

// Sessiya cookie faqat admin API yo'liga yuboriladi.
const adminCookiePath = "/api/v1/admin"

func (h *Handler) setSessionCookie(c *fiber.Ctx, token string, expires time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     h.cfg.CookieName,
		Value:    token,
		Path:     adminCookiePath,
		Domain:   h.cfg.CookieDomain,
		Expires:  expires,
		HTTPOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: fiber.CookieSameSiteStrictMode,
	})
}

// POST /api/v1/admin/auth/login
func (h *Handler) Login(c *fiber.Ctx) error {
	var in model.LoginInput
	if err := h.bind(c, &in); err != nil {
		return err
	}
	token, admin, err := h.svc.Auth.Login(c.UserContext(), &in)
	if err != nil {
		return err
	}
	h.setSessionCookie(c, token, time.Now().Add(h.cfg.JWTTTL))
	return ok(c, admin)
}

// POST /api/v1/admin/auth/logout
func (h *Handler) Logout(c *fiber.Ctx) error {
	h.setSessionCookie(c, "", time.Unix(0, 0))
	return noContent(c)
}

// GET /api/v1/admin/auth/me
func (h *Handler) Me(c *fiber.Ctx) error {
	return ok(c, middleware.CurrentAdmin(c))
}

// PUT /api/v1/admin/auth/password
func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	var in model.PasswordInput
	if err := h.bind(c, &in); err != nil {
		return err
	}
	token, err := h.svc.Auth.ChangePassword(c.UserContext(), middleware.CurrentAdmin(c), &in)
	if err != nil {
		return err
	}
	h.setSessionCookie(c, token, time.Now().Add(h.cfg.JWTTTL))
	return noContent(c)
}

// ─── Adminlar (Users) ───────────────────────────────────────

func (h *Handler) ListAdmins(c *fiber.Ctx) error {
	items, err := h.svc.Auth.ListAdmins(c.UserContext())
	if err != nil {
		return err
	}
	return ok(c, items)
}

func (h *Handler) CreateAdmin(c *fiber.Ctx) error {
	var in model.AdminInput
	if err := h.bind(c, &in); err != nil {
		return err
	}
	admin, err := h.svc.Auth.CreateAdmin(c.UserContext(), &in)
	if err != nil {
		return err
	}
	return created(c, admin)
}

func (h *Handler) DeleteAdmin(c *fiber.Ctx) error {
	id, err := paramID(c)
	if err != nil {
		return err
	}
	if err := h.svc.Auth.DeleteAdmin(c.UserContext(), middleware.CurrentAdmin(c), id); err != nil {
		return err
	}
	return noContent(c)
}
