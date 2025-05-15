CREATE TABLE IF NOT EXISTS songs (
                                     id           BIGSERIAL PRIMARY KEY,
                                     title        TEXT    NOT NULL,
                                     artist       TEXT    NOT NULL,
                                     album        TEXT,
                                     genre        TEXT,
                                     duration_sec INT     NOT NULL,
                                     file_url     TEXT,
                                     released_at  TIMESTAMP,
                                     created_at   TIMESTAMP DEFAULT now(),
    updated_at   TIMESTAMP DEFAULT now()
    );
