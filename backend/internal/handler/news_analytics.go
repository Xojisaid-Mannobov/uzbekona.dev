package handler

import (
	"github.com/gofiber/fiber/v2"

	"uzbekona.dev/backend/internal/model"
)

// ─── Yangiliklar (public) ───────────────────────────────────

// GET /api/v1/news?page=1&limit=12
func (h *Handler) ListNews(c *fiber.Ctx) error {
	items, meta, err := h.svc.News.PublicList(c.UserContext(), c.QueryInt("page", 1), c.QueryInt("limit", 12))
	if err != nil {
		return err
	}
	return paged(c, items, meta)
}

// GET /api/v1/news/:slug
func (h *Handler) GetNews(c *fiber.Ctx) error {
	item, err := h.svc.News.PublicGet(c.UserContext(), c.Params("slug"))
	if err != nil {
		return err
	}
	return ok(c, item)
}

// ─── Yangiliklar (admin) ────────────────────────────────────

func (h *Handler) AdminListNews(c *fiber.Ctx) error {
	items, meta, err := h.svc.News.AdminList(c.UserContext(), listParams(c))
	if err != nil {
		return err
	}
	return paged(c, items, meta)
}

func (h *Handler) AdminNews() (get, save, del fiber.Handler) {
	s := h.svc.News
	return getHandler(s.AdminGet), saveHandler(h, s.Save), deleteHandler(s.Delete)
}

// ─── Statistika ─────────────────────────────────────────────

// POST /api/v1/track — brauzer sahifa ochilganda yuboradi (navigator.sendBeacon).
// Javob doim 204: tracker xatosi foydalanuvchi tajribasiga ta'sir qilmasligi kerak.
func (h *Handler) Track(c *fiber.Ctx) error {
	var in model.TrackInput
	if err := h.bind(c, &in); err != nil {
		return noContent(c)
	}
	if err := h.svc.Analytics.Track(c.UserContext(), &in, c.IP(), c.Get(fiber.HeaderUserAgent), c.Hostname()); err != nil {
		return err
	}
	return noContent(c)
}

// GET /api/v1/admin/analytics?days=30
func (h *Handler) Analytics(c *fiber.Ctx) error {
	report, err := h.svc.Analytics.Report(c.UserContext(), c.QueryInt("days", 30))
	if err != nil {
		return err
	}
	return ok(c, report)
}
