CREATE TABLE clicks (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    url_id     UUID        NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
    ip_hash    VARCHAR(64) NOT NULL DEFAULT '',
    user_agent TEXT        NOT NULL DEFAULT '',
    country    CHAR(2)     NOT NULL DEFAULT '',
    clicked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_clicks_url_id ON clicks(url_id, clicked_at DESC);

ALTER TABLE urls ALTER COLUMN slug TYPE VARCHAR(8192);

ALTER TABLE urls ADD COLUMN slug_hash BYTEA;
UPDATE urls SET slug_hash = sha256(convert_to(slug, 'UTF8'));
ALTER TABLE urls ALTER COLUMN slug_hash SET NOT NULL;

DROP INDEX idx_urls_slug;
ALTER TABLE urls DROP CONSTRAINT urls_slug_key;
ALTER TABLE urls ADD CONSTRAINT urls_slug_hash_key UNIQUE (slug_hash);
