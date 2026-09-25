package model

import "time"

// ─── Admins ─────────────────────────────────────────────────

type Admin struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email,max=160"`
	Password string `json:"password" validate:"required,max=200"`
}

type AdminInput struct {
	Name     string `json:"name" validate:"required,max=120"`
	Email    string `json:"email" validate:"required,email,max=160"`
	Password string `json:"password" validate:"required,min=10,max=72"`
}

type PasswordInput struct {
	CurrentPassword string `json:"current_password" validate:"required,max=200"`
	NewPassword     string `json:"new_password" validate:"required,min=10,max=72"`
}

// ─── Media ──────────────────────────────────────────────────

type Media struct {
	MediaRef
	Path         string    `json:"-"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	CreatedAt    time.Time `json:"created_at"`
}

type MediaUpdateInput struct {
	Alt string `json:"alt" validate:"max=300"`
}

// ─── Contacts ───────────────────────────────────────────────

type Contact struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Contact     string    `json:"contact"`
	Email       string    `json:"email"`
	ProjectType string    `json:"project_type"`
	Budget      string    `json:"budget"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
	Note        string    `json:"note"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ContactInput struct {
	Name        string `json:"name" validate:"required,min=2,max=120"`
	Contact     string `json:"contact" validate:"required_without=Email,max=120"`
	Email       string `json:"email" validate:"omitempty,email,max=160"`
	ProjectType string `json:"project_type" validate:"max=80"`
	Budget      string `json:"budget" validate:"max=80"`
	Message     string `json:"message" validate:"required,min=10,max=5000"`
	// Honeypot — botlar to'ldiradi, odamlar ko'rmaydi
	Website string `json:"website"`
}

type ContactUpdateInput struct {
	Status string `json:"status" validate:"required,oneof=new in_progress done spam"`
	Note   string `json:"note" validate:"max=5000"`
}

// ─── Settings ───────────────────────────────────────────────

type SiteSettings struct {
	Name      string `json:"name" validate:"required,max=80"`
	Tagline   string `json:"tagline" validate:"max=160"`
	Email     string `json:"email" validate:"omitempty,email,max=160"`
	Phone     string `json:"phone" validate:"max=40"`
	Telegram  string `json:"telegram" validate:"max=80"`
	Address   string `json:"address" validate:"max=240"`
	Available bool   `json:"available"`
}

type SocialLink struct {
	ID       int64  `json:"id"`
	Platform string `json:"platform" validate:"required,max=40"`
	Label    string `json:"label" validate:"max=60"`
	URL      string `json:"url" validate:"required,url,max=500"`
	Position int    `json:"position"`
}

// Settings — public sayt uchun barcha sozlamalar bitta obyektda.
type Settings struct {
	Site    SiteSettings `json:"site"`
	Metrics []Metric     `json:"metrics"`
	SEO     SEO          `json:"seo"`
	Socials []SocialLink `json:"socials"`
}

type SettingsInput struct {
	Site    SiteSettings `json:"site"`
	Metrics []Metric     `json:"metrics" validate:"max=8,dive"`
	SEO     SEO          `json:"seo"`
	Socials []SocialLink `json:"socials" validate:"max=12,dive"`
}

// ─── Dashboard ──────────────────────────────────────────────

type DashboardStats struct {
	ActiveProjects    int       `json:"active_projects"`
	PublishedProjects int       `json:"published_projects"`
	Articles          int       `json:"articles"`
	News              int       `json:"news"`
	ViewsToday        int       `json:"views_today"`
	VisitorsToday     int       `json:"visitors_today"`
	IncomingRequests  int       `json:"incoming_requests"`
	TotalRequests     int       `json:"total_requests"`
	Media             int       `json:"media"`
	Admins            int       `json:"admins"`
	RecentRequests    []Contact `json:"recent_requests"`
	RecentProjects    []Project `json:"recent_projects"`
	RequestsByDay     []DayStat `json:"requests_by_day"`
}

type DayStat struct {
	Day   string `json:"day"`
	Count int    `json:"count"`
}

// ─── Analytics (tashriflar statistikasi) ────────────────────

// TrackInput — brauzer har sahifa ochilganda yuboradigan ma'lumot.
// vid/sid — brauzerdagi tasodifiy identifikatorlar; serverda faqat xeshi saqlanadi.
type TrackInput struct {
	Path     string `json:"path" validate:"required,max=300"`
	Referrer string `json:"referrer" validate:"max=1000"`
	VID      string `json:"vid" validate:"max=64"`
	SID      string `json:"sid" validate:"max=64"`
	Landing  bool   `json:"landing"`
}

// PageView — repository'ga yoziladigan tayyor yozuv.
type PageView struct {
	Path         string
	Visitor      string
	Session      string
	ReferrerHost string
	Device       string
}

type AnalyticsSummary struct {
	Views           int `json:"views"`
	Visits          int `json:"visits"`
	Visitors        int `json:"visitors"`
	PrevViews       int `json:"prev_views"`
	PrevVisits      int `json:"prev_visits"`
	PrevVisitors    int `json:"prev_visitors"`
	TodayViews      int `json:"today_views"`
	TodayVisitors   int `json:"today_visitors"`
	OnlineNow       int `json:"online_now"`
	AllTimeViews    int `json:"all_time_views"`
	AllTimeVisitors int `json:"all_time_visitors"`
}

type AnalyticsDay struct {
	Day      string `json:"day"`
	Views    int    `json:"views"`
	Visitors int    `json:"visitors"`
}

type AnalyticsRow struct {
	Key      string `json:"key"`
	Views    int    `json:"views"`
	Visitors int    `json:"visitors"`
}

type NewsViews struct {
	ID    int64  `json:"id"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
	Views int64  `json:"views"`
}

type Analytics struct {
	Days      int              `json:"days"`
	Summary   AnalyticsSummary `json:"summary"`
	ByDay     []AnalyticsDay   `json:"by_day"`
	TopPages  []AnalyticsRow   `json:"top_pages"`
	Referrers []AnalyticsRow   `json:"referrers"`
	Devices   []AnalyticsRow   `json:"devices"`
	TopNews   []NewsViews      `json:"top_news"`
}
