CREATE TABLE IF NOT EXISTS rewards (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    check_point INTEGER      NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

INSERT INTO rewards (name, check_point) VALUES
    ('รางวัล 1', 500),
    ('รางวัล 2', 1000),
    ('รางวัล 3', 10000)
ON CONFLICT DO NOTHING;