ALTER TABLE urls ADD COLUMN expires_at TIMESTAMPTZ;
UPDATE urls
SET expires_at = TO_TIMESTAMP(0);
CREATE INDEX idx_urls_expires_at
ON urls(expires_at);
