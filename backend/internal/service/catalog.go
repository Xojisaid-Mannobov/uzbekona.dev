package service

import (
	"context"
	"strings"

	"uzbekona.dev/backend/internal/model"
	"uzbekona.dev/backend/internal/repository"
)

// ─── Xizmatlar ──────────────────────────────────────────────

type ServiceCatalog struct {
	repo *repository.ServiceRepo
}

const serviceNotFound = "Xizmat topilmadi"

func (s *ServiceCatalog) PublicList(ctx context.Context) ([]model.Service, error) {
	return s.repo.ListPublished(ctx)
}

func (s *ServiceCatalog) PublicGet(ctx context.Context, slug string) (*model.Service, error) {
	item, err := s.repo.GetPublishedBySlug(ctx, slug)
	return item, mapErr(err, serviceNotFound)
}

func (s *ServiceCatalog) AdminList(ctx context.Context) ([]model.Service, error) {
	return s.repo.ListAll(ctx)
}

func (s *ServiceCatalog) AdminGet(ctx context.Context, id int64) (*model.Service, error) {
	item, err := s.repo.GetByID(ctx, id)
	return item, mapErr(err, serviceNotFound)
}

func (s *ServiceCatalog) Save(ctx context.Context, id int64, in *model.ServiceInput) (*model.Service, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.Features = trimAll(in.Features)
	in.Stack = trimAll(in.Stack)

	slug, err := resolveSlug(ctx, strings.TrimSpace(in.Slug), in.Title, id, s.repo.SlugExists)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		id, err = s.repo.Create(ctx, slug, in)
	} else {
		err = s.repo.Update(ctx, id, slug, in)
	}
	if err != nil {
		return nil, mapErr(err, serviceNotFound)
	}
	return s.AdminGet(ctx, id)
}

func (s *ServiceCatalog) Reorder(ctx context.Context, ids []int64) error {
	return s.repo.Reorder(ctx, ids)
}

func (s *ServiceCatalog) Delete(ctx context.Context, id int64) error {
	return mapErr(s.repo.Delete(ctx, id), serviceNotFound)
}

// ─── Jamoa ──────────────────────────────────────────────────

type TeamService struct {
	repo *repository.TeamRepo
}

const memberNotFound = "Jamoa a’zosi topilmadi"

func (s *TeamService) PublicList(ctx context.Context) ([]model.TeamMember, error) {
	return s.repo.ListPublished(ctx)
}

func (s *TeamService) AdminList(ctx context.Context) ([]model.TeamMember, error) {
	return s.repo.ListAll(ctx)
}

func (s *TeamService) AdminGet(ctx context.Context, id int64) (*model.TeamMember, error) {
	item, err := s.repo.GetByID(ctx, id)
	return item, mapErr(err, memberNotFound)
}

func (s *TeamService) Save(ctx context.Context, id int64, in *model.TeamMemberInput) (*model.TeamMember, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Role = strings.TrimSpace(in.Role)
	var err error
	if id == 0 {
		id, err = s.repo.Create(ctx, in)
	} else {
		err = s.repo.Update(ctx, id, in)
	}
	if err != nil {
		return nil, mapErr(err, memberNotFound)
	}
	return s.AdminGet(ctx, id)
}

func (s *TeamService) Reorder(ctx context.Context, ids []int64) error {
	return s.repo.Reorder(ctx, ids)
}

func (s *TeamService) Delete(ctx context.Context, id int64) error {
	return mapErr(s.repo.Delete(ctx, id), memberNotFound)
}

// ─── Labs ───────────────────────────────────────────────────

type LabService struct {
	repo *repository.LabRepo
}

const labNotFound = "Labs loyihasi topilmadi"

func (s *LabService) PublicList(ctx context.Context) ([]model.Lab, error) {
	return s.repo.ListPublished(ctx)
}

func (s *LabService) AdminList(ctx context.Context) ([]model.Lab, error) {
	return s.repo.ListAll(ctx)
}

func (s *LabService) AdminGet(ctx context.Context, id int64) (*model.Lab, error) {
	item, err := s.repo.GetByID(ctx, id)
	return item, mapErr(err, labNotFound)
}

func (s *LabService) Save(ctx context.Context, id int64, in *model.LabInput) (*model.Lab, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.Stack = trimAll(in.Stack)

	slug, err := resolveSlug(ctx, strings.TrimSpace(in.Slug), in.Title, id, s.repo.SlugExists)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		id, err = s.repo.Create(ctx, slug, in)
	} else {
		err = s.repo.Update(ctx, id, slug, in)
	}
	if err != nil {
		return nil, mapErr(err, labNotFound)
	}
	return s.AdminGet(ctx, id)
}

func (s *LabService) Reorder(ctx context.Context, ids []int64) error {
	return s.repo.Reorder(ctx, ids)
}

func (s *LabService) Delete(ctx context.Context, id int64) error {
	return mapErr(s.repo.Delete(ctx, id), labNotFound)
}
