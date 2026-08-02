-- +goose Up
-- +goose StatementBegin

CREATE TYPE challenge_side_to_play AS ENUM ('white', 'black', 'random');
CREATE TYPE challenge_status AS ENUM ('open', 'matched', 'suspended');
CREATE TYPE challenge_time_control AS ENUM (
   '1_0',
   '2_1',
   '3_0',
   '3_2',
   '5_0',
   '10_0',
   '15_0',
   '15_10',
   '30_0',
   '60_60'
);

CREATE TABLE challenges (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   matched_with UUID REFERENCES challenges(id),
   created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
   side_to_play challenge_side_to_play NOT NULL,
   status challenge_status NOT NULL DEFAULT 'open',
   time_control challenge_time_control NOT NULL,
   max_rating INTEGER NOT NULL,
   min_rating INTEGER NOT NULL,
   rating INTEGER NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   matched_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX challenges_created_at_idx
    ON challenges (created_at);

CREATE INDEX challenges_matching_idx
    ON challenges (side_to_play, status, time_control, rating);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS challenges;
DROP TYPE IF EXISTS challenge_time_control;
DROP TYPE IF EXISTS challenge_status;
DROP TYPE IF EXISTS challenge_side_to_play;

-- +goose StatementEnd
