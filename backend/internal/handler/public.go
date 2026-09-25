package handler

import (
	"github.com/gofiber/fiber/v2"

	"uzbekona.dev/backend/internal/model"
)

// GET /api/v1/projects?featured=true&limit=6
func (h *Handler) ListProjects(c *fiber.Ctx) error {
	items, err := h.svc.Projects.PublicList(c.UserContext(), c.QueryBool("featured"), c.QueryInt("limit", 100))
	if err != nil {
		return err
	}
	return ok(c, items)
}

// GET /api/v1/projects/:slug
func (h *Handler) GetProject(c *fiber.Ctx) error {
	item, err := h.svc.Projects.PublicGet(c.UserContext(), c.Params("slug"))
	if err != nil {
		return err
	}
	return ok(c, item)
}

// GET /api/v1/services
func (h *Handler) ListServices(c *fiber.Ctx) error {
	items, err := h.svc.Catalog.PublicList(c.UserContext())
	if err != nil {
		return err
	}
	return ok(c, items)
}

// GET /api/v1/services/:slug
func (h *Handler) GetService(c *fiber.Ctx) error {
	item, err := h.svc.Catalog.PublicGet(c.UserContext(), c.Params("slug"))
	if err != nil {
		return err
	}
	return ok(c, item)
}

// GET /api/v1/team
func (h *Handler) ListTeam(c *fiber.Ctx) error {
	items, err := h.svc.Team.PublicList(c.UserContext())
	if err != nil {
		return err
	}
	return ok(c, items)
}

// GET /api/v1/labs
func (h *Handler) ListLabs(c *fiber.Ctx) error {
	items, err := h.svc.Labs.PublicList(c.UserContext())
	if err != nil {
		return err
	}
	return ok(c, items)
}

// GET /api/v1/articles?category=engineering&page=1&limit=12
func (h *Handler) ListArticles(c *fiber.Ctx) error {
	items, meta, err := h.svc.Articles.PublicList(c.UserContext(), model.ArticleFilter{
		Category: c.Query("category"),
		Page:     c.QueryInt("page", 1),
		Limit:    c.QueryInt("limit", 12),
	})
	if err != nil {
		return err
	}
	return paged(c, items, meta)
}

// GET /api/v1/articles/:slug
func (h *Handler) GetArticle(c *fiber.Ctx) error {
	item, err := h.svc.Articles.PublicGet(c.UserContext(), c.Params("slug"))
	if err != nil {
		return err
	}
	return ok(c, item)
}

// GET /api/v1/article-categories
func (h *Handler) ListCategories(c *fiber.Ctx) error {
	items, err := h.svc.Articles.Categories(c.UserContext())
	if err != nil {
		return err
	}
	return ok(c, items)
}

// GET /api/v1/settings
func (h *Handler) GetSettings(c *fiber.Ctx) error {
	s, err := h.svc.Settings.Get(c.UserContext())
	if err != nil {
		return err
	}
	return ok(c, s)
}

// POST /api/v1/contact
func (h *Handler) SubmitContact(c *fiber.Ctx) error {
	var in model.ContactInput
	if err := h.bind(c, &in); err != nil {
		return err
	}
	if err := h.svc.Contacts.Submit(c.UserContext(), &in, c.IP(), c.Get(fiber.HeaderUserAgent)); err != nil {
		return err
	}
	return created(c, fiber.Map{"message": "So‘rovingiz qabul qilindi. Tez orada bog‘lanamiz."})
}
