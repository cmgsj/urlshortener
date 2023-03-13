CREATE TABLE
	IF NOT EXISTS urls (
		url_id TEXT PRIMARY KEY,
		redirect_url TEXT NOT NULL UNIQUE
	);

CREATE INDEX IF NOT EXISTS redirect_url_index ON urls (redirect_url);