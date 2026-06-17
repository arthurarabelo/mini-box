CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    username      TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    name          TEXT NOT NULL,
    email         TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- mesmos usuários que existiam hardcoded em auth/handler.go
-- senha "admin" e "alice123" já hasheadas com bcrypt
INSERT INTO users (username, password_hash, name, email) VALUES
    ('admin', '$2a$10$iBb/wpora1Pp1Uwlu3hN.eS4MHo9SiAxDZ5L12rk0YVf5Mey9aSdm', 'Admin', 'admin@cern.ch'),
    ('alice', '$2a$10$3iye.gSZZWMaPS3/s6lMT.wildymnGuY1OsNugbUky0FlDg2xCEPu', 'Alice Smith', 'alice@cern.ch')
ON CONFLICT (username) DO NOTHING;
