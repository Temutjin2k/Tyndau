CREATE TABLE IF NOT EXISTS songs (
    id           BIGSERIAL PRIMARY KEY,
    title        TEXT NOT NULL,
    artist       TEXT NOT NULL,
    album        TEXT NOT NULL,
    genre        TEXT NOT NULL,
    duration_sec INT NOT NULL,
    file_url     TEXT NOT NULL,
    released_at  TIMESTAMP,
    created_at   TIMESTAMP DEFAULT now(),
    updated_at   TIMESTAMP DEFAULT now()
);