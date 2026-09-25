// Package model — domen strukturalari va kiruvchi ma'lumotlar (input DTO).
package model

import "encoding/json"

// Status qiymatlari
const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusArchived  = "archived"
)

// MediaRef — boshqa jadvallarga biriktirilgan media faylning qisqa ko'rinishi.
type MediaRef struct {
	ID       int64          `json:"id"`
	URL      string         `json:"url"`
	Kind     string         `json:"kind"`
	Mime     string         `json:"mime"`
	Width    *int           `json:"width"`
	Height   *int           `json:"height"`
	Alt      string         `json:"alt"`
	Variants []MediaVariant `json:"variants"`
}

// MediaVariant — rasmning responsive (srcset) uchun kichraytirilgan nusxasi.
type MediaVariant struct {
	Width int    `json:"width"`
	URL   string `json:"url"`
}

type Metric struct {
	Value string `json:"value" validate:"required,max=24"`
	Label string `json:"label" validate:"required,max=80"`
}

type SEO struct {
	Title       string `json:"title" validate:"max=160"`
	Description string `json:"description" validate:"max=320"`
}

type SocialRef struct {
	Platform string `json:"platform" validate:"required,max=40"`
	URL      string `json:"url" validate:"required,url,max=500"`
}

// Block — content builder bloki (loyiha case-study va maqolalar uchun).
type Block struct {
	ID   int64           `json:"id,omitempty"`
	Type string          `json:"type" validate:"required,oneof=heading text large_text image full_image gallery video stats quote two_columns three_columns technology process before_after"`
	Data json.RawMessage `json:"data" validate:"required"`
}

// ListParams — admin ro'yxatlari uchun umumiy filter/pagination.
type ListParams struct {
	Query  string
	Status string
	Page   int
	Limit  int
}

func (p ListParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

type PageMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

// ReorderInput — drag & drop tartiblash natijasi.
type ReorderInput struct {
	IDs []int64 `json:"ids" validate:"required,min=1,max=500"`
}
