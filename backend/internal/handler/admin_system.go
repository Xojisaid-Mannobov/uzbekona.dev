package handler

import (
	"github.com/gofiber/fiber/v2"

	"uzbekona.dev/backend/internal/apperr"
	"uzbekona.dev/backend/internal/model"
)

// GET /api/v1/admin/dashboard
func (h *Handler) Dashboard(c *fiber.Ctx) error {
	stats, err := h.svc.Dashboard.Stats(c.UserContext())
	if err != nil {
		return err
	}
	return ok(c, stats)
}

// ─── Media ──────────────────────────────────────────────────

// GET /api/v1/admin/media?kind=image&q=&page=
func (h *Handler) ListMedia(c *fiber.Ctx) error {
	items, meta, err := h.svc.Media.List(c.UserContext(), c.Query("kind"), listParams(c))
	if err != nil {
		return err
	}
	return paged(c, items, meta)
}

// POST /api/v1/admin/media (multipart, "files" maydoni — bir nechta fayl)
func (h *Handler) UploadMedia(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return apperr.BadRequest("Fayl yuborilmadi")
	}
	files := append(form.File["files"], form.File["file"]...)
	if len(files) == 0 {
		return apperr.BadRequest("Fayl yuborilmadi")
	}
	if len(files) > 20 {
		return apperr.BadRequest("Bir martada ko‘pi bilan 20 ta fayl")
	}

	uploaded := []model.Media{}
	failed := []fiber.Map{}
	for _, fh := range files {
		m, err := h.svc.Media.Upload(c.UserContext(), fh)
		if err != nil {
			msg := "Yuklab bo‘lmadi"
			if ae, ok := err.(*apperr.Error); ok {
				msg = ae.Message
			} else {
				return err
			}
			failed = append(failed, fiber.Map{"name": fh.Filename, "message": msg})
			continue
		}
		uploaded = append(uploaded, *m)
	}
	if len(uploaded) == 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":  apperr.New(fiber.StatusUnprocessableEntity, "upload_failed", failed[0]["message"].(string)),
			"failed": failed,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": uploaded, "failed": failed})
}

// PUT /api/v1/admin/media/:id
func (h *Handler) UpdateMedia(c *fiber.Ctx) error {
	id, err := paramID(c)
	if err != nil {
		return err
	}
	var in model.MediaUpdateInput
	if err := h.bind(c, &in); err != nil {
		return err
	}
	m, err := h.svc.Media.UpdateAlt(c.UserContext(), id, in.Alt)
	if err != nil {
		return err
	}
	return ok(c, m)
}

// GET /api/v1/admin/media/:id/usage
func (h *Handler) MediaUsage(c *fiber.Ctx) error {
	id, err := paramID(c)
	if err != nil {
		return err
	}
	n, err := h.svc.Media.Usage(c.UserContext(), id)
	if err != nil {
		return err
	}
	return ok(c, fiber.Map{"count": n})
}

// DELETE /api/v1/admin/media/:id
func (h *Handler) DeleteMedia(c *fiber.Ctx) error {
	return deleteHandler(h.svc.Media.Delete)(c)
}

// ─── So'rovlar (contacts) ───────────────────────────────────

// GET /api/v1/admin/requests?status=new&q=&page=
func (h *Handler) ListRequests(c *fiber.Ctx) error {
	items, meta, err := h.svc.Contacts.List(c.UserContext(), listParams(c))
	if err != nil {
		return err
	}
	return paged(c, items, meta)
}

func (h *Handler) GetRequest(c *fiber.Ctx) error {
	return getHandler(h.svc.Contacts.Get)(c)
}

// PATCH /api/v1/admin/requests/:id
func (h *Handler) UpdateRequest(c *fiber.Ctx) error {
	id, err := paramID(c)
	if err != nil {
		return err
	}
	var in model.ContactUpdateInput
	if err := h.bind(c, &in); err != nil {
		return err
	}
	item, err := h.svc.Contacts.Update(c.UserContext(), id, &in)
	if err != nil {
		return err
	}
	return ok(c, item)
}

func (h *Handler) DeleteRequest(c *fiber.Ctx) error {
	return deleteHandler(h.svc.Contacts.Delete)(c)
}

// ─── Sozlamalar ─────────────────────────────────────────────

// PUT /api/v1/admin/settings
func (h *Handler) UpdateSettings(c *fiber.Ctx) error {
	var in model.SettingsInput
	if err := h.bind(c, &in); err != nil {
		return err
	}
	s, err := h.svc.Settings.Update(c.UserContext(), &in)
	if err != nil {
		return err
	}
	return ok(c, s)
}
