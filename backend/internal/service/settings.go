package service

import (
	"context"

	"github.com/jackc/pgx/v5"

	"uzbekona.dev/backend/internal/model"
	"uzbekona.dev/backend/internal/repository"
)

type SettingsService struct {
	repos *repository.Repositories
}

// Get — sayt sozlamalari, metrikalar, SEO va ijtimoiy tarmoqlar.
func (s *SettingsService) Get(ctx context.Context) (*model.Settings, error) {
	out := &model.Settings{
		Site:    model.SiteSettings{Name: "Uzbekona.dev"},
		Metrics: []model.Metric{},
	}
	if err := s.repos.Settings.Get(ctx, "site", &out.Site); err != nil {
		return nil, err
	}
	if err := s.repos.Settings.Get(ctx, "metrics", &out.Metrics); err != nil {
		return nil, err
	}
	if err := s.repos.Settings.Get(ctx, "seo", &out.SEO); err != nil {
		return nil, err
	}
	socials, err := s.repos.Settings.Socials(ctx)
	if err != nil {
		return nil, err
	}
	out.Socials = socials
	out.Metrics = nonNilMetrics(out.Metrics)
	return out, nil
}

func (s *SettingsService) Update(ctx context.Context, in *model.SettingsInput) (*model.Settings, error) {
	err := s.repos.InTx(ctx, func(tx pgx.Tx) error {
		if err := s.repos.Settings.Set(ctx, tx, "site", in.Site); err != nil {
			return err
		}
		if err := s.repos.Settings.Set(ctx, tx, "metrics", nonNilMetrics(in.Metrics)); err != nil {
			return err
		}
		if err := s.repos.Settings.Set(ctx, tx, "seo", in.SEO); err != nil {
			return err
		}
		return s.repos.Settings.ReplaceSocials(ctx, tx, in.Socials)
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx)
}

func nonNilMetrics(m []model.Metric) []model.Metric {
	if m == nil {
		return []model.Metric{}
	}
	return m
}

// ─── Dashboard ──────────────────────────────────────────────

type DashboardService struct {
	repos *repository.Repositories
}

func (s *DashboardService) Stats(ctx context.Context) (*model.DashboardStats, error) {
	stats := &model.DashboardStats{}
	if err := s.repos.Settings.DashboardCounts(ctx, stats); err != nil {
		return nil, err
	}
	var err error
	if stats.RecentRequests, err = s.repos.Contacts.Recent(ctx, 5); err != nil {
		return nil, err
	}
	if stats.RecentProjects, err = s.repos.Projects.Recent(ctx, 5); err != nil {
		return nil, err
	}
	if stats.RequestsByDay, err = s.repos.Contacts.ByDay(ctx, 14); err != nil {
		return nil, err
	}
	if stats.ViewsToday, stats.VisitorsToday, err = s.repos.Analytics.Today(ctx); err != nil {
		return nil, err
	}
	return stats, nil
}
