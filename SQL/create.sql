CREATE TABLE IF NOT EXISTS links (
	short_link VARCHAR(10) PRIMARY KEY,
	original_link VARCHAR(260) NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS index_links_original_link ON links (original_link);
