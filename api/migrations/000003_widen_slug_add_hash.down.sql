ALTER TABLE urls DROP CONSTRAINT urls_slug_hash_key;
ALTER TABLE urls DROP COLUMN slug_hash;

ALTER TABLE urls ADD CONSTRAINT urls_slug_key UNIQUE (slug);
CREATE INDEX idx_urls_slug ON urls(slug);

ALTER TABLE urls ALTER COLUMN slug TYPE VARCHAR(16);
