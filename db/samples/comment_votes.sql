-- ONLY RUN THIS AFTER COMMENTS AND USERS ARE INSERTED
-- (inserts 30 random votes, from random users on random comments)
-- note that this could fail due to multiple votes from a single user being attempted on a single comment; i'll find a way to prevent this in the future
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