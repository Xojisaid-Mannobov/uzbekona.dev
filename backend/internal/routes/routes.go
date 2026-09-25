// Package routes — Fiber ilovasi, global middleware'lar va barcha /api/v1 marshrutlari.
package routes

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"uzbekona.dev/backend/internal/cache"
	"uzbekona.dev/backend/internal/config"
	"uzbekona.dev/backend/internal/handler"
	"uzbekona.dev/backend/internal/middleware"
	"uzbekona.dev/backend/internal/service"
)

type Deps struct {
	Config   *config.Config
	Handler  *handler.Handler
	Services *service.Services
	Cache    cache.Cache
	Storage  fiber.Storage // rate limit hisoblagichlari (nil → xotira)
}

func NewApp(d Deps) *fiber.App {
	cfg := d.Config

	fiberCfg := fiber.Config{
		AppName:               "uzbekona-api",
		DisableStartupMessage: true,
		ErrorHandler:          handler.ErrorHandler,
		BodyLimit:             (cfg.MaxUploadMB + 5) << 20,
		ReadTimeout:           60 * time.Second,
		WriteTimeout:          60 * time.Second,
		IdleTimeout:           120 * time.Second,
	}
	if cfg.TrustProxy {
		// Nginx orqasida: haqiqiy IP X-Real-IP sarlavhasidan, faqat ichki tarmoqdagi proxy'dan
		fiberCfg.ProxyHeader = "X-Real-IP"
		fiberCfg.EnableTrustedProxyCheck = true
		fiberCfg.TrustedProxies = []string{"127.0.0.1", "::1", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}
	}
	app := fiber.New(fiberCfg)

	app.Use(recover.New())
	app.Use(middleware.RequestLogger())
	app.Use(helmet.New(helmet.Config{
		ContentSecurityPolicy:     "default-src 'none'; img-src 'self'; media-src 'self'; style-src 'unsafe-inline'; frame-ancestors 'none'",
		CrossOriginEmbedderPolicy: "unsafe-none",
		CrossOriginResourcePolicy: "same-site",
		ReferrerPolicy:            "strict-origin-when-cross-origin",
		HSTSMaxAge:                31536000,
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.CORSOrigins, ","),
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,X-Request-ID",
		AllowCredentials: true,
		MaxAge:           600,
	}))
	app.Use(compress.New())

	// Development'da yuklangan fayllarni API o'zi beradi (productionda — Nginx)
	app.Static(cfg.UploadURL, cfg.UploadDir, fiber.Static{MaxAge: 31536000, ByteRange: true})

	h := d.Handler
	api := app.Group("/api/v1", middleware.JSONBodyLimit(cfg.MaxJSONBodyKB))
	api.Get("/health", h.Health)

	// ─── Public API ─────────────────────────────────────────
	// Middleware'lar har bir route'ga aniq biriktiriladi: Fiber'da Group("") + middleware
	// butun /api/v1 prefiksiga (admin yo'llariga ham) qo'llanib ketadi.
	limit := middleware.RateLimit(d.Storage, "api", 300, time.Minute)
	cached := middleware.PublicCache(d.Cache, cfg.PublicCacheTTL)
	get := func(path string, h fiber.Handler) { api.Get(path, limit, cached, h) }

	get("/settings", h.GetSettings)
	get("/projects", h.ListProjects)
	get("/projects/:slug", h.GetProject)
	get("/services", h.ListServices)
	get("/services/:slug", h.GetService)
	get("/team", h.ListTeam)
	get("/labs", h.ListLabs)
	get("/articles", h.ListArticles)
	get("/articles/:slug", h.GetArticle)
	get("/article-categories", h.ListCategories)
	api.Post("/contact", limit, middleware.RateLimit(d.Storage, "contact", 5, 10*time.Minute), h.SubmitContact)

	// ─── Admin API ──────────────────────────────────────────
	admin := api.Group("/admin", middleware.OriginGuard(cfg.CORSOrigins))
	admin.Post("/auth/login", middleware.RateLimit(d.Storage, "login", 10, 15*time.Minute), h.Login)
	admin.Post("/auth/logout", h.Logout)

	secured := admin.Group("",
		middleware.RequireAdmin(d.Services.Auth, cfg.CookieName),
		middleware.RateLimit(d.Storage, "admin", 600, time.Minute),
		middleware.InvalidateCache(d.Cache),
	)
	secured.Get("/auth/me", h.Me)
	secured.Put("/auth/password", h.ChangePassword)
	secured.Get("/dashboard", h.Dashboard)

	secured.Get("/projects", h.AdminListProjects)
	secured.Post("/projects", h.AdminCreateProject)
	secured.Put("/projects/reorder", h.AdminReorderProjects)
	secured.Get("/projects/:id", h.AdminGetProject)
	secured.Put("/projects/:id", h.AdminUpdateProject)
	secured.Patch("/projects/:id/status", h.AdminProjectStatus)
	secured.Patch("/projects/:id/featured", h.AdminProjectFeatured)
	secured.Delete("/projects/:id", h.AdminDeleteProject)

	resource := func(prefix string, r handler.Resource) {
		secured.Get(prefix, r.List)
		secured.Post(prefix, r.Save)
		secured.Put(prefix+"/reorder", r.Reorder)
		secured.Get(prefix+"/:id", r.Get)
		secured.Put(prefix+"/:id", r.Save)
		secured.Delete(prefix+"/:id", r.Delete)
	}
	resource("/services", h.AdminServices())
	resource("/team", h.AdminTeam())
	resource("/labs", h.AdminLabs())

	articleGet, articleSave, articleDelete := h.AdminArticles()
	secured.Get("/articles", h.AdminListArticles)
	secured.Post("/articles", articleSave)
	secured.Get("/articles/:id", articleGet)
	secured.Put("/articles/:id", articleSave)
	secured.Delete("/articles/:id", articleDelete)

	secured.Get("/article-categories", h.ListCategories)
	secured.Post("/article-categories", h.AdminSaveCategory)
	secured.Put("/article-categories/:id", h.AdminSaveCategory)
	secured.Delete("/article-categories/:id", h.AdminDeleteCategory())

	secured.Get("/media", h.ListMedia)
	secured.Post("/media", h.UploadMedia)
	secured.Put("/media/:id", h.UpdateMedia)
	secured.Get("/media/:id/usage", h.MediaUsage)
	secured.Delete("/media/:id", h.DeleteMedia)

	secured.Get("/requests", h.ListRequests)
	secured.Get("/requests/:id", h.GetRequest)
	secured.Patch("/requests/:id", h.UpdateRequest)
	secured.Delete("/requests/:id", h.DeleteRequest)

	secured.Get("/settings", h.GetSettings)
	secured.Put("/settings", h.UpdateSettings)

	secured.Get("/users", h.ListAdmins)
	secured.Post("/users", h.CreateAdmin)
	secured.Delete("/users/:id", h.DeleteAdmin)

	// Noma'lum /api yo'llari uchun JSON 404
	api.Use(func(c *fiber.Ctx) error {
		return fiber.ErrNotFound
	})

	return app
}
