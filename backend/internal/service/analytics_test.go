package service

import "testing"

func TestNormalizePath(t *testing.T) {
	cases := map[string]string{
		"/":                         "/",
		"/projects/":                "/projects",
		"/projects/kuaf?utm=tg#top": "/projects/kuaf",
		"/news/yangi-sayt":          "/news/yangi-sayt",
		"/admin":                    "",
		"/admin/news":               "",
		"/projects/kuaf/extra":      "",
		"/<script>":                 "",
		"/unknown":                  "",
		"":                          "",
	}
	for in, want := range cases {
		if got := NormalizePath(in); got != want {
			t.Errorf("NormalizePath(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDeviceOfAndBots(t *testing.T) {
	iphone := "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 Mobile/15E148 Safari/604.1"
	androidPhone := "Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 Chrome/126.0 Mobile Safari/537.36"
	androidTab := "Mozilla/5.0 (Linux; Android 13; SM-X700) AppleWebKit/537.36 Chrome/126.0 Safari/537.36"
	ipad := "Mozilla/5.0 (iPad; CPU OS 17_0 like Mac OS X) AppleWebKit/605.1.15 Safari/604.1"
	desktop := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/126.0 Safari/537.36"
	tgApp := "Mozilla/5.0 (Linux; Android 14) AppleWebKit/537.36 Chrome/126.0 Mobile Safari/537.36 Telegram-Android/10.14"

	for ua, want := range map[string]string{iphone: "mobile", androidPhone: "mobile", androidTab: "tablet", ipad: "tablet", desktop: "desktop"} {
		if got := DeviceOf(ua); got != want {
			t.Errorf("DeviceOf(%q) = %q, want %q", ua, got, want)
		}
	}
	for _, ua := range []string{desktop, iphone, tgApp} {
		if IsBot(ua) {
			t.Errorf("IsBot(%q) = true, haqiqiy foydalanuvchi", ua)
		}
	}
	for _, ua := range []string{"", "Googlebot/2.1 (+http://www.google.com/bot.html)", "TelegramBot (like TwitterBot)",
		"Mozilla/5.0 HeadlessChrome/126.0", "curl/8.5.0", "Chrome-Lighthouse"} {
		if !IsBot(ua) {
			t.Errorf("IsBot(%q) = false, robot", ua)
		}
	}
}

func TestReferrerHost(t *testing.T) {
	cases := []struct{ ref, own, want string }{
		{"", "uzbekona.dev", "direct"},
		{"https://www.google.com/search?q=x", "uzbekona.dev", "google.com"},
		{"https://uzbekona.dev/projects", "uzbekona.dev", "direct"},
		{"https://t.me/uzbekona_dev", "localhost:5173", "t.me"},
		{"not a url", "uzbekona.dev", "direct"},
	}
	for _, c := range cases {
		if got := referrerHost(c.ref, c.own); got != c.want {
			t.Errorf("referrerHost(%q, %q) = %q, want %q", c.ref, c.own, got, c.want)
		}
	}
}
