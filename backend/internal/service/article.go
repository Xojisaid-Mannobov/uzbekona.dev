package service

import (
	"context"
	"strings"

	"uzbekona.dev/backend/internal/model"
	"uzbekona.dev/backend/internal/repository"
)

type ArticleService struct {
	repo       *repository.ArticleRepo
	categories *repository.CategoryRepo
}

const (
	articleNotFound  = "Maqola topilmadi"
	categoryNotFound = "Kategoriya topilmadi"
)

func (s *ArticleService) PublicList(ctx context.Context, f model.ArticleFilter) ([]model.Article, model.PageMeta, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 50 {
		f.Limit = 12
	}
	items, total, err := s.repo.ListPublished(ctx, f)
	return items, model.PageMeta{Page: f.Page, Limit: f.Limit, Total: total}, err
}

func (s *ArticleService) PublicGet(ctx context.Context, slug string) (*model.ArticleDetail, error) {
	a, err := s.repo.GetPublishedBySlug(ctx, slug)
	if err != nil {
		return nil, mapErr(err, articleNotFound)
	}
	related, err := s.repo.Related(ctx, a, 3)
	if err != nil {
		return nil, err
	}
	return &model.ArticleDetail{Article: a, Related: related}, nil
}

func (s *ArticleService) AdminList(ctx context.Context, params model.ListParams) ([]model.Article, model.PageMeta, error) {
	params = normalizeList(params, 20, 100)
	items, total, err := s.repo.ListAdmin(ctx, params)
	return items, model.PageMeta{Page: params.Page, Limit: params.Limit, Total: total}, err
}

func (s *ArticleService) AdminGet(ctx context.Context, id int64) (*model.Article, error) {
	a, err := s.repo.GetByID(ctx, id)
	return a, mapErr(err, articleNotFound)
}

func (s *ArticleService) Save(ctx context.Context, id int64, in *model.ArticleInput) (*model.Article, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.AuthorName = strings.TrimSpace(in.AuthorName)
	if err := validateBlocks("content", in.Content); err != nil {
		return nil, err
	}

	slug, err := resolveSlug(ctx, strings.TrimSpace(in.Slug), in.Title, id, s.repo.SlugExists)
	if err != nil {
		return nil, err
	}
	readingTime := ReadingTime(in.Content)

	if id == 0 {
		id, err = s.repo.Create(ctx, slug, readingTime, in)
	} else {
		err = s.repo.Update(ctx, id, slug, readingTime, in)
	}
	if err != nil {
		return nil, mapErr(err, articleNotFound)
	}
	return s.AdminGet(ctx, id)
}

func (s *ArticleService) Delete(ctx context.Context, id int64) error {
	return mapErr(s.repo.Delete(ctx, id), articleNotFound)
}

// ─── Kategoriyalar ──────────────────────────────────────────

func (s *ArticleService) Categories(ctx context.Context) ([]model.ArticleCategory, error) {
	return s.categories.List(ctx)
}

func (s *ArticleService) SaveCategory(ctx context.Context, id int64, in *model.CategoryInput) error {
	in.Name = strings.TrimSpace(in.Name)
	slug, err := resolveSlug(ctx, strings.TrimSpace(in.Slug), in.Name, id, s.categories.SlugExists)
	if err != nil {
		return err
	}
	if id == 0 {
		_, err = s.categories.Create(ctx, slug, in)
	} else {
		err = s.categories.Update(ctx, id, slug, in)
	}
	return mapErr(err, categoryNotFound)
}

func (s *ArticleService) DeleteCategory(ctx context.Context, id int64) error {
	return mapErr(s.categories.Delete(ctx, id), categoryNotFound)
}
