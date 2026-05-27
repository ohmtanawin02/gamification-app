CREATE TABLE IF NOT EXISTS user_rewards (
    id          BIGSERIAL    PRIMARY KEY,
    user_id     BIGINT       NOT NULL REFERENCES users(id),
    reward_id   BIGINT       NOT NULL REFERENCES rewards(id),
    claimed_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_user_rewards UNIQUE (user_id, reward_id)
);

CREATE INDEX IF NOT EXISTS idx_user_rewards_user_id   ON user_rewards (user_id);
CREATE INDEX IF NOT EXISTS idx_user_rewards_reward_id ON user_rewards (reward_id);
