-- Note: for proper functionality of UUIDs, you may need to install the extension
-- by running the following in the psql console:
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE audit_event AS ENUM (
    'login',
    'logout',
    'failed_login',
    'post',
    'comment',
    'post_click',
    'sent_email'
);
-- no plans to use passwords
-- instead i'm going to email magic links
-- saves the hastle of hashing passwords, plus security improvements
CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(100) PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    registered_ip VARCHAR(100) NOT NULL
);
CREATE TABLE IF NOT EXISTS magic_links (
    username VARCHAR(100) PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- known as "UserMetadata" when represented as a Go struct
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
    metadata VARCHAR(255),
    -- other information, for instance, what post was clicked, what username the failed login used
    ip VARCHAR(100) NOT NULL,
    ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- soon: submission link table
CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    link VARCHAR(255) NOT NULL,
    body TEXT,
    -- optional body text (when you visit a submission page on HN, sometimes there will be additonal text)
    flagged BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (username) REFERENCES users(username) -- notice the lack of cascade
);

-- handy function that allows you to query count of posts between timestamps
-- Usage (for posts in the last 48 hours):
-- SELECT days_between((NOW() - INTERVAL '2 days')::TIMESTAMP, NOW()::TIMESTAMP);
CREATE OR REPLACE FUNCTION days_between(ts1 TIMESTAMP, ts2 TIMESTAMP)
RETURNS INTEGER AS $$
DECLARE
  result INTEGER;
BEGIN
  SELECT COUNT(*) INTO result
  FROM submissions
  WHERE created_at BETWEEN ts1 AND ts2;

  RETURN result;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE IF NOT EXISTS votes (
    id SERIAL PRIMARY KEY,
    submission_id UUID NOT NULL,
    voter_username VARCHAR(100) NOT NULL,
    ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    positive BOOLEAN NOT NULL,
    -- true = upvote, false = downvote
    FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE,
    FOREIGN KEY (voter_username) REFERENCES users(username) ON DELETE CASCADE,
    UNIQUE(submission_id, voter_username)
);
-- note: use a recursive query to build a comment chain (via self join)
-- FK summary:
--  author --> users(username)
--  parent_comment --> submissions(id)
--  in_reponse_to --> comments(id) - OPTIONAL
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    in_response_to UUID NOT NULL,
    content TEXT NOT NULL,
    author VARCHAR(100) NOT NULL,
    parent_comment UUID NULL,
    flagged BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_author FOREIGN KEY (author) REFERENCES users(username),
    CONSTRAINT fk_parent_comment FOREIGN KEY (parent_comment) REFERENCES comments(id),
    CONSTRAINT fk_in_response_to FOREIGN KEY (in_response_to) REFERENCES submissions(id)
);

-- essentially a clone of the normal votes table, yet this time it's for comments
CREATE TABLE IF NOT EXISTS comment_votes (
    id SERIAL PRIMARY KEY,
    comment_id UUID NOT NULL,
    voter_username VARCHAR(100) NOT NULL,
    ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    positive BOOLEAN NOT NULL,
    -- true = upvote, false = downvote
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (voter_username) REFERENCES users(username) ON DELETE CASCADE,
    UNIQUE(comment_id, voter_username)
);


CREATE TABLE IF NOT EXISTS admins (
    username VARCHAR(100) PRIMARY KEY,
    remarks TEXT,
    FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE
);

-- users for automated access (API)
CREATE TABLE api_tokens (
    username VARCHAR(100) PRIMARY KEY,
    token VARCHAR(255) NOT NULL,
    FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE
);

-- table of reports: both comments and posts
-- REPORT WEIGHT CHART
-- |------------------|---------------|
-- | Account Age      | Report Weight |
-- |------------------|---------------|
-- | < 1 day          |   0.1         |
-- | 1 day - 7 days   |   0.25        |
-- | 7 days - 28 days |   0.33        |
-- | 28+ days         |   0.5         |
-- |------------------|---------------|

CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    reporter VARCHAR(100) NOT NULL REFERENCES users(username),
    target_type VARCHAR(20) NOT NULL,  -- 'post', 'comment'
    target_id UUID NOT NULL,            -- references post.id OR comment.id (i don't know if there is a way to check this)
    target_user VARCHAR(100) NOT NULL REFERENCES users(username),
    rweight FLOAT NOT NULL, -- "weight" of the report (logic determined on frontend)
    created_at TIMESTAMP DEFAULT NOW()
);