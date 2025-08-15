CREATE TABLE IF NOT EXISTS statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(16) NOT NULL UNIQUE
);

INSERT INTO statuses (name) VALUES
  ('SUCCESS'), ('REJECT'), ('PENDING'), ('OLD')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256),
    ip INET NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS exploits (
    id UUID PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    type VARCHAR(32),
    is_running_on_farm BOOLEAN NOT NULL,
    executable_path VARCHAR(256),
    requirements_path VARCHAR(256),
    is_running BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS flags (
    id SERIAL PRIMARY KEY,
    value VARCHAR(256) NOT NULL UNIQUE,
    status_id INT NOT NULL REFERENCES statuses(id),
    exploit_id UUID REFERENCES exploits(id),
    get_from INT REFERENCES teams(id),
    message_from_server TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_flags_status_id ON flags(status_id);
CREATE INDEX IF NOT EXISTS idx_flags_exploit_id ON flags(exploit_id);
CREATE INDEX IF NOT EXISTS idx_flags_get_from ON flags(get_from);

CREATE TABLE IF NOT EXISTS config (
    id SERIAL PRIMARY KEY,
    key VARCHAR(256) NOT NULL UNIQUE,
    value TEXT NOT NULL
);
