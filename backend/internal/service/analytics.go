package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"net/url"
	"regexp"
	"strings"
	"time"

	"uzbekona.dev/backend/internal/config"
	"uzbekona.dev/backend/internal/model"
	"uzbekona.dev/backend/internal/repository"
)

// AnalyticsService — o'z serverimizdagi, cookie'siz tashriflar statistikasi.
// IP va User-Agent saqlanmaydi: ulardan (yoki brauzerdagi tasodifiy ID'dan) maxfiy kalit bilan
// xesh olinadi, xom ma'lumot bazaga tushmaydi.
type AnalyticsService struct {
	key  []byte
	repo *repository.AnalyticsRepo
	news *repository.NewsRepo
}

func newAnalyticsService(cfg *config.Config, repos *repository.Repositories) *AnalyticsService {
	// JWT kalitidan alohida maqsad uchun hosila kalit
	m := hmac.New(sha256.New, []byte(cfg.JWTSecret))
	m.Write([]byte("uzbekona/analytics/v1"))
	return &AnalyticsService{key: m.Sum(nil), repo: repos.Analytics, news: repos.News}
}

var (
	// Faqat saytdagi haqiqiy marshrutlar hisoblanadi — ixtiyoriy yo'llar bilan statistikani "ifloslantirib" bo'lmaydi
	trackablePath = regexp.MustCompile(`^/(projects|services|about|team|journal|news|contact|join)?(/[a-z0-9-]{1,120})?$`)
	newsPath      = regexp.MustCompile(`^/news/([a-z0-9-]{1,120})$`)
	clientID      = regexp.MustCompile(`^[A-Za-z0-9-]{8,64}$`)
	// Telegram ichki brauzeri ("Telegram-Android") haqiqiy foydalanuvchi — faqat "TelegramBot" robot
	botUA    = regexp.MustCompile(`(?i)bot\b|bot/|crawl|spider|slurp|headless|lighthouse|pagespeed|preview|monitor|curl|wget|python|go-http|java/|okhttp|axios|node-fetch|scrapy|facebookexternalhit`)
	mobileUA = regexp.MustCompile(`(?i)mobi|iphone|ipod|windows phone`)
	tabletUA = regexp.MustCompile(`(?i)ipad|tablet|kindle|silk|playbook`)
	android  = regexp.MustCompile(`(?i)android`)
)

func (s *AnalyticsService) hash(parts ...string) string {
	m := hmac.New(sha256.New, s.key)
	m.Write([]byte(strings.Join(parts, "\x00")))
	return hex.EncodeToString(m.Sum(nil))[:32]
}

// NormalizePath — so'rov satri va oxirgi "/" olib tashlanadi; ruxsat etilmagan yo'l uchun "" qaytaradi.
func NormalizePath(p string) string {
	if i := strings.IndexAny(p, "?#"); i >= 0 {
		p = p[:i]
	}
	if len(p) > 1 {
		p = strings.TrimRight(p, "/")
	}
	if !trackablePath.MatchString(p) {
		return ""
	}
	return p
}

// DeviceOf — User-Agent bo'yicha qurilma turi.
func DeviceOf(ua string) string {
	switch {
	// Android planshetlar UA'sida "Mobile" so'zi bo'lmaydi
	case tabletUA.MatchString(ua), android.MatchString(ua) && !mobileUA.MatchString(ua):
		return "tablet"
	case mobileUA.MatchString(ua):
		return "mobile"
	default:
		return "desktop"
	}
}

// IsBot — qidiruv robotlari, monitoring va avtomatlashtirilgan so'rovlar hisobga olinmaydi.
func IsBot(ua string) bool {
	return ua == "" || botUA.MatchString(ua)
}

// referrerHost — kirish sahifasining manbasi: tashqi sayt domeni yoki "direct".
func referrerHost(ref, ownHost string) string {
	u, err := url.Parse(strings.TrimSpace(ref))
	if err != nil || u.Host == "" {
		return "direct"
	}
	host := strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
	own := strings.TrimPrefix(strings.ToLower(strings.Split(ownHost, ":")[0]), "www.")
	if host == own {
		return "direct"
	}
	if len(host) > 120 {
		host = host[:120]
	}
	return host
}

// Track — bitta sahifa ko'rilishini yozadi. Bot yoki noto'g'ri yo'l jimgina e'tiborsiz qoldiriladi.
func (s *AnalyticsService) Track(ctx context.Context, in *model.TrackInput, ip, ua, host string) error {
	if IsBot(ua) {
		return nil
	}
	path := NormalizePath(in.Path)
	if path == "" {
		return nil
	}

	v := &model.PageView{Path: path, Device: DeviceOf(ua)}
	if clientID.MatchString(in.VID) {
		v.Visitor = s.hash("v", in.VID)
	} else {
		// Brauzer ID bermasa — kunlik (ertasi kuni o'zgaradigan) anonim xesh
		v.Visitor = s.hash("d", time.Now().UTC().Format("2006-01-02"), ip, ua)
	}
	if clientID.MatchString(in.SID) {
		v.Session = s.hash("s", in.SID)
	}
	if in.Landing {
		v.ReferrerHost = referrerHost(in.Referrer, host)
	}

	if m := newsPath.FindStringSubmatch(path); m != nil {
		seen, err := s.repo.SeenRecently(ctx, v.Visitor, path)
		if err != nil {
			return err
		}
		if !seen {
			if err := s.news.IncrementViews(ctx, m[1]); err != nil {
				return err
			}
		}
	}
	return s.repo.Insert(ctx, v)
}

// Report — admin statistika sahifasi: davr (7/30/90 kun) bo'yicha barcha ko'rsatkichlar.
func (s *AnalyticsService) Report(ctx context.Context, days int) (*model.Analytics, error) {
	switch days {
	case 7, 30, 90:
	default:
		days = 30
	}
	out := &model.Analytics{Days: days}
	var err error
	if out.Summary, err = s.repo.Summary(ctx, days); err != nil {
		return nil, err
	}
	if out.ByDay, err = s.repo.ByDay(ctx, days); err != nil {
		return nil, err
	}
	if out.TopPages, err = s.repo.TopPages(ctx, days, 10); err != nil {
		return nil, err
	}
	if out.Referrers, err = s.repo.Referrers(ctx, days, 8); err != nil {
		return nil, err
	}
	if out.Devices, err = s.repo.Devices(ctx, days); err != nil {
		return nil, err
	}
	if out.TopNews, err = s.news.TopViewed(ctx, 5); err != nil {
		return nil, err
	}
	return out, nil
}

// Cleanup — 13 oydan eski ko'rishlarni o'chiradi (yillik taqqoslash uchun yetarli).
func (s *AnalyticsService) Cleanup(ctx context.Context) {
	n, err := s.repo.Cleanup(ctx, 400)
	if err != nil {
		slog.Warn("statistika tozalanmadi", "err", err)
		return
	}
	if n > 0 {
		slog.Info("eski statistika tozalandi", "rows", n)
	}
}
