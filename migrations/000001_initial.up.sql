-- migrations/000001_initial.up.sql
-- Создание таблицы ссылок
CREATE TABLE short_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    original_url VARCHAR(1000) NOT NULL UNIQUE,
    short_url VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);