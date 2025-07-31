CREATE TABLE IF NOT EXISTS statuses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(16) NOT NULL UNIQUE
);

INSERT INTO statuses (name)
VALUES ('SUCCESS'), ('REJECT'), ('PENDING'), ('OLD');

CREATE TABLE IF NOT EXISTS teams (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(256),
    ip VARCHAR(16) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS exploits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(256) NOT NULL,
    type VARCHAR(32),
    is_local BOOLEAN NOT NULL,
    executable_path VARCHAR(256),
    requirements_path VARCHAR(256)
);

CREATE TABLE IF NOT EXISTS flags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    value VARCHAR(256) NOT NULL,
    status_id INTEGER NOT NULL,
    exploit_id INTEGER,
    get_from INTEGER,
    message_from_server TEXT,

    FOREIGN KEY(get_from) REFERENCES teams(id),
    FOREIGN KEY(status_id) REFERENCES statuses(id),
    FOREIGN KEY(exploit_id) REFERENCES exploits(id)
);
CREATE INDEX IF NOT EXISTS idx_flags_status_id ON flags(status_id);
CREATE INDEX IF NOT EXISTS idx_flags_exploit_id ON flags(exploit_id);
CREATE INDEX IF NOT EXISTS idx_flags_get_from ON flags(get_from);

CREATE TABLE IF NOT EXISTS config (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key VARCHAR(256) NOT NULL UNIQUE,
    value TEXT NOT NULL
);