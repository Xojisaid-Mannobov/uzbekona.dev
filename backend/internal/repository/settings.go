package repository

import (
	"context"
	"encoding/json"

	"uzbekona.dev/backend/internal/model"
)

type SettingsRepo struct {
	db DBTX
}

// Get — kalit bo'yicha JSON qiymatni dest ichiga o'qiydi. Kalit bo'lmasa dest o'zgarmaydi.
func (r *SettingsRepo) Get(ctx context.Context, key string, dest any) error {
	var raw []byte
	err := r.db.QueryRow(ctx, `SELECT value FROM settings WHERE key = $1`, key).Scan(&raw)
	if err != nil {
		if notFound(err) == ErrNotFound {
			return nil
		}
		return err
	}
	return json.Unmarshal(raw, dest)
}

func (r *SettingsRepo) Set(ctx context.Context, db DBTX, key string, value any) error {
	_, err := db.Exec(ctx, `INSERT INTO settings (key, value) VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value`, key, value)
	return err
}

func (r *SettingsRepo) Socials(ctx context.Context) ([]model.SocialLink, error) {
	rows, err := r.db.Query(ctx, `SELECT id, platform, label, url, position FROM social_links ORDER BY position, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.SocialLink{}
	for rows.Next() {
		var s model.SocialLink
		if err := rows.Scan(&s.ID, &s.Platform, &s.Label, &s.URL, &s.Position); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *SettingsRepo) ReplaceSocials(ctx context.Context, db DBTX, links []model.SocialLink) error {
	if _, err := db.Exec(ctx, `DELETE FROM social_links`); err != nil {
		return err
	}
	if len(links) == 0 {
		return nil
	}
	rows := make([][]any, len(links))
	for i, l := range links {
		rows[i] = []any{l.Platform, l.Label, l.URL, i + 1}
	}
	_, err := copyRows(ctx, db, "social_links", []string{"platform", "label", "url", "position"}, rows)
	return err
}

// Dashboard — bitta so'rov bilan barcha hisoblagichlar.
func (r *SettingsRepo) DashboardCounts(ctx context.Context, s *model.DashboardStats) error {
	return r.db.QueryRow(ctx, `SELECT
		(SELECT count(*) FROM projects WHERE status <> 'archived'),
		(SELECT count(*) FROM projects WHERE status = 'published'),
		(SELECT count(*) FROM articles),
		(SELECT count(*) FROM news),
		(SELECT count(*) FROM contacts WHERE status = 'new'),
		(SELECT count(*) FROM job_applications WHERE status = 'new'),
		(SELECT count(*) FROM contacts WHERE status <> 'spam'),
		(SELECT count(*) FROM media),
		(SELECT count(*) FROM admins)`,
	).Scan(&s.ActiveProjects, &s.PublishedProjects, &s.Articles, &s.News, &s.IncomingRequests, &s.NewApplications, &s.TotalRequests, &s.Media, &s.Admins)
}
