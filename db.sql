CREATE TABLE IF NOT EXISTS "user" (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    password VARCHAR(128) NOT NULL
);

CREATE UNIQUE INDEX unique_user_name ON "user" (name);