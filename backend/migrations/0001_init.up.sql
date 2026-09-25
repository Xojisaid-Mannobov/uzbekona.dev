-- Uzbekona.dev — boshlang'ich sxema

-- updated_at ustunini avtomatik yangilovchi trigger funksiyasi
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Adminlar (public foydalanuvchilardan butunlay alohida)
CREATE TABLE admins (
    id            BIGSERIAL PRIMARY KEY,
    name          TEXT        NOT NULL,
    email         TEXT        NOT NULL,
    password_hash TEXT        NOT NULL,
    last_login_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX admins_email_unique ON admins (lower(email));

-- Media kutubxonasi
CREATE TABLE media (
    id            BIGSERIAL PRIMARY KEY,
    kind          TEXT        NOT NULL CHECK (kind IN ('image', 'vector', 'video')),
    path          TEXT        NOT NULL,
    url           TEXT        NOT NULL,
    original_name TEXT        NOT NULL,
    mime          TEXT        NOT NULL,
    size          BIGINT      NOT NULL,
    width         INT,
    height        INT,
    variants      JSONB       NOT NULL DEFAULT '[]',
    alt           TEXT        NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX media_created_idx ON media (created_at DESC);

-- Loyihalar
CREATE TABLE projects (
    id                BIGSERIAL PRIMARY KEY,
    slug              TEXT        NOT NULL UNIQUE,
    title             TEXT        NOT NULL,
    tagline           TEXT        NOT NULL DEFAULT '',
    short_description TEXT        NOT NULL DEFAULT '',
    full_description  TEXT        NOT NULL DEFAULT '',
    cover_id          BIGINT REFERENCES media (id) ON DELETE SET NULL,
    year              INT,
    client            TEXT        NOT NULL DEFAULT '',
    industry          TEXT        NOT NULL DEFAULT '',
    platforms         TEXT[]      NOT NULL DEFAULT '{}',
    services          TEXT[]      NOT NULL DEFAULT '{}',
    stack             TEXT[]      NOT NULL DEFAULT '{}',
    metrics           JSONB       NOT NULL DEFAULT '[]',
    live_url          TEXT        NOT NULL DEFAULT '',
    accent            TEXT        NOT NULL DEFAULT '',
    status            TEXT        NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    featured          BOOLEAN     NOT NULL DEFAULT false,
    position          INT         NOT NULL DEFAULT 0,
    seo               JSONB       NOT NULL DEFAULT '{}',
    published_at      TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX projects_status_position_idx ON projects (status, position);

-- Case-study content builder bloklari
CREATE TABLE project_blocks (
    id         BIGSERIAL PRIMARY KEY,
    project_id BIGINT      NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    type       TEXT        NOT NULL,
    position   INT         NOT NULL DEFAULT 0,
    data       JSONB       NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX project_blocks_project_idx ON project_blocks (project_id, position);

-- Loyiha galereyasi
CREATE TABLE project_images (
    id         BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    media_id   BIGINT NOT NULL REFERENCES media (id) ON DELETE CASCADE,
    caption    TEXT   NOT NULL DEFAULT '',
    position   INT    NOT NULL DEFAULT 0
);
CREATE INDEX project_images_project_idx ON project_images (project_id, position);

-- Xizmatlar
CREATE TABLE services (
    id          BIGSERIAL PRIMARY KEY,
    slug        TEXT        NOT NULL UNIQUE,
    title       TEXT        NOT NULL,
    summary     TEXT        NOT NULL DEFAULT '',
    description TEXT        NOT NULL DEFAULT '',
    features    TEXT[]      NOT NULL DEFAULT '{}',
    stack       TEXT[]      NOT NULL DEFAULT '{}',
    preview_id  BIGINT REFERENCES media (id) ON DELETE SET NULL,
    status      TEXT        NOT NULL DEFAULT 'published' CHECK (status IN ('draft', 'published')),
    position    INT         NOT NULL DEFAULT 0,
    seo         JSONB       NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Jamoa a'zolari
CREATE TABLE team_members (
    id           BIGSERIAL PRIMARY KEY,
    name         TEXT        NOT NULL,
    role         TEXT        NOT NULL DEFAULT '',
    bio          TEXT        NOT NULL DEFAULT '',
    photo_id     BIGINT REFERENCES media (id) ON DELETE SET NULL,
    socials      JSONB       NOT NULL DEFAULT '[]',
    is_published BOOLEAN     NOT NULL DEFAULT true,
    position     INT         NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Journal kategoriyalari
CREATE TABLE article_categories (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT        NOT NULL,
    slug       TEXT        NOT NULL UNIQUE,
    position   INT         NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Maqolalar (kontent — bloklar ro'yxati, JSONB)
CREATE TABLE articles (
    id           BIGSERIAL PRIMARY KEY,
    slug         TEXT        NOT NULL UNIQUE,
    title        TEXT        NOT NULL,
    excerpt      TEXT        NOT NULL DEFAULT '',
    cover_id     BIGINT REFERENCES media (id) ON DELETE SET NULL,
    category_id  BIGINT REFERENCES article_categories (id) ON DELETE SET NULL,
    content      JSONB       NOT NULL DEFAULT '[]',
    reading_time INT         NOT NULL DEFAULT 1,
    author_name  TEXT        NOT NULL DEFAULT '',
    status       TEXT        NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    featured     BOOLEAN     NOT NULL DEFAULT false,
    seo          JSONB       NOT NULL DEFAULT '{}',
    published_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX articles_status_published_idx ON articles (status, published_at DESC);

-- Uzbekona Labs (ichki / open-source loyihalar)
CREATE TABLE labs (
    id           BIGSERIAL PRIMARY KEY,
    slug         TEXT        NOT NULL UNIQUE,
    title        TEXT        NOT NULL,
    description  TEXT        NOT NULL DEFAULT '',
    stage        TEXT        NOT NULL DEFAULT 'in_development' CHECK (stage IN ('open_source', 'experimental', 'in_development')),
    url          TEXT        NOT NULL DEFAULT '',
    repo_url     TEXT        NOT NULL DEFAULT '',
    cover_id     BIGINT REFERENCES media (id) ON DELETE SET NULL,
    stack        TEXT[]      NOT NULL DEFAULT '{}',
    is_published BOOLEAN     NOT NULL DEFAULT true,
    position     INT         NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Contact form so'rovlari
CREATE TABLE contacts (
    id           BIGSERIAL PRIMARY KEY,
    name         TEXT        NOT NULL,
    contact      TEXT        NOT NULL DEFAULT '',
    email        TEXT        NOT NULL DEFAULT '',
    project_type TEXT        NOT NULL DEFAULT '',
    budget       TEXT        NOT NULL DEFAULT '',
    message      TEXT        NOT NULL,
    status       TEXT        NOT NULL DEFAULT 'new' CHECK (status IN ('new', 'in_progress', 'done', 'spam')),
    note         TEXT        NOT NULL DEFAULT '',
    ip           TEXT        NOT NULL DEFAULT '',
    user_agent   TEXT        NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX contacts_status_created_idx ON contacts (status, created_at DESC);

-- Sayt sozlamalari (key → JSON)
CREATE TABLE settings (
    key        TEXT PRIMARY KEY,
    value      JSONB       NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Ijtimoiy tarmoq havolalari
CREATE TABLE social_links (
    id         BIGSERIAL PRIMARY KEY,
    platform   TEXT        NOT NULL,
    label      TEXT        NOT NULL DEFAULT '',
    url        TEXT        NOT NULL,
    position   INT         NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER admins_updated_at BEFORE UPDATE ON admins FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER projects_updated_at BEFORE UPDATE ON projects FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER services_updated_at BEFORE UPDATE ON services FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER team_members_updated_at BEFORE UPDATE ON team_members FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER articles_updated_at BEFORE UPDATE ON articles FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER labs_updated_at BEFORE UPDATE ON labs FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER contacts_updated_at BEFORE UPDATE ON contacts FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER settings_updated_at BEFORE UPDATE ON settings FOR EACH ROW EXECUTE FUNCTION set_updated_at();
