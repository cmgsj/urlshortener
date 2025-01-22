CREATE TABLE IF NOT EXISTS urls (
    url_id TEXT PRIMARY KEY,
    redirect_url TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS urls_redirect_url ON urls (redirect_url);
