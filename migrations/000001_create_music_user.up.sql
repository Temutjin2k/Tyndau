CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email text UNIQUE NOT NULL,
    avatar_link text,
    password_hash bytea NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    version integer NOT NULL DEFAULT 1
);


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