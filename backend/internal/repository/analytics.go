package repository

import (
	"context"

	"uzbekona.dev/backend/internal/model"
)

// StatsTimezone — "bugun" va kunlik grafik chegaralari shu vaqt mintaqasi bo'yicha hisoblanadi.
const StatsTimezone = "Asia/Tashkent"

type AnalyticsRepo struct {
	db DBTX
}

// periodStart — oxirgi `days` kun boshlanishi ($1 — kunlar, $2 — vaqt mintaqasi), timestamptz.
const periodStart = `(((now() AT TIME ZONE $2)::date - ($1::int - 1))::timestamp AT TIME ZONE $2)`

// visitKey — sessiya (brauzer tabi) bo'yicha tashrif; sessiya yo'q bo'lsa — kunlik tashrif buyuruvchi.
const visitKey = `CASE WHEN session <> '' THEN session ELSE visitor || to_char(created_at, 'YYYYMMDD') END`

func (r *AnalyticsRepo) Insert(ctx context.Context, v *model.PageView) error {
	_, err := r.db.Exec(ctx, `INSERT INTO page_views (path, visitor, session, referrer_host, device)
		VALUES ($1, $2, $3, $4, $5)`, v.Path, v.Visitor, v.Session, v.ReferrerHost, v.Device)
	return err
}

// SeenRecently — shu tashrif buyuruvchi bu sahifani oxirgi 24 soatda ko'rganmi (yangilik ko'rishlari
// hisoblagichi bir odam sahifani qayta-qayta yangilasa oshib ketmasligi uchun).
func (r *AnalyticsRepo) SeenRecently(ctx context.Context, visitor, path string) (bool, error) {
	var seen bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM page_views
		WHERE visitor = $1 AND path = $2 AND created_at >= now() - interval '24 hours')`, visitor, path).Scan(&seen)
	return seen, err
}

func (r *AnalyticsRepo) Summary(ctx context.Context, days int) (model.AnalyticsSummary, error) {
	var s model.AnalyticsSummary
	err := r.db.QueryRow(ctx, `WITH b AS (
			SELECT `+periodStart+` AS cur_start,
				(((now() AT TIME ZONE $2)::date - (2 * $1::int - 1))::timestamp AT TIME ZONE $2) AS prev_start,
				(((now() AT TIME ZONE $2)::date)::timestamp AT TIME ZONE $2) AS today_start
		)
		SELECT
			count(*) FILTER (WHERE p.created_at >= b.cur_start),
			count(DISTINCT `+visitKey+`) FILTER (WHERE p.created_at >= b.cur_start),
			count(DISTINCT p.visitor) FILTER (WHERE p.created_at >= b.cur_start),
			count(*) FILTER (WHERE p.created_at < b.cur_start),
			count(DISTINCT `+visitKey+`) FILTER (WHERE p.created_at < b.cur_start),
			count(DISTINCT p.visitor) FILTER (WHERE p.created_at < b.cur_start),
			count(*) FILTER (WHERE p.created_at >= b.today_start),
			count(DISTINCT p.visitor) FILTER (WHERE p.created_at >= b.today_start),
			count(DISTINCT p.visitor) FILTER (WHERE p.created_at >= now() - interval '5 minutes')
		FROM b LEFT JOIN page_views p ON p.created_at >= b.prev_start`, days, StatsTimezone).Scan(
		&s.Views, &s.Visits, &s.Visitors, &s.PrevViews, &s.PrevVisits, &s.PrevVisitors,
		&s.TodayViews, &s.TodayVisitors, &s.OnlineNow)
	if err != nil {
		return s, err
	}
	err = r.db.QueryRow(ctx, `SELECT count(*), count(DISTINCT visitor) FROM page_views`).
		Scan(&s.AllTimeViews, &s.AllTimeVisitors)
	return s, err
}

func (r *AnalyticsRepo) ByDay(ctx context.Context, days int) ([]model.AnalyticsDay, error) {
	rows, err := r.db.Query(ctx, `SELECT to_char(d.day, 'YYYY-MM-DD'), count(p.id), count(DISTINCT p.visitor)
		FROM generate_series((now() AT TIME ZONE $2)::date - ($1::int - 1), (now() AT TIME ZONE $2)::date, interval '1 day') AS d(day)
		LEFT JOIN page_views p
			ON p.created_at >= (d.day::timestamp AT TIME ZONE $2)
			AND p.created_at < ((d.day + interval '1 day')::timestamp AT TIME ZONE $2)
		GROUP BY d.day ORDER BY d.day`, days, StatsTimezone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.AnalyticsDay{}
	for rows.Next() {
		var d model.AnalyticsDay
		if err := rows.Scan(&d.Day, &d.Views, &d.Visitors); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// grouped — davr ichida bitta ustun bo'yicha guruhlangan reyting. col — faqat kod ichidagi konstanta.
func (r *AnalyticsRepo) grouped(ctx context.Context, col, extraCond string, days, limit int) ([]model.AnalyticsRow, error) {
	rows, err := r.db.Query(ctx, `SELECT `+col+`, count(*), count(DISTINCT visitor)
		FROM page_views WHERE created_at >= `+periodStart+extraCond+`
		GROUP BY `+col+` ORDER BY 3 DESC, 2 DESC LIMIT $3`, days, StatsTimezone, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.AnalyticsRow{}
	for rows.Next() {
		var row model.AnalyticsRow
		if err := rows.Scan(&row.Key, &row.Views, &row.Visitors); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (r *AnalyticsRepo) TopPages(ctx context.Context, days, limit int) ([]model.AnalyticsRow, error) {
	return r.grouped(ctx, "path", "", days, limit)
}

// Referrers — faqat sessiyaning birinchi (kirish) sahifasida manba yoziladi, shuning uchun
// ko'rishlar soni = shu manbadan kelgan tashriflar soni.
func (r *AnalyticsRepo) Referrers(ctx context.Context, days, limit int) ([]model.AnalyticsRow, error) {
	return r.grouped(ctx, "referrer_host", " AND referrer_host <> ''", days, limit)
}

func (r *AnalyticsRepo) Devices(ctx context.Context, days int) ([]model.AnalyticsRow, error) {
	return r.grouped(ctx, "device", "", days, 3)
}

// Today — dashboard kartasi uchun bugungi ko'rishlar va tashrif buyuruvchilar.
func (r *AnalyticsRepo) Today(ctx context.Context) (views, visitors int, err error) {
	err = r.db.QueryRow(ctx, `SELECT count(*), count(DISTINCT visitor) FROM page_views
		WHERE created_at >= (((now() AT TIME ZONE $1)::date)::timestamp AT TIME ZONE $1)`, StatsTimezone).Scan(&views, &visitors)
	return views, visitors, err
}

// Cleanup — saqlash muddatidan eski yozuvlarni o'chiradi.
func (r *AnalyticsRepo) Cleanup(ctx context.Context, keepDays int) (int64, error) {
	tag, err := r.db.Exec(ctx, `DELETE FROM page_views WHERE created_at < now() - make_interval(days => $1)`, keepDays)
	return tag.RowsAffected(), err
}
