package service

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"

	"uzbekona.dev/backend/internal/model"
	"uzbekona.dev/backend/internal/repository"
)

type ProjectService struct {
	repos *repository.Repositories
}

const projectNotFound = "Loyiha topilmadi"

func (s *ProjectService) PublicList(ctx context.Context, featuredOnly bool, limit int) ([]model.Project, error) {
	if limit < 1 || limit > 100 {
		limit = 100
	}
	return s.repos.Projects.ListPublished(ctx, featuredOnly, limit)
}

// PublicGet — case-study sahifasi: loyiha, bloklar, galereya va keyingi loyiha.
func (s *ProjectService) PublicGet(ctx context.Context, slug string) (*model.ProjectDetail, error) {
	p, err := s.repos.Projects.GetPublishedBySlug(ctx, slug)
	if err != nil {
		return nil, mapErr(err, projectNotFound)
	}
	if err := s.loadRelations(ctx, p); err != nil {
		return nil, err
	}
	next, err := s.repos.Projects.Next(ctx, p)
	if err != nil {
		return nil, err
	}
	return &model.ProjectDetail{Project: p, Next: next}, nil
}

func (s *ProjectService) AdminList(ctx context.Context, params model.ListParams) ([]model.Project, model.PageMeta, error) {
	params = normalizeList(params, 50, 100)
	items, total, err := s.repos.Projects.ListAdmin(ctx, params)
	return items, model.PageMeta{Page: params.Page, Limit: params.Limit, Total: total}, err
}

func (s *ProjectService) AdminGet(ctx context.Context, id int64) (*model.Project, error) {
	p, err := s.repos.Projects.GetByID(ctx, id)
	if err != nil {
		return nil, mapErr(err, projectNotFound)
	}
	if err := s.loadRelations(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProjectService) loadRelations(ctx context.Context, p *model.Project) error {
	var err error
	if p.Blocks, err = s.repos.Projects.Blocks(ctx, p.ID); err != nil {
		return err
	}
	p.Gallery, err = s.repos.Projects.Gallery(ctx, p.ID)
	return err
}

func (s *ProjectService) Create(ctx context.Context, in *model.ProjectInput) (*model.Project, error) {
	return s.save(ctx, 0, in)
}

func (s *ProjectService) Update(ctx context.Context, id int64, in *model.ProjectInput) (*model.Project, error) {
	return s.save(ctx, id, in)
}

// save — loyiha, bloklar va galereyani bitta tranzaksiyada yozadi.
func (s *ProjectService) save(ctx context.Context, id int64, in *model.ProjectInput) (*model.Project, error) {
	s.normalize(in)
	if err := validateBlocks("blocks", in.Blocks); err != nil {
		return nil, err
	}
	slug, err := resolveSlug(ctx, in.Slug, in.Title, id, s.repos.Projects.SlugExists)
	if err != nil {
		return nil, err
	}

	err = s.repos.InTx(ctx, func(tx pgx.Tx) error {
		if id == 0 {
			if id, err = s.repos.Projects.Create(ctx, tx, slug, in); err != nil {
				return err
			}
		} else if err := s.repos.Projects.Update(ctx, tx, id, slug, in); err != nil {
			return err
		}
		if err := s.repos.Projects.ReplaceBlocks(ctx, tx, id, in.Blocks); err != nil {
			return err
		}
		return s.repos.Projects.ReplaceGallery(ctx, tx, id, in.Gallery)
	})
	if err != nil {
		return nil, mapErr(err, projectNotFound)
	}
	return s.AdminGet(ctx, id)
}

func (s *ProjectService) normalize(in *model.ProjectInput) {
	in.Title = strings.TrimSpace(in.Title)
	in.Slug = strings.TrimSpace(in.Slug)
	in.Accent = strings.ToUpper(strings.TrimSpace(in.Accent))
	in.Platforms = trimAll(in.Platforms)
	in.Services = trimAll(in.Services)
	in.Stack = trimAll(in.Stack)
}

func (s *ProjectService) SetStatus(ctx context.Context, id int64, status string) error {
	return mapErr(s.repos.Projects.SetStatus(ctx, id, status), projectNotFound)
}

func (s *ProjectService) SetFeatured(ctx context.Context, id int64, featured bool) error {
	return mapErr(s.repos.Projects.SetFeatured(ctx, id, featured), projectNotFound)
}

func (s *ProjectService) Reorder(ctx context.Context, ids []int64) error {
	return s.repos.Projects.Reorder(ctx, ids)
}

func (s *ProjectService) Delete(ctx context.Context, id int64) error {
	return mapErr(s.repos.Projects.Delete(ctx, id), projectNotFound)
}
