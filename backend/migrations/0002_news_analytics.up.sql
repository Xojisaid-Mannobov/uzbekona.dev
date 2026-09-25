-- Yangiliklar (studiya xabarlari) va sayt tashriflari statistikasi

CREATE TABLE news (
    id           BIGSERIAL PRIMARY KEY,
    slug         TEXT        NOT NULL UNIQUE,
    title        TEXT        NOT NULL,
    excerpt      TEXT        NOT NULL DEFAULT '',
    tag          TEXT        NOT NULL DEFAULT '',
    cover_id     BIGINT REFERENCES media (id) ON DELETE SET NULL,
    content      JSONB       NOT NULL DEFAULT '[]',
    status       TEXT        NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    pinned       BOOLEAN     NOT NULL DEFAULT false,
    views        BIGINT      NOT NULL DEFAULT 0,
    seo          JSONB       NOT NULL DEFAULT '{}',
    published_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX news_status_published_idx ON news (status, published_at DESC);
CREATE TRIGGER news_updated_at BEFORE UPDATE ON news FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Har bir sahifa ko'rilishi. IP va User-Agent saqlanmaydi: visitor/session — maxfiy kalit bilan
-- olingan xesh (qaytarib tiklab bo'lmaydi), shuning uchun cookie'siz va shaxsiy ma'lumotsiz.
CREATE TABLE page_views (
    id            BIGSERIAL PRIMARY KEY,
    path          TEXT        NOT NULL,
    visitor       TEXT        NOT NULL,
    session       TEXT        NOT NULL DEFAULT '',
    referrer_host TEXT        NOT NULL DEFAULT '',
    device        TEXT        NOT NULL DEFAULT 'desktop' CHECK (device IN ('desktop', 'mobile', 'tablet')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX page_views_created_idx ON page_views (created_at);
CREATE INDEX page_views_visitor_idx ON page_views (visitor, path, created_at);
