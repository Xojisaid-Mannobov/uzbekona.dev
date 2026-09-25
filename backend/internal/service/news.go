package service

import (
	"context"
	"strings"

	"uzbekona.dev/backend/internal/model"
	"uzbekona.dev/backend/internal/repository"
)

type NewsService struct {
	repo *repository.NewsRepo
}

const newsNotFound = "Yangilik topilmadi"

func (s *NewsService) PublicList(ctx context.Context, page, limit int) ([]model.News, model.PageMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 12
	}
	items, total, err := s.repo.ListPublished(ctx, page, limit)
	return items, model.PageMeta{Page: page, Limit: limit, Total: total}, err
}

func (s *NewsService) PublicGet(ctx context.Context, slug string) (*model.NewsDetail, error) {
	n, err := s.repo.GetPublishedBySlug(ctx, slug)
	if err != nil {
		return nil, mapErr(err, newsNotFound)
	}
	related, err := s.repo.Latest(ctx, n.ID, 3)
	if err != nil {
		return nil, err
	}
	return &model.NewsDetail{News: n, Related: related}, nil
}

func (s *NewsService) AdminList(ctx context.Context, params model.ListParams) ([]model.News, model.PageMeta, error) {
	params = normalizeList(params, 20, 100)
	items, total, err := s.repo.ListAdmin(ctx, params)
	return items, model.PageMeta{Page: params.Page, Limit: params.Limit, Total: total}, err
}

func (s *NewsService) AdminGet(ctx context.Context, id int64) (*model.News, error) {
	n, err := s.repo.GetByID(ctx, id)
	return n, mapErr(err, newsNotFound)
}

func (s *NewsService) Save(ctx context.Context, id int64, in *model.NewsInput) (*model.News, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.Tag = strings.TrimSpace(in.Tag)
	if err := validateBlocks("content", in.Content); err != nil {
		return nil, err
	}

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
		return nil, mapErr(err, newsNotFound)
	}
	// Mahkamlangan yangilik bittagina bo'ladi — yangisi eskisining o'rnini egallaydi
	if in.Pinned {
		if err := s.repo.UnpinOthers(ctx, id); err != nil {
			return nil, err
		}
	}
	return s.AdminGet(ctx, id)
}

func (s *NewsService) Delete(ctx context.Context, id int64) error {
	return mapErr(s.repo.Delete(ctx, id), newsNotFound)
}
