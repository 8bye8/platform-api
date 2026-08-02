-- +goose Up
-- +goose StatementBegin

CREATE TABLE user_stats (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
   bullet_rating INTEGER NOT NULL,
   blitz_rating INTEGER NOT NULL,
   rapid_rating INTEGER NOT NULL,
   classical_rating INTEGER NOT NULL,
   bullet_games_won INTEGER NOT NULL DEFAULT 0,
   bullet_games_lost INTEGER NOT NULL DEFAULT 0,
   bullet_games_drawn INTEGER NOT NULL DEFAULT 0,
   blitz_games_won INTEGER NOT NULL DEFAULT 0,
   blitz_games_lost INTEGER NOT NULL DEFAULT 0,
   blitz_games_drawn INTEGER NOT NULL DEFAULT 0,
   rapid_games_won INTEGER NOT NULL DEFAULT 0,
   rapid_games_lost INTEGER NOT NULL DEFAULT 0,
   rapid_games_drawn INTEGER NOT NULL DEFAULT 0,
   classical_games_won INTEGER NOT NULL DEFAULT 0,
   classical_games_lost INTEGER NOT NULL DEFAULT 0,
   classical_games_drawn INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX users_stats_user_id_idx
    ON user_stats USING hash (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS user_stats;

-- +goose StatementEnd
