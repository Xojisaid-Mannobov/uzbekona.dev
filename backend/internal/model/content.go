package model

import "time"

// ─── Projects ───────────────────────────────────────────────

type Project struct {
	ID               int64         `json:"id"`
	Slug             string        `json:"slug"`
	Title            string        `json:"title"`
	Tagline          string        `json:"tagline"`
	ShortDescription string        `json:"short_description"`
	FullDescription  string        `json:"full_description"`
	CoverID          *int64        `json:"cover_id"`
	Cover            *MediaRef     `json:"cover"`
	Year             *int          `json:"year"`
	Client           string        `json:"client"`
	Industry         string        `json:"industry"`
	Platforms        []string      `json:"platforms"`
	Services         []string      `json:"services"`
	Stack            []string      `json:"stack"`
	Metrics          []Metric      `json:"metrics"`
	LiveURL          string        `json:"live_url"`
	Accent           string        `json:"accent"`
	Status           string        `json:"status"`
	Featured         bool          `json:"featured"`
	Position         int           `json:"position"`
	SEO              SEO           `json:"seo"`
	PublishedAt      *time.Time    `json:"published_at"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
	Blocks           []Block       `json:"blocks,omitempty"`
	Gallery          []GalleryItem `json:"gallery,omitempty"`
}

type GalleryItem struct {
	MediaID int64     `json:"media_id" validate:"required,min=1"`
	Caption string    `json:"caption" validate:"max=240"`
	Media   *MediaRef `json:"media,omitempty"`
}

// ProjectDetail — public case-study sahifasi uchun (keyingi loyiha bilan).
type ProjectDetail struct {
	*Project
	Next *Project `json:"next"`
}

type ProjectInput struct {
	Title            string        `json:"title" validate:"required,max=160"`
	Slug             string        `json:"slug" validate:"omitempty,max=160,slug"`
	Tagline          string        `json:"tagline" validate:"max=240"`
	ShortDescription string        `json:"short_description" validate:"max=500"`
	FullDescription  string        `json:"full_description" validate:"max=20000"`
	CoverID          *int64        `json:"cover_id"`
	Year             *int          `json:"year" validate:"omitempty,min=1990,max=2100"`
	Client           string        `json:"client" validate:"max=160"`
	Industry         string        `json:"industry" validate:"max=160"`
	Platforms        []string      `json:"platforms" validate:"max=10,dive,required,max=40"`
	Services         []string      `json:"services" validate:"max=20,dive,required,max=80"`
	Stack            []string      `json:"stack" validate:"max=30,dive,required,max=60"`
	Metrics          []Metric      `json:"metrics" validate:"max=8,dive"`
	LiveURL          string        `json:"live_url" validate:"omitempty,url,max=500"`
	Accent           string        `json:"accent" validate:"omitempty,hexcolor"`
	Status           string        `json:"status" validate:"required,oneof=draft published archived"`
	Featured         bool          `json:"featured"`
	SEO              SEO           `json:"seo"`
	Gallery          []GalleryItem `json:"gallery" validate:"max=60,dive"`
	Blocks           []Block       `json:"blocks" validate:"max=200,dive"`
}

type StatusInput struct {
	Status string `json:"status" validate:"required,oneof=draft published archived"`
}

type FeaturedInput struct {
	Featured bool `json:"featured"`
}

// ─── Services ───────────────────────────────────────────────

type Service struct {
	ID          int64     `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Features    []string  `json:"features"`
	Stack       []string  `json:"stack"`
	PreviewID   *int64    `json:"preview_id"`
	Preview     *MediaRef `json:"preview"`
	Status      string    `json:"status"`
	Position    int       `json:"position"`
	SEO         SEO       `json:"seo"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ServiceInput struct {
	Title       string   `json:"title" validate:"required,max=120"`
	Slug        string   `json:"slug" validate:"omitempty,max=120,slug"`
	Summary     string   `json:"summary" validate:"max=400"`
	Description string   `json:"description" validate:"max=10000"`
	Features    []string `json:"features" validate:"max=20,dive,required,max=120"`
	Stack       []string `json:"stack" validate:"max=20,dive,required,max=60"`
	PreviewID   *int64   `json:"preview_id"`
	Status      string   `json:"status" validate:"required,oneof=draft published"`
	SEO         SEO      `json:"seo"`
}

// ─── Team ───────────────────────────────────────────────────

type TeamMember struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Role        string      `json:"role"`
	Bio         string      `json:"bio"`
	PhotoID     *int64      `json:"photo_id"`
	Photo       *MediaRef   `json:"photo"`
	Socials     []SocialRef `json:"socials"`
	IsPublished bool        `json:"is_published"`
	Position    int         `json:"position"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type TeamMemberInput struct {
	Name        string      `json:"name" validate:"required,max=120"`
	Role        string      `json:"role" validate:"max=160"`
	Bio         string      `json:"bio" validate:"max=1000"`
	PhotoID     *int64      `json:"photo_id"`
	Socials     []SocialRef `json:"socials" validate:"max=8,dive"`
	IsPublished bool        `json:"is_published"`
}

// ─── Labs ───────────────────────────────────────────────────

type Lab struct {
	ID          int64     `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Stage       string    `json:"stage"`
	URL         string    `json:"url"`
	RepoURL     string    `json:"repo_url"`
	CoverID     *int64    `json:"cover_id"`
	Cover       *MediaRef `json:"cover"`
	Stack       []string  `json:"stack"`
	IsPublished bool      `json:"is_published"`
	Position    int       `json:"position"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LabInput struct {
	Title       string   `json:"title" validate:"required,max=120"`
	Slug        string   `json:"slug" validate:"omitempty,max=120,slug"`
	Description string   `json:"description" validate:"max=1000"`
	Stage       string   `json:"stage" validate:"required,oneof=open_source experimental in_development"`
	URL         string   `json:"url" validate:"omitempty,url,max=500"`
	RepoURL     string   `json:"repo_url" validate:"omitempty,url,max=500"`
	CoverID     *int64   `json:"cover_id"`
	Stack       []string `json:"stack" validate:"max=20,dive,required,max=60"`
	IsPublished bool     `json:"is_published"`
}

// ─── Articles ───────────────────────────────────────────────

type ArticleCategory struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Position int    `json:"position"`
	Count    int    `json:"count"`
}

type CategoryInput struct {
	Name     string `json:"name" validate:"required,max=80"`
	Slug     string `json:"slug" validate:"omitempty,max=80,slug"`
	Position int    `json:"position" validate:"min=0,max=1000"`
}

type Article struct {
	ID          int64            `json:"id"`
	Slug        string           `json:"slug"`
	Title       string           `json:"title"`
	Excerpt     string           `json:"excerpt"`
	CoverID     *int64           `json:"cover_id"`
	Cover       *MediaRef        `json:"cover"`
	CategoryID  *int64           `json:"category_id"`
	Category    *ArticleCategory `json:"category"`
	Content     []Block          `json:"content,omitempty"`
	ReadingTime int              `json:"reading_time"`
	AuthorName  string           `json:"author_name"`
	Status      string           `json:"status"`
	Featured    bool             `json:"featured"`
	SEO         SEO              `json:"seo"`
	PublishedAt *time.Time       `json:"published_at"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type ArticleDetail struct {
	*Article
	Related []Article `json:"related"`
}

type ArticleInput struct {
	Title       string     `json:"title" validate:"required,max=200"`
	Slug        string     `json:"slug" validate:"omitempty,max=200,slug"`
	Excerpt     string     `json:"excerpt" validate:"max=500"`
	CoverID     *int64     `json:"cover_id"`
	CategoryID  *int64     `json:"category_id"`
	Content     []Block    `json:"content" validate:"max=300,dive"`
	AuthorName  string     `json:"author_name" validate:"max=120"`
	Status      string     `json:"status" validate:"required,oneof=draft published archived"`
	Featured    bool       `json:"featured"`
	SEO         SEO        `json:"seo"`
	PublishedAt *time.Time `json:"published_at"`
}

type ArticleFilter struct {
	Category string
	Page     int
	Limit    int
}
