-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL
);

CREATE TABLE sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT NOT NULL UNIQUE,
    user_id INTEGER,
    created_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL,

    FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE monitors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT NOT NULL UNIQUE,
    `url` TEXT NOT NULL,
    created_at DATETIME NOT NULL
);

CREATE TABLE checks(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT NOT NULL UNIQUE,
    monitor_id INTEGER,
    created_at DATETIME NOT NULL,
    status_code INTEGER NOT NULL,
    response_time_ms INTEGER NOT NULL,
    
    FOREIGN KEY(monitor_id) REFERENCES monitors(id)
);

CREATE TABLE incidents(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT NOT NULL UNIQUE,
    monitor_id INTEGER,
    status_text TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    response_time_ms INTEGER NOT NULL,
    body TEXT,
    headers TEXT,
    created_at DATETIME NOT NULL,

    FOREIGN KEY(monitor_id) REFERENCES monitors(id)
);
-- +goose StatementEnd