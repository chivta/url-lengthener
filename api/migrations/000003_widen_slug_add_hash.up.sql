ALTER TABLE urls ALTER COLUMN slug TYPE VARCHAR(8192);

ALTER TABLE urls ADD COLUMN slug_hash BYTEA;
UPDATE urls SET slug_hash = sha256(convert_to(slug, 'UTF8'));
ALTER TABLE urls ALTER COLUMN slug_hash SET NOT NULL;

DROP INDEX idx_urls_slug;
ALTER TABLE urls DROP CONSTRAINT urls_slug_key;
ALTER TABLE urls ADD CONSTRAINT urls_slug_hash_key UNIQUE (slug_hash);
