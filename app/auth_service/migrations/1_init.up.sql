CREATE TABLE IF NOT EXISTS Users
(
    id            INTEGER PRIMARY KEY,
    first_name    TEXT NOT NULL,
    last_name     TEXT NOT NULL,
    login         TEXT NOT NULL UNIQUE,
    avatar        TEXT NOT NULL DEFAULT 'default.jpg',
    password_hash BLOB NOT NULL
);