package service

import (
	"context"
	"strings"

	"uzbekona.dev/backend/internal/apperr"
	"uzbekona.dev/backend/internal/model"
	"uzbekona.dev/backend/internal/repository"
)

// CareersService — "Jamoaga qo'shilish": nomzod o'zi haqida yozadi, jamoa ko'rib chiqib bog'lanadi.
type CareersService struct {
	repo     *repository.ApplicationRepo
	notifier Notifier
}

const applicationNotFound = "Ariza topilmadi"

// httpURL — admin panelda havola sifatida ochiladi, shuning uchun faqat http(s) (javascript: va h.k. emas).
func httpURL(field, v string) error {
	if v == "" || strings.HasPrefix(v, "https://") || strings.HasPrefix(v, "http://") {
		return nil
	}
	return apperr.Validation(map[string]string{field: "Havola http:// yoki https:// bilan boshlanishi kerak"})
}

// Submit — public forma. Honeypot to'ldirilgan bo'lsa (bot) — saqlanmaydi, javob esa bir xil.
func (s *CareersService) Submit(ctx context.Context, in *model.ApplicationInput, ip, userAgent string) error {
	if strings.TrimSpace(in.Website) != "" {
		return nil
	}
	a := &model.Application{
		FullName:     strings.TrimSpace(in.FullName),
		Email:        strings.ToLower(strings.TrimSpace(in.Email)),
		Phone:        strings.TrimSpace(in.Phone),
		Telegram:     strings.TrimSpace(in.Telegram),
		Position:     strings.TrimSpace(in.Position),
		Experience:   strings.TrimSpace(in.Experience),
		PortfolioURL: strings.TrimSpace(in.PortfolioURL),
		ResumeURL:    strings.TrimSpace(in.ResumeURL),
		About:        strings.TrimSpace(in.About),
		IP:           ip,
		UserAgent:    truncate(userAgent, 400),
	}
	if err := httpURL("portfolio_url", a.PortfolioURL); err != nil {
		return err
	}
	if err := httpURL("resume_url", a.ResumeURL); err != nil {
		return err
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return err
	}
	go s.notifier.NewApplication(*a)
	return nil
}

func (s *CareersService) List(ctx context.Context, params model.ListParams) ([]model.Application, model.PageMeta, error) {
	params = normalizeList(params, 30, 100)
	items, total, err := s.repo.List(ctx, params)
	return items, model.PageMeta{Page: params.Page, Limit: params.Limit, Total: total}, err
}

func (s *CareersService) Get(ctx context.Context, id int64) (*model.Application, error) {
	a, err := s.repo.GetByID(ctx, id)
	return a, mapErr(err, applicationNotFound)
}

func (s *CareersService) Update(ctx context.Context, id int64, in *model.ApplicationUpdateInput) (*model.Application, error) {
	if err := s.repo.Update(ctx, id, in); err != nil {
		return nil, mapErr(err, applicationNotFound)
	}
	return s.Get(ctx, id)
}

func (s *CareersService) Delete(ctx context.Context, id int64) error {
	return mapErr(s.repo.Delete(ctx, id), applicationNotFound)
}
