CREATE TABLE clicks (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    url_id     UUID        NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
    ip_hash    VARCHAR(64) NOT NULL DEFAULT '',
    user_agent TEXT        NOT NULL DEFAULT '',
    country    CHAR(2)     NOT NULL DEFAULT '',
    clicked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_clicks_url_id ON clicks(url_id, clicked_at DESC);
