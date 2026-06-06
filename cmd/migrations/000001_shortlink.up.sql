CREATE SCHEMA IF NOT EXISTS shortlink;

CREATE TABLE IF NOT EXISTS shortlink.user (
    id BIGSERIAL PRIMARY KEY,

    username VARCHAR(32) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    register_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(username),
    UNIQUE(email)
)

CREATE TABLE IF NOT EXISTS shortlink.shortlink (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL
    REFERENCES shortlink.user(id)
    ON DELETE CASCADE,

    shortlink VARCHAR(50) NOT NULL,
    original_link VARCHAR(255) NOT NULL,
    count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(user_id),
    UNIQUE(shortlink)
)

CREATE TABLE IF NOT EXISTS shortlink.shortlink_activate (
    id BIGSERIAL PRIMARY KEY,

    ip_addr VARCHAR(100) NOT NULL,
    user_agent VARCHAR(150) NOT NULL,
    shortlink_id BIGINT NOT NULL
    REFERENCES shortlink.shortlink(id)
    ON DELETE CASCADE,
    last_activate TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(ip_addr)
)