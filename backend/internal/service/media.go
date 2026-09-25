package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"

	"uzbekona.dev/backend/internal/apperr"
	"uzbekona.dev/backend/internal/config"
	"uzbekona.dev/backend/internal/model"
	"uzbekona.dev/backend/internal/repository"
	"uzbekona.dev/backend/pkg/imageproc"
)

type MediaService struct {
	cfg  *config.Config
	repo *repository.MediaRepo
}

const mediaNotFound = "Fayl topilmadi"

type allowedType struct {
	ext  string
	kind string
}

// Ruxsat etilgan formatlar: PNG, JPG, WEBP, AVIF, SVG, MP4.
// Tur fayl kengaytmasi bo'yicha emas, tarkibi (magic bytes) bo'yicha aniqlanadi.
var allowedTypes = map[string]allowedType{
	"image/png":     {".png", "image"},
	"image/jpeg":    {".jpg", "image"},
	"image/webp":    {".webp", "image"},
	"image/avif":    {".avif", "image"},
	"image/svg+xml": {".svg", "vector"},
	"video/mp4":     {".mp4", "video"},
}

// Rasterlar uchun alohida chegara — juda katta rasmlar xotirani to'ldirmasin
const maxImageBytes = 25 << 20

var unsafeSVG = regexp.MustCompile(`(?i)<script|<foreignobject|javascript:|\son[a-z]+\s*=|<iframe|<embed|<object|xlink:href\s*=\s*["']?\s*data:`)

func (s *MediaService) List(ctx context.Context, kind string, params model.ListParams) ([]model.Media, model.PageMeta, error) {
	params = normalizeList(params, 40, 100)
	items, total, err := s.repo.List(ctx, kind, params)
	return items, model.PageMeta{Page: params.Page, Limit: params.Limit, Total: total}, err
}

// Upload faylni tekshiradi, diskka yozadi, rasmlar uchun responsive variantlar yaratadi.
func (s *MediaService) Upload(ctx context.Context, fh *multipart.FileHeader) (*model.Media, error) {
	maxBytes := int64(s.cfg.MaxUploadMB) << 20
	if fh.Size > maxBytes {
		return nil, apperr.BadRequest(fmt.Sprintf("Fayl hajmi %d MB dan oshmasin", s.cfg.MaxUploadMB))
	}

	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, maxBytes+1))
	if err != nil {
		return nil, err
	}

	mime := mimetype.Detect(data)
	var (
		t     allowedType
		found bool
	)
	for m := mime; m != nil && !found; m = m.Parent() {
		t, found = allowedTypes[m.String()]
	}
	if !found {
		return nil, apperr.BadRequest("Format qo‘llab-quvvatlanmaydi. Ruxsat: PNG, JPG, WEBP, AVIF, SVG, MP4")
	}
	if t.kind == "image" && len(data) > maxImageBytes {
		return nil, apperr.BadRequest("Rasm hajmi 25 MB dan oshmasin")
	}
	if t.kind == "vector" && unsafeSVG.Match(data) {
		return nil, apperr.BadRequest("SVG ichida skript yoki xavfli atributlar bor")
	}

	// Fayllar oylar bo'yicha papkalarga joylanadi: 2026/09/<random>.jpg
	sub := time.Now().Format("2006/01")
	dir := filepath.Join(s.cfg.UploadDir, filepath.FromSlash(sub))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	base := randomName()
	fileName := base + t.ext
	if err := os.WriteFile(filepath.Join(dir, fileName), data, 0o644); err != nil {
		return nil, err
	}

	m := &model.Media{
		MediaRef: model.MediaRef{
			URL:      s.cfg.UploadURL + "/" + path.Join(sub, fileName),
			Kind:     t.kind,
			Mime:     mimeName(t.ext),
			Alt:      altFromName(fh.Filename),
			Variants: []model.MediaVariant{},
		},
		Path:         path.Join(sub, fileName),
		OriginalName: truncate(fh.Filename, 200),
		Size:         int64(len(data)),
	}

	switch t.ext {
	case ".png", ".jpg", ".webp":
		res, err := imageproc.Process(data, dir, base)
		if err != nil {
			_ = os.Remove(filepath.Join(dir, fileName))
			return nil, apperr.BadRequest("Rasm buzilgan yoki o‘qib bo‘lmadi")
		}
		m.Width, m.Height = &res.Width, &res.Height
		for _, v := range res.Variants {
			m.Variants = append(m.Variants, model.MediaVariant{
				Width: v.Width,
				URL:   s.cfg.UploadURL + "/" + path.Join(sub, v.File),
			})
		}
	case ".svg":
		if w, h, ok := svgSize(data); ok {
			m.Width, m.Height = &w, &h
		}
	}

	if err := s.repo.Create(ctx, m); err != nil {
		s.removeFiles(m)
		return nil, err
	}
	return m, nil
}

func (s *MediaService) UpdateAlt(ctx context.Context, id int64, alt string) (*model.Media, error) {
	if err := s.repo.UpdateAlt(ctx, id, strings.TrimSpace(alt)); err != nil {
		return nil, mapErr(err, mediaNotFound)
	}
	m, err := s.repo.GetByID(ctx, id)
	return m, mapErr(err, mediaNotFound)
}

func (s *MediaService) Usage(ctx context.Context, id int64) (int, error) {
	return s.repo.Usage(ctx, id)
}

// Delete — yozuvni va diskdagi barcha variantlarni o'chiradi.
func (s *MediaService) Delete(ctx context.Context, id int64) error {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return mapErr(err, mediaNotFound)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return mapErr(err, mediaNotFound)
	}
	s.removeFiles(m)
	return nil
}

func (s *MediaService) removeFiles(m *model.Media) {
	paths := []string{m.Path}
	for _, v := range m.Variants {
		paths = append(paths, strings.TrimPrefix(v.URL, s.cfg.UploadURL+"/"))
	}
	for _, p := range paths {
		// Path traversal'dan himoya: faqat upload papkasi ichidagi fayllar
		clean := filepath.Clean(filepath.Join(s.cfg.UploadDir, filepath.FromSlash(p)))
		root := filepath.Clean(s.cfg.UploadDir) + string(os.PathSeparator)
		if strings.HasPrefix(clean, root) {
			_ = os.Remove(clean)
		}
	}
}

func randomName() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func mimeName(ext string) string {
	for mime, t := range allowedTypes {
		if t.ext == ext {
			return mime
		}
	}
	return "application/octet-stream"
}

// altFromName — "kuaf-dashboard_v2.png" → "kuaf dashboard v2"
func altFromName(name string) string {
	base := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
	base = strings.NewReplacer("-", " ", "_", " ").Replace(base)
	return truncate(strings.TrimSpace(base), 200)
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}

var (
	svgWidth   = regexp.MustCompile(`(?i)<svg[^>]*\swidth=["']([\d.]+)(px)?["']`)
	svgHeight  = regexp.MustCompile(`(?i)<svg[^>]*\sheight=["']([\d.]+)(px)?["']`)
	svgViewBox = regexp.MustCompile(`(?i)<svg[^>]*\sviewBox=["'][\d.\-]+[\s,]+[\d.\-]+[\s,]+([\d.]+)[\s,]+([\d.]+)["']`)
)

// svgSize — width/height atributlari yoki viewBox'dan o'lchamni oladi.
func svgSize(data []byte) (int, int, bool) {
	head := data
	if len(head) > 4096 {
		head = head[:4096]
	}
	var w, h float64
	if m := svgWidth.FindSubmatch(head); m != nil {
		fmt.Sscanf(string(m[1]), "%g", &w)
	}
	if m := svgHeight.FindSubmatch(head); m != nil {
		fmt.Sscanf(string(m[1]), "%g", &h)
	}
	if (w == 0 || h == 0) && svgViewBox.Match(head) {
		m := svgViewBox.FindSubmatch(head)
		fmt.Sscanf(string(m[1]), "%g", &w)
		fmt.Sscanf(string(m[2]), "%g", &h)
	}
	if w <= 0 || h <= 0 {
		return 0, 0, false
	}
	return int(w), int(h), true
}
