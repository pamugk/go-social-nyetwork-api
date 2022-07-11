CREATE TYPE gender AS ENUM ('FEMALE', 'MALE');

CREATE TABLE IF NOT EXISTS "user"(
    -- Technical information
    id BIGSERIAL PRIMARY KEY,
    created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    active BOOLEAN DEFAULT TRUE,
    login VARCHAR(100) NOT NULL,
    password CHAR(60) NOT NULL,
    preferred_locale VARCHAR(6) NOT NULL,
    -- Personal data
    name VARCHAR(150) NOT NULL,
    surname VARCHAR(200) NOT NULL,
    patronymic VARCHAR(175),
    birthday DATE NOT NULL,
    gender gender NOT NULL,
    country CHAR(3) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(400),
    about VARCHAR(300)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_login_idx ON "user"(login) WHERE active;
