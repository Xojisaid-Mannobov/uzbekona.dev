package routes

import (
	"net/http/httptest"
	"testing"
	"time"

	"uzbekona.dev/backend/internal/cache"
	"uzbekona.dev/backend/internal/config"
	"uzbekona.dev/backend/internal/handler"
	"uzbekona.dev/backend/internal/repository"
	"uzbekona.dev/backend/internal/service"
	"uzbekona.dev/backend/internal/validator"
)

func testApp(t *testing.T) *config.Config {
	t.Helper()
	return &config.Config{
		JWTSecret:      "test-secret-test-secret-test-secret-123",
		JWTTTL:         time.Hour,
		CookieName:     "uzb_admin",
		CORSOrigins:    []string{"http://localhost:5173"},
		UploadDir:      t.TempDir(),
		UploadURL:      "/uploads",
		MaxUploadMB:    10,
		MaxJSONBodyKB:  64,
		PublicCacheTTL: time.Minute,
	}
}

// Admin yo'llari autentifikatsiyasiz 401 qaytarishi va hech qachon public keshga tushmasligi kerak.
func TestAdminRoutesRequireAuthAndAreNotCached(t *testing.T) {
	cfg := testApp(t)
	repos := &repository.Repositories{Admins: &repository.AdminRepo{}}
	svc := service.New(cfg, repos)
	app := NewApp(Deps{Config: cfg, Handler: handler.New(cfg, svc, repos, validator.New()), Services: svc, Cache: cache.NewMemory()})

	for _, path := range []string{"/api/v1/admin/auth/me", "/api/v1/admin/dashboard", "/api/v1/admin/settings", "/api/v1/admin/projects"} {
		for i := 0; i < 2; i++ {
			resp, err := app.Test(httptest.NewRequest("GET", path, nil))
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != 401 {
				t.Fatalf("%s: status %d, kutilgan 401", path, resp.StatusCode)
			}
			if resp.Header.Get("X-Cache") != "" {
				t.Fatalf("%s: admin javobi keshlanmoqda", path)
			}
			if resp.Header.Get("Cache-Control") != "no-store" {
				t.Fatalf("%s: Cache-Control no-store emas", path)
			}
		}
	}
}

// Begona saytdan kelgan o'zgartiruvchi so'rov (CSRF) rad etiladi.
func TestOriginGuardBlocksForeignOrigin(t *testing.T) {
	cfg := testApp(t)
	repos := &repository.Repositories{Admins: &repository.AdminRepo{}}
	svc := service.New(cfg, repos)
	app := NewApp(Deps{Config: cfg, Handler: handler.New(cfg, svc, repos, validator.New()), Services: svc, Cache: cache.NewMemory()})

	req := httptest.NewRequest("POST", "/api/v1/admin/auth/login", nil)
	req.Header.Set("Origin", "https://evil.example")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 403 {
		t.Fatalf("status %d, kutilgan 403", resp.StatusCode)
	}
}
