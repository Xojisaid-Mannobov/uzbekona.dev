package service

import (
	"context"
	"strings"

	"uzbekona.dev/backend/internal/model"
	"uzbekona.dev/backend/internal/repository"
)

type ContactService struct {
	repo     *repository.ContactRepo
	notifier Notifier
}

const contactNotFound = "So‘rov topilmadi"

// Submit — public contact form. Honeypot to'ldirilgan bo'lsa (bot), yozuv saqlanmaydi,
// lekin javob bir xil — bot farqni sezmasligi uchun.
func (s *ContactService) Submit(ctx context.Context, in *model.ContactInput, ip, userAgent string) error {
	if strings.TrimSpace(in.Website) != "" {
		return nil
	}

	c := &model.Contact{
		Name:        strings.TrimSpace(in.Name),
		Contact:     strings.TrimSpace(in.Contact),
		Email:       strings.ToLower(strings.TrimSpace(in.Email)),
		ProjectType: strings.TrimSpace(in.ProjectType),
		Budget:      strings.TrimSpace(in.Budget),
		Message:     strings.TrimSpace(in.Message),
		IP:          ip,
		UserAgent:   truncate(userAgent, 400),
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return err
	}

	// Bildirishnoma so'rovni sekinlashtirmasligi uchun fonda yuboriladi
	go s.notifier.NewContact(*c)
	return nil
}

func (s *ContactService) List(ctx context.Context, params model.ListParams) ([]model.Contact, model.PageMeta, error) {
	params = normalizeList(params, 30, 100)
	items, total, err := s.repo.List(ctx, params)
	return items, model.PageMeta{Page: params.Page, Limit: params.Limit, Total: total}, err
}

func (s *ContactService) Get(ctx context.Context, id int64) (*model.Contact, error) {
	c, err := s.repo.GetByID(ctx, id)
	return c, mapErr(err, contactNotFound)
}

func (s *ContactService) Update(ctx context.Context, id int64, in *model.ContactUpdateInput) (*model.Contact, error) {
	if err := s.repo.Update(ctx, id, in); err != nil {
		return nil, mapErr(err, contactNotFound)
	}
	return s.Get(ctx, id)
}

func (s *ContactService) Delete(ctx context.Context, id int64) error {
	return mapErr(s.repo.Delete(ctx, id), contactNotFound)
}
