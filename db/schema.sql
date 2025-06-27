-- Note: for proper functionality of UUIDs, you may need to install the extension
-- by running the following in the psql console:
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE audit_event AS ENUM ('login', 'logout', 'failed_login', 'post', 'comment', 'post_click', 'sent_email');

-- no plans to use passwords
-- instead i'm going to email magic links
-- saves the hastle of hashing passwords, plus security improvements
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
    metadata VARCHAR(255), -- other information, for instance, what post was clicked, what username the failed login used
    ip VARCHAR(100) NOT NULL,
    ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- soon: submission link table

CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(100) PRIMARY KEY,
    link VARCHAR(255) NOT NULL,
    body TEXT, -- optional body text (when you visit a submission page on HN, sometimes there will be additonal text)
    flagged BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (username) REFERENCES users(username) -- notice the lack of cascade
);

CREATE TABLE IF NOT EXISTS votes (
    id SERIAL PRIMARY KEY,
    submission_id UUID NOT NULL,
    voter_username VARCHAR(100) NOT NULL,
    ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    positive BOOLEAN NOT NULL, -- true = upvote, false = downvote
    FOREIGN KEY (submission_id) REFERENCES submissions(id),
    FOREIGN KEY (voter_username) REFERENCES users(username)
);