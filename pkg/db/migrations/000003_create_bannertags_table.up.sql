CREATE TABLE IF NOT EXISTS bannertags(
    id serial PRIMARY KEY,
    banner_id INT NOT NULL REFERENCES banners ON DELETE CASCADE,
    tag_id INT NOT NULL,
    UNIQUE(banner_id, tag_id)
);