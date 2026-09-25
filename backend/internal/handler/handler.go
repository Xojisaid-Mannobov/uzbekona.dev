// Package handler — HTTP qatlami: so'rovni o'qish, tekshirish, servisni chaqirish va javob qaytarish.
// Biznes logika bu yerda emas — service paketida.
package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"uzbekona.dev/backend/internal/apperr"
	"uzbekona.dev/backend/internal/config"
	"uzbekona.dev/backend/internal/model"
	"uzbekona.dev/backend/internal/repository"
	"uzbekona.dev/backend/internal/service"
	"uzbekona.dev/backend/internal/validator"
)

type Handler struct {
	cfg   *config.Config
	svc   *service.Services
	repos *repository.Repositories
	v     *validator.Validator
}

func New(cfg *config.Config, svc *service.Services, repos *repository.Repositories, v *validator.Validator) *Handler {
	return &Handler{cfg: cfg, svc: svc, repos: repos, v: v}
}

// ─── Javob yordamchilari ────────────────────────────────────

func ok(c *fiber.Ctx, data any) error {
	return c.JSON(fiber.Map{"data": data})
}

func created(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": data})
}

func paged(c *fiber.Ctx, data any, meta model.PageMeta) error {
	return c.JSON(fiber.Map{"data": data, "meta": meta})
}

func noContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// bind — JSON tanani o'qiydi va struct teglari bo'yicha tekshiradi.
func (h *Handler) bind(c *fiber.Ctx, dst any) error {
	if !strings.HasPrefix(c.Get(fiber.HeaderContentType), fiber.MIMEApplicationJSON) {
		return apperr.New(fiber.StatusUnsupportedMediaType, "unsupported_media_type", "Content-Type: application/json bo‘lishi kerak")
	}
	if err := json.Unmarshal(c.Body(), dst); err != nil {
		return apperr.BadRequest("JSON formati noto‘g‘ri")
	}
	return h.v.Struct(dst)
}

func paramID(c *fiber.Ctx) (int64, error) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, apperr.BadRequest("ID noto‘g‘ri")
	}
	return id, nil
}

func listParams(c *fiber.Ctx) model.ListParams {
	return model.ListParams{
		Query:  c.Query("q"),
		Status: c.Query("status"),
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 0),
	}
}

// ErrorHandler — barcha xatolarni yagona JSON formatga keltiradi:
// {"error": {"code": "...", "message": "...", "fields": {...}}}
func ErrorHandler(c *fiber.Ctx, err error) error {
	var ae *apperr.Error
	if errors.As(err, &ae) {
		return c.Status(ae.Status).JSON(fiber.Map{"error": ae})
	}

	var fe *fiber.Error
	if errors.As(err, &fe) {
		msg := fe.Message
		switch fe.Code {
		case fiber.StatusNotFound:
			msg = "Manzil topilmadi"
		case fiber.StatusRequestEntityTooLarge:
			msg = "So‘rov hajmi juda katta"
		case fiber.StatusMethodNotAllowed:
			msg = "Metod ruxsat etilmagan"
		}
		return c.Status(fe.Code).JSON(fiber.Map{"error": apperr.New(fe.Code, "http_error", msg)})
	}

	// Kutilmagan xato — tafsilotlar faqat logda, klientga umumiy xabar
	slog.Error("kutilmagan xato", "path", c.Path(), "error", err)
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": apperr.New(fiber.StatusInternalServerError, "internal_error", "Xatolik yuz berdi. Qayta urinib ko‘ring"),
	})
}

// Health — konteyner/nginx uchun holat tekshiruvi.
func (h *Handler) Health(c *fiber.Ctx) error {
	if err := h.repos.Ping(c.UserContext()); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "degraded", "database": "down"})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}
