-- ONLY RUN THIS AFTER SUBMISSIONS AND USERS ARE INSERTED
-- (inserts 30 gibberish comments, from random users on random posts)
DO $$ BEGIN FOR i IN 1..30 LOOP
INSERT INTO comment_votes (
        comment_id,
        voter_username,
        positive
    )
VALUES (
        (
            SELECT id
            FROM comments
            ORDER BY random()
            LIMIT 1
        ),
        (
            SELECT username
            FROM users
            ORDER BY random()
            LIMIT 1
        ),
        (
            SELECT val
            FROM (VALUES (TRUE), (FALSE)) AS t(val)
            ORDER BY random()
            LIMIT 1
        )
    );
END LOOP;
END $$;