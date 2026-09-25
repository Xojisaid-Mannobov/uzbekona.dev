// Package service — biznes logika qatlami. Handler faqat HTTP bilan ishlaydi,
// barcha qoidalar (slug, status, fayl qayta ishlash, xavfsizlik) shu yerda.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"unicode"

	"uzbekona.dev/backend/internal/apperr"
	"uzbekona.dev/backend/internal/config"
	"uzbekona.dev/backend/internal/model"
	"uzbekona.dev/backend/internal/repository"
	"uzbekona.dev/backend/pkg/slug"
)

type Services struct {
	Auth      *AuthService
	Projects  *ProjectService
	Catalog   *ServiceCatalog
	Team      *TeamService
	Labs      *LabService
	Articles  *ArticleService
	News      *NewsService
	Analytics *AnalyticsService
	Media     *MediaService
	Contacts  *ContactService
	Settings  *SettingsService
	Dashboard *DashboardService
}

func New(cfg *config.Config, repos *repository.Repositories) *Services {
	return &Services{
		Auth:      &AuthService{cfg: cfg, admins: repos.Admins},
		Projects:  &ProjectService{repos: repos},
		Catalog:   &ServiceCatalog{repo: repos.Services},
		Team:      &TeamService{repo: repos.Team},
		Labs:      &LabService{repo: repos.Labs},
		Articles:  &ArticleService{repo: repos.Articles, categories: repos.Categories},
		News:      &NewsService{repo: repos.News},
		Analytics: newAnalyticsService(cfg, repos),
		Media:     &MediaService{cfg: cfg, repo: repos.Media},
		Contacts:  &ContactService{repo: repos.Contacts, notifier: NewTelegramNotifier(cfg)},
		Settings:  &SettingsService{repos: repos},
		Dashboard: &DashboardService{repos: repos},
	}
}

// mapErr — repository xatolarini foydalanuvchiga tushunarli apperr'ga aylantiradi.
func mapErr(err error, notFoundMsg string) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, repository.ErrNotFound):
		return apperr.NotFound(notFoundMsg)
	case repository.IsUniqueViolation(err):
		return apperr.FieldConflict("slug", "Bu slug allaqachon band")
	case repository.IsForeignKeyViolation(err):
		return apperr.BadRequest("Bog‘langan yozuv (media yoki kategoriya) topilmadi")
	default:
		return err
	}
}

type slugChecker func(ctx context.Context, slug string, excludeID int64) (bool, error)

// resolveSlug — slug berilmasa sarlavhadan yasaydi va band bo'lsa -2, -3 … qo'shadi.
// Admin slugni o'zi kiritgan bo'lsa va u band bo'lsa — aniq xato qaytaradi.
func resolveSlug(ctx context.Context, requested, title string, excludeID int64, exists slugChecker) (string, error) {
	if requested != "" {
		taken, err := exists(ctx, requested, excludeID)
		if err != nil {
			return "", err
		}
		if taken {
			return "", apperr.FieldConflict("slug", "Bu slug allaqachon band")
		}
		return requested, nil
	}

	base := slug.Make(title)
	if base == "" {
		base = "item"
	}
	candidate := base
	for i := 2; i < 1000; i++ {
		taken, err := exists(ctx, candidate, excludeID)
		if err != nil {
			return "", err
		}
		if !taken {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
	return "", apperr.FieldConflict("slug", "Bo‘sh slug topilmadi")
}

const maxBlockBytes = 64 * 1024

// validateBlocks — content builder bloklari: data JSON obyekt va o'lchami chegaralangan.
func validateBlocks(field string, blocks []model.Block) error {
	fields := map[string]string{}
	for i, b := range blocks {
		data := strings.TrimSpace(string(b.Data))
		if !strings.HasPrefix(data, "{") || !json.Valid(b.Data) {
			fields[fmt.Sprintf("%s.%d.data", field, i)] = "Blok ma’lumoti noto‘g‘ri"
		} else if len(b.Data) > maxBlockBytes {
			fields[fmt.Sprintf("%s.%d.data", field, i)] = "Blok juda katta"
		}
	}
	if len(fields) > 0 {
		return apperr.Validation(fields)
	}
	return nil
}

// ReadingTime — bloklardagi barcha matnlar bo'yicha o'qish vaqti (daqiqa, 200 so'z/daqiqa).
func ReadingTime(blocks []model.Block) int {
	words := 0
	for _, b := range blocks {
		var v any
		if err := json.Unmarshal(b.Data, &v); err == nil {
			words += countWords(v)
		}
	}
	minutes := (words + 199) / 200
	if minutes < 1 {
		return 1
	}
	return minutes
}

func countWords(v any) int {
	switch t := v.(type) {
	case string:
		if strings.HasPrefix(t, "http") || strings.HasPrefix(t, "/") {
			return 0
		}
		return len(strings.FieldsFunc(t, unicode.IsSpace))
	case []any:
		n := 0
		for _, item := range t {
			n += countWords(item)
		}
		return n
	case map[string]any:
		n := 0
		for k, item := range t {
			// media obyektlaridagi texnik maydonlar hisobga olinmaydi
			if k == "url" || k == "mime" || k == "kind" || k == "variants" {
				continue
			}
			n += countWords(item)
		}
		return n
	}
	return 0
}

// normalizeList — sahifa va limitni xavfsiz oraliqqa keltiradi.
func normalizeList(p model.ListParams, defLimit, maxLimit int) model.ListParams {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = defLimit
	}
	if p.Limit > maxLimit {
		p.Limit = maxLimit
	}
	p.Query = strings.TrimSpace(p.Query)
	return p
}

func trimAll(items []string) []string {
	out := make([]string, 0, len(items))
	for _, s := range items {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func logErr(msg string, err error) {
	if err != nil {
		slog.Error(msg, "error", err)
	}
}
