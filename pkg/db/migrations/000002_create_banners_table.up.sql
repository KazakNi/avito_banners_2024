CREATE TABLE IF NOT EXISTS banners(
    id serial PRIMARY KEY,
    feature_id INT NOT NULL,
    title VARCHAR(150) NOT NULL,
    text VARCHAR(300) NOT NULL,
    url VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
