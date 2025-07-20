-- ONLY RUN THIS AFTER SUBMISSIONS AND USERS ARE INSERTED
-- (inserts 30 gibberish comments, from random users on random posts)
DO $$ BEGIN FOR i IN 1..30 LOOP
INSERT INTO comments (
        in_response_to,
        content,
        author,
        flagged
    )
VALUES (
        -- in response to a random post
        (
            SELECT id
            FROM submissions
            ORDER BY random()
            LIMIT 1
        ), -- generated gibberish from Postgres
        (
            SELECT string_agg(chr(97 + floor(random() * 26)::int), '')
            FROM generate_series(1, 100)
        ),
        -- author is a random user
        (
            SELECT username
            FROM users
            ORDER BY random()
            LIMIT 1
        ), FALSE
    );
END LOOP;
END $$;