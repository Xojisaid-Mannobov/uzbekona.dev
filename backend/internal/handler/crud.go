package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"uzbekona.dev/backend/internal/model"
)

// Oddiy CRUD resurslar (xizmatlar, jamoa, labs, maqolalar) uchun umumiy handler'lar.
// Har bir resurs o'z servis funksiyalarini beradi — handler'da takroriy kod yo'q.

func listHandler[T any](fn func(context.Context) ([]T, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		items, err := fn(c.UserContext())
		if err != nil {
			return err
		}
		return ok(c, items)
	}
}

func getHandler[T any](fn func(context.Context, int64) (*T, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := paramID(c)
		if err != nil {
			return err
		}
		item, err := fn(c.UserContext(), id)
		if err != nil {
			return err
		}
		return ok(c, item)
	}
}

// saveHandler — POST (id yo'q → yaratish) va PUT /:id (yangilash) uchun.
func saveHandler[In any, T any](h *Handler, fn func(context.Context, int64, *In) (*T, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var id int64
		if c.Params("id") != "" {
			var err error
			if id, err = paramID(c); err != nil {
				return err
			}
		}
		var in In
		if err := h.bind(c, &in); err != nil {
			return err
		}
		item, err := fn(c.UserContext(), id, &in)
		if err != nil {
			return err
		}
		if id == 0 {
			return created(c, item)
		}
		return ok(c, item)
	}
}

func deleteHandler(fn func(context.Context, int64) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := paramID(c)
		if err != nil {
			return err
		}
		if err := fn(c.UserContext(), id); err != nil {
			return err
		}
		return noContent(c)
	}
}

func reorderHandler(h *Handler, fn func(context.Context, []int64) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in model.ReorderInput
		if err := h.bind(c, &in); err != nil {
			return err
		}
		if err := fn(c.UserContext(), in.IDs); err != nil {
			return err
		}
		return noContent(c)
	}
}

// ─── Resurslar ──────────────────────────────────────────────

// Resource — bitta CRUD resursning barcha handler'lari.
type Resource struct {
	List, Get, Save, Delete, Reorder fiber.Handler
}

func (h *Handler) AdminServices() Resource {
	s := h.svc.Catalog
	return Resource{listHandler(s.AdminList), getHandler(s.AdminGet), saveHandler(h, s.Save), deleteHandler(s.Delete), reorderHandler(h, s.Reorder)}
}

func (h *Handler) AdminTeam() Resource {
	s := h.svc.Team
	return Resource{listHandler(s.AdminList), getHandler(s.AdminGet), saveHandler(h, s.Save), deleteHandler(s.Delete), reorderHandler(h, s.Reorder)}
}

func (h *Handler) AdminLabs() Resource {
	s := h.svc.Labs
	return Resource{listHandler(s.AdminList), getHandler(s.AdminGet), saveHandler(h, s.Save), deleteHandler(s.Delete), reorderHandler(h, s.Reorder)}
}

// GET /api/v1/admin/articles?q=&status=&page=
func (h *Handler) AdminListArticles(c *fiber.Ctx) error {
	items, meta, err := h.svc.Articles.AdminList(c.UserContext(), listParams(c))
	if err != nil {
		return err
	}
	return paged(c, items, meta)
}

func (h *Handler) AdminArticles() (get, save, del fiber.Handler) {
	s := h.svc.Articles
	return getHandler(s.AdminGet), saveHandler(h, s.Save), deleteHandler(s.Delete)
}

// POST /api/v1/admin/article-categories, PUT /:id
func (h *Handler) AdminSaveCategory(c *fiber.Ctx) error {
	var id int64
	if c.Params("id") != "" {
		var err error
		if id, err = paramID(c); err != nil {
			return err
		}
	}
	var in model.CategoryInput
	if err := h.bind(c, &in); err != nil {
		return err
	}
	if err := h.svc.Articles.SaveCategory(c.UserContext(), id, &in); err != nil {
		return err
	}
	items, err := h.svc.Articles.Categories(c.UserContext())
	if err != nil {
		return err
	}
	return ok(c, items)
}

func (h *Handler) AdminDeleteCategory() fiber.Handler {
	return deleteHandler(h.svc.Articles.DeleteCategory)
}
