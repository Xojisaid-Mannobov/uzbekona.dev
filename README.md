# Uzbekona.dev

Digital Product Studio sayti va kontent boshqaruv tizimi: web platformalar, mobil ilovalar, Telegram tizimlari va avtomatlashtirish mahsulotlarini namoyish qiluvchi public sayt + shu frontend ichidagi admin panel + Go API.

```
Cloudflare → Nginx ─┬─ /            → Vue 3 (public sayt)
                    ├─ /admin/*     → Vue 3 (admin panel, shu SPA ichida)
                    ├─ /api/v1/*    → Go (Fiber) → PostgreSQL, Redis
                    └─ /uploads/*   → media fayllar (Nginx to'g'ridan-to'g'ri beradi)
```

## Stack

| Qatlam | Texnologiya |
|---|---|
| Frontend | Vue 3, TypeScript, Vite, Vue Router, Pinia, Axios, GSAP, Lenis |
| Backend | Go 1.26, Fiber v2, pgx v5 (PostgreSQL), go-redis, golang-jwt |
| Ma’lumotlar | PostgreSQL 16, Redis 7 (kesh + rate limit, ixtiyoriy) |
| Infratuzilma | Docker Compose, Nginx, Cloudflare |

UI butunlay custom — tayyor UI kutubxonasi yoki admin template ishlatilmagan.

## Loyiha tuzilishi

```
uzbekona.dev/
├── frontend/
│   └── src/
│       ├── components/   # UI, layout, media, content (BlockRenderer)
│       ├── layouts/      # PublicLayout.vue
│       ├── pages/        # Home, Projects, ProjectDetail, Services, About, Team, Journal, Contact …
│       ├── modules/      # home bo'limlari, projects, journal, team, contact
│       ├── composables/  # useAsync, useSeo, v-reveal, useTheme, smooth scroll
│       ├── stores/       # Pinia: sayt sozlamalari
│       ├── router/       # meta.layout orqali public / admin / blank
│       ├── services/     # Axios + public API
│       ├── types/ utils/ styles/ content/
│       └── admin/        # AdminLayout, sahifalar, content builder, media manager
├── backend/
│   ├── cmd/api/main.go
│   ├── internal/         # config, database, handler, middleware, model, repository, service, validator, routes, cache
│   ├── migrations/       # SQL migratsiyalar (binary ichiga embed qilinadi)
│   └── pkg/              # slug, imageproc
├── nginx/nginx.conf
├── docker-compose.yml
└── .env.example
```

Backend arxitekturasi: **Handler → Service → Repository → PostgreSQL**. Handler faqat HTTP (parse, validatsiya, javob), biznes qoidalar service'da, SQL faqat repository'da.

## Tez boshlash (Docker)

```bash
cp .env.example .env
# .env: POSTGRES_PASSWORD, JWT_SECRET (openssl rand -hex 32), ADMIN_EMAIL, ADMIN_PASSWORD ni to'ldiring
# Lokal http://localhost uchun: COOKIE_SECURE=false, CORS_ORIGINS=http://localhost
docker compose up -d --build
```

- Sayt: http://localhost
- Admin: http://localhost/admin (ADMIN_EMAIL / ADMIN_PASSWORD bilan)
- API: http://localhost/api/v1/health

Birinchi ishga tushishda migratsiyalar avtomatik qo‘llanadi, `SEED_DEMO=true` bo‘lsa boshlang‘ich kontent (xizmatlar, kategoriyalar, demo loyihalar) yoziladi.

## Lokal development

Talablar: Go 1.26+, Node 22+, PostgreSQL 16, Redis (ixtiyoriy).

```bash
# 1. Backend
cd backend
cp .env.example .env          # DATABASE_URL va ADMIN_* ni moslang
go run ./cmd/api              # http://localhost:8080

# 2. Frontend (boshqa terminalda)
cd frontend
npm install
npm run dev                   # http://localhost:5173 — /api va /uploads 8080'ga proxy qilinadi
```

Backend buyruqlari:

```bash
go run ./cmd/api migrate      # faqat migratsiyalar
go run ./cmd/api rollback     # oxirgi migratsiyani bekor qilish
go run ./cmd/api seed         # demo kontent (faqat bo'sh bazaga)
go test ./...
```

Frontend buyruqlari: `npm run build` (vue-tsc + vite), `npm run format`, `npm run lint`.

## Vercel’da demo (faqat frontend)

Vercel serverless platforma — Go server, PostgreSQL va diskka fayl yuklash u yerda doimiy ishlamaydi. Shuning uchun Vercel’da **demo rejim** ishlatiladi:

1. https://vercel.com/new → GitHub repo’ni import qiling, **Root Directory: `frontend`** (qolgan sozlamalar `frontend/vercel.json` da).
2. `/api/v1/*` so‘rovlari `frontend/public/demo-api/*.json` statik snapshot’ga rewrite qilinadi — public sayt to‘liq ishlaydi.
3. Kontakt forma va admin panel `api/demo.js` orqali aniq “Demo rejim” xabarini qaytaradi.

Snapshot’ni yangilash (lokal API ishlab turganda): `cd frontend && npm run snapshot`, keyin commit + push.

To‘liq versiya (admin, forma, media) uchun backend’ni VPS’da `docker compose up -d` bilan ishga tushiring.

## Environment o‘zgaruvchilari (backend)

| O‘zgaruvchi | Standart | Izoh |
|---|---|---|
| `APP_ENV` | `development` | `production` da `JWT_SECRET` majburiy (≥32 belgi) |
| `DATABASE_URL` | lokal postgres | pgx connection string |
| `REDIS_URL` | — | bo‘lmasa kesh va rate limit xotirada |
| `JWT_SECRET`, `JWT_TTL` | —, `12h` | admin sessiya tokeni |
| `COOKIE_SECURE` | `false` | HTTPS’da `true` |
| `CORS_ORIGINS` | `http://localhost:5173` | admin panel domenlari (CSRF tekshiruvi ham shu ro‘yxat bo‘yicha) |
| `TRUST_PROXY` | `false` | Nginx orqasida `true` — IP `X-Real-IP` dan olinadi |
| `UPLOAD_DIR`, `UPLOAD_URL`, `MAX_UPLOAD_MB` | `./storage/uploads`, `/uploads`, `100` | media |
| `PUBLIC_CACHE_TTL` | `5m` | public GET javoblari keshi |
| `AUTO_MIGRATE`, `SEED_DEMO` | `true`, `false` | |
| `ADMIN_NAME`, `ADMIN_EMAIL`, `ADMIN_PASSWORD` | — | adminlar jadvali bo‘sh bo‘lsa birinchi admin |
| `TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID` | — | yangi so‘rovlar haqida bildirishnoma (ixtiyoriy) |

## API

Prefiks: `/api/v1`. Javob formati: `{"data": …, "meta": …}`, xato: `{"error": {"code", "message", "fields"}}`.

**Public**

```
GET  /settings                 GET  /team
GET  /projects?featured=true   GET  /labs
GET  /projects/:slug           GET  /articles?category=&page=&limit=
GET  /services                 GET  /articles/:slug
GET  /services/:slug           GET  /article-categories
GET  /news?page=&limit=        GET  /news/:slug
POST /contact                  POST /track   (sahifa ko'rilishi, 204)
POST /careers/apply            (jamoaga qo'shilish arizasi)
GET  /health
```

**Admin** (HttpOnly cookie sessiya; `auth/login` va `auth/logout` dan tashqari hammasi autentifikatsiya talab qiladi)

```
POST /admin/auth/login   POST /admin/auth/logout   GET /admin/auth/me   PUT /admin/auth/password
GET  /admin/dashboard
GET|POST /admin/projects   GET|PUT|DELETE /admin/projects/:id
PATCH /admin/projects/:id/status   PATCH /admin/projects/:id/featured   PUT /admin/projects/reorder
GET|POST /admin/{services,team,labs}   GET|PUT|DELETE /admin/{…}/:id   PUT /admin/{…}/reorder
GET|POST /admin/articles   GET|PUT|DELETE /admin/articles/:id
GET|POST /admin/news       GET|PUT|DELETE /admin/news/:id
GET /admin/analytics?days=7|30|90
GET /admin/applications   PATCH|DELETE /admin/applications/:id
GET|POST /admin/article-categories   PUT|DELETE /admin/article-categories/:id
GET|POST /admin/media   PUT|DELETE /admin/media/:id   GET /admin/media/:id/usage
GET /admin/requests   GET|PATCH|DELETE /admin/requests/:id
GET|PUT /admin/settings
GET|POST /admin/users   DELETE /admin/users/:id
```

## Admin panel

`/admin/login` → Dashboard (faol/e’lon qilingan loyihalar, maqolalar, yangi so‘rovlar, 14 kunlik grafik), Loyihalar (draft/publish/archive, featured, drag & drop tartiblash, galereya, SEO), **content builder** (14 blok: Heading, Text, Large Text, Image, Full-width Image, Gallery, Video, Stats, Quote, Two/Three Columns, Technology, Process, Before/After — qo‘shish, tahrirlash, nusxalash, o‘chirish, tartiblash), Xizmatlar, Jamoa, Labs, Maqolalar (+kategoriyalar, rejalashtirilgan e’lon), Yangiliklar (teg, mahkamlash, ko‘rishlar soni), Statistika, Nomzodlar (jamoaga qo‘shilish arizalari), Media manager, So‘rovlar, Sozlamalar (kontaktlar, metrikalar, SEO, ijtimoiy tarmoqlar, parol), Adminlar.

Saqlanmagan o‘zgarishlar bilan sahifadan chiqishda ogohlantirish, `Ctrl+S` — tezkor saqlash.

## Yangiliklar va statistika

- **Yangiliklar** (`/news`): admin paneldan CRUD, content builder bloklari, teg (rangi teg bo‘yicha), bitta mahkamlangan asosiy yangilik, rejalashtirilgan e’lon. Bosh sahifada — hero ostidagi yuguruvchi yangiliklar lentasi va “Studiyada nimalar bo‘lyapti?” bo‘limi.
- **Statistika** (`/admin/analytics`): ko‘rishlar, tashriflar (sessiyalar), noyob tashrif buyuruvchilar, hozir saytda, kunlik grafik, eng ko‘p ko‘rilgan sahifalar, manbalar (Telegram, Google …), qurilmalar va eng ko‘p o‘qilgan yangiliklar; 7/30/90 kun va oldingi davrga nisbatan o‘zgarish.
- **Maxfiylik**: tashqi analytics skriptlari va cookie yo‘q. Brauzer tasodifiy ID’ni localStorage/sessionStorage’da saqlaydi, server faqat uning HMAC-xeshini yozadi; IP va User-Agent bazaga tushmaydi. Robotlar, admin sahifalari, 404 va “Do Not Track” yoqilgan brauzerlar hisobga olinmaydi. 400 kundan eski yozuvlar kuniga bir marta tozalanadi.
- Kunlar `Asia/Tashkent` vaqti bo‘yicha hisoblanadi.

## Jamoaga qo‘shilish

`/join` sahifasida nomzod o‘zi haqida yozadi (ism, email, telefon/Telegram, yo‘nalish, tajriba, portfolio va rezyume havolalari, “o‘zim haqimda” — kamida 50 belgi, rozilik belgisi). Ariza admin paneldagi **Nomzodlar** bo‘limiga tushadi: holat (Yangi → Ko‘rib chiqilmoqda → Suhbat → Qabul qilindi / Rad etildi), ichki izoh, “Javob yozish” (email). `TELEGRAM_BOT_TOKEN` berilgan bo‘lsa, yangi ariza Telegram’ga ham yuboriladi. Honeypot va rate limit (30 daqiqada 3 ta), havolalar faqat `http(s)://`.

## Muqova o‘lchami

Yangilik va maqola muqovasi uchun admin kartadagi nisbatni tanlaydi (asl, 16:9, 4:3, 1:1, 3:4, 21:9) va kesilganda qaysi qism saqlanishini (yuqori/markaz/past) — yonida jonli ko‘rinish. Matn ichidagi “Rasm” blokida o‘lcham (kichik / o‘rta / keng) va nisbat tanlanadi.

## Media

PNG, JPG, WEBP, AVIF, SVG, MP4. Tur fayl tarkibi (magic bytes) bo‘yicha aniqlanadi. Raster rasmlardan 640/1024/1600/2400px variantlar yaratiladi (frontend `srcset` + lazy loading bilan ishlatadi). SVG ichida skript/`on*` atributlari bo‘lsa rad etiladi, Nginx `/uploads` uchun qat’iy CSP qo‘yadi.

## Xavfsizlik

- Admin sessiya: JWT **HttpOnly + SameSite=Strict** cookie (`/api/v1/admin` yo‘li bilan cheklangan); parol o‘zgarsa eski tokenlar bekor bo‘ladi; o‘chirilgan admin darhol chiqariladi.
- CSRF: admin’dagi o‘zgartiruvchi so‘rovlar faqat `CORS_ORIGINS` dagi Origin’dan.
- Rate limiting (Redis yoki xotira): public API, login (10/15 daq), contact (5/10 daq); Nginx’da qo‘shimcha `limit_req`.
- Validatsiya (go-playground/validator, o‘zbekcha xabarlar), parametrlangan SQL, `ILIKE` wildcard ekranlash.
- XSS: kontent faqat text interpolation bilan chiqariladi (`v-html` yo‘q); CSP, `X-Frame-Options`, HSTS va boshqa sarlavhalar (Fiber helmet + Nginx).
- Request size limit: JSON uchun 1 MB, upload uchun `MAX_UPLOAD_MB`.
- Contact formada honeypot; bcrypt (cost 12), timing-safe login.
- Strukturali loglar (`slog`, productionda JSON) va har bir so‘rovga `X-Request-ID`.

## Performance

Lokal production build (Lighthouse, mobil emulyatsiya): bosh sahifa — Performance 94, Accessibility 96, Best Practices 100, SEO 100; case study va maqola sahifalari — Performance 97–98. Hero sof CSS bilan animatsiya qilinadi, GSAP/Lenis alohida chunk’da kechiktirib yuklanadi; sahifalar lazy-load; asset’lar immutable kesh.

## Almashtirilishi kerak bo‘lgan boshlang‘ich kontent

Seed ma’lumotlari admin paneldan tahrirlanadi, lekin quyidagilar **real ma’lumot bilan almashtirilishi shart**:

- **Metrikalar** (75K+, 12+, 1M+, 99.9%) — TZ’dagi misol qiymatlar. Sozlamalar → Metrikalar.
- **Kontaktlar va ijtimoiy tarmoqlar** (`hello@uzbekona.dev`, `@uzbekona_dev`, GitHub/Instagram/LinkedIn havolalari).
- **OTMEDU, MyLearn** loyihalari tavsiflari va barcha loyihalar uchun real screenshot’lar (cover yuklanmaguncha loyiha accent rangidagi interfeys eskizi ko‘rsatiladi).
- Jamoa a’zolari va rasmlari.
- **Demo yangiliklar** (4 ta: yangi sayt, KUAF, Telekit, vakansiya) — admin paneldagi Yangiliklar bo‘limida tahrirlang yoki o‘chiring.
