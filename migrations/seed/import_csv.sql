BEGIN;

CREATE TABLE csv_staging (
    nickname VARCHAR,
    point    INT,
    datetime TIMESTAMPTZ
);

COPY csv_staging
FROM '/tmp/mock_data.csv'
WITH (FORMAT csv, HEADER true);


WITH agg AS (
    SELECT
        nickname,
        LEAST(SUM(point), 10000) AS total_points
    FROM csv_staging
    GROUP BY nickname
),
updated AS (
    UPDATE users u
    SET total_points = a.total_points,
        updated_at   = NOW()
    FROM agg a
    WHERE u.nickname = a.nickname
    RETURNING u.nickname
)
INSERT INTO users (nickname, total_points, created_at, updated_at)
SELECT a.nickname, a.total_points, NOW(), NOW()
FROM agg a
WHERE a.nickname NOT IN (SELECT nickname FROM updated);


INSERT INTO user_game_histories (user_id, points_earned, created_at)
SELECT
    u.id,
    s.point,
    s.datetime
FROM csv_staging s
JOIN users u ON u.nickname = s.nickname;


DROP TABLE csv_staging;

COMMIT;