CREATE TABLE IF NOT EXISTS user_game_histories (
    id            BIGSERIAL    PRIMARY KEY,
    user_id       BIGINT       NOT NULL REFERENCES users(id),
    points_earned INTEGER      NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_game_histories_user_id    ON user_game_histories (user_id);
CREATE INDEX IF NOT EXISTS idx_user_game_histories_created_at ON user_game_histories (created_at);
