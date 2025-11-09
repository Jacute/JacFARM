CREATE SCHEMA IF NOT EXISTS audit;

CREATE TABLE IF NOT EXISTS audit.log_levels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(16) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS audit.modules (
    id SERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS audit.logs (
    id BIGSERIAL PRIMARY KEY,
    module_id INT REFERENCES audit.modules (id),
    operation VARCHAR(256),
    exploit_id UUID REFERENCES exploits (id),
    value TEXT NOT NULL,
    attrs JSONB,
    created_at TIMESTAMP NOT NULL,
    log_level_id INT NOT NULL REFERENCES audit.log_levels (id)
);
