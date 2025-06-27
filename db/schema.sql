CREATE TYPE audit_event AS ENUM ('login', 'logout', 'failed_login');

CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(100) PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    registered_ip VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS bio (
    username VARCHAR(100) PRIMARY KEY,
    full_name VARCHAR(100),
    birthdate DATE,
    bio_text TEXT,
    FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS audit_log (
    username VARCHAR(100) PRIMARY KEY,
    event_type audit_event NOT NULL,
    ip VARCHAR(100) NOT NULL,
    ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- soon: submission link table