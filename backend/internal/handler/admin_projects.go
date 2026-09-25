package handler

import (
	"github.com/gofiber/fiber/v2"

	"uzbekona.dev/backend/internal/model"
)

// GET /api/v1/admin/projects?q=&status=&page=
func (h *Handler) AdminListProjects(c *fiber.Ctx) error {
	items, meta, err := h.svc.Projects.AdminList(c.UserContext(), listParams(c))
	if err != nil {
		return err
	}
	return paged(c, items, meta)
}

// GET /api/v1/admin/projects/:id
func (h *Handler) AdminGetProject(c *fiber.Ctx) error {
	id, err := paramID(c)
	if err != nil {
		return err
	}
	item, err := h.svc.Projects.AdminGet(c.UserContext(), id)
	if err != nil {
		return err
	}
	return ok(c, item)
}

// POST /api/v1/admin/projects
func (h *Handler) AdminCreateProject(c *fiber.Ctx) error {
	var in model.ProjectInput
	if err := h.bind(c, &in); err != nil {
		return err
	}
	item, err := h.svc.Projects.Create(c.UserContext(), &in)
	if err != nil {
		return err
	}
	return created(c, item)
}

// PUT /api/v1/admin/projects/:id
func (h *Handler) AdminUpdateProject(c *fiber.Ctx) error {
	id, err := paramID(c)
	if err != nil {
		return err
	}
	var in model.ProjectInput
	if err := h.bind(c, &in); err != nil {
		return err
	}
	item, err := h.svc.Projects.Update(c.UserContext(), id, &in)
	if err != nil {
		return err
	}
	return ok(c, item)
}

// PATCH /api/v1/admin/projects/:id/status  — Draft / Publish / Archive
func (h *Handler) AdminProjectStatus(c *fiber.Ctx) error {
	id, err := paramID(c)
	if err != nil {
		return err
	}
	var in model.StatusInput
	if err := h.bind(c, &in); err != nil {
		return err
	}
	if err := h.svc.Projects.SetStatus(c.UserContext(), id, in.Status); err != nil {
		return err
	}
	return noContent(c)
}

// PATCH /api/v1/admin/projects/:id/featured
func (h *Handler) AdminProjectFeatured(c *fiber.Ctx) error {
	id, err := paramID(c)
	if err != nil {
		return err
	}
	var in model.FeaturedInput
	if err := h.bind(c, &in); err != nil {
		return err
	}
	if err := h.svc.Projects.SetFeatured(c.UserContext(), id, in.Featured); err != nil {
		return err
	}
	return noContent(c)
}

// PUT /api/v1/admin/projects/reorder
func (h *Handler) AdminReorderProjects(c *fiber.Ctx) error {
	var in model.ReorderInput
	if err := h.bind(c, &in); err != nil {
		return err
	}
	if err := h.svc.Projects.Reorder(c.UserContext(), in.IDs); err != nil {
		return err
	}
	return noContent(c)
}

// DELETE /api/v1/admin/projects/:id
func (h *Handler) AdminDeleteProject(c *fiber.Ctx) error {
	id, err := paramID(c)
	if err != nil {
		return err
	}
	if err := h.svc.Projects.Delete(c.UserContext(), id); err != nil {
		return err
	}
	return noContent(c)
}
