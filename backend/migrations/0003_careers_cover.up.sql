-- Jamoaga qo'shilish arizalari va muqova rasmining ko'rinish o'lchami

CREATE TABLE job_applications (
    id            BIGSERIAL PRIMARY KEY,
    full_name     TEXT        NOT NULL,
    email         TEXT        NOT NULL,
    phone         TEXT        NOT NULL DEFAULT '',
    telegram      TEXT        NOT NULL DEFAULT '',
    position      TEXT        NOT NULL,
    experience    TEXT        NOT NULL DEFAULT '',
    portfolio_url TEXT        NOT NULL DEFAULT '',
    resume_url    TEXT        NOT NULL DEFAULT '',
    about         TEXT        NOT NULL,
    status        TEXT        NOT NULL DEFAULT 'new'
                  CHECK (status IN ('new', 'reviewing', 'interview', 'accepted', 'rejected')),
    note          TEXT        NOT NULL DEFAULT '',
    ip            TEXT        NOT NULL DEFAULT '',
    user_agent    TEXT        NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX job_applications_status_created_idx ON job_applications (status, created_at DESC);
CREATE TRIGGER job_applications_updated_at BEFORE UPDATE ON job_applications FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Muqova: kartada qaysi nisbatda ko'rsatilishi va kesishda qaysi qismi saqlanishi
ALTER TABLE news
    ADD COLUMN cover_ratio TEXT NOT NULL DEFAULT 'auto' CHECK (cover_ratio IN ('auto', '16:9', '4:3', '1:1', '3:4', '21:9')),
    ADD COLUMN cover_focus TEXT NOT NULL DEFAULT 'center' CHECK (cover_focus IN ('center', 'top', 'bottom'));
ALTER TABLE articles
    ADD COLUMN cover_ratio TEXT NOT NULL DEFAULT 'auto' CHECK (cover_ratio IN ('auto', '16:9', '4:3', '1:1', '3:4', '21:9')),
    ADD COLUMN cover_focus TEXT NOT NULL DEFAULT 'center' CHECK (cover_focus IN ('center', 'top', 'bottom'));
