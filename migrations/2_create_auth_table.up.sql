CREATE TABLE IF NOT EXISTS auth_users (
    id INTEGER PRIMARY KEY autoincrement,
    created_at TEXT,
    email TEXT,
    password TEXT,
    is_admin TINYINT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS email_index ON auth_users(email);
