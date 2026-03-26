CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    hostname VARCHAR(255) NOT NULL,
    ip VARCHAR(45) NOT NULL,
    location VARCHAR(255),
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_devices_hostname_trgm ON devices USING gin (hostname gin_trgm_ops);

CREATE INDEX idx_devices_active_deleted_id ON devices(is_active, is_deleted, id)
WHERE
    is_deleted = false;

CREATE INDEX idx_devices_is_active ON devices(is_active)
WHERE
    is_deleted = false;