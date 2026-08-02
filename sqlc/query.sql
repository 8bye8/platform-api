-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    password_hash,
    password_salt
) VALUES (
  @username,
  @email,
  @password_hash,
  @password_salt
) RETURNING *;

-- name: GetUserByUsername :one
SELECT id, password_hash, password_salt FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateChallenge :one
INSERT INTO challenges (
                        created_by,
                        side_to_play,
                        time_control,
                        max_rating,
                        min_rating,
                        rating
) VALUES (
           @created_by,
           @side_to_play,
           @time_control,
           @max_rating,
           @min_rating,
           @rating
) RETURNING *;


-- name: GetChallengeToMatch :one
SELECT * FROM challenges
WHERE
    side_to_play = @side_to_play AND
    status = 'open' AND
    time_control = @time_control AND
    rating <= @max_rating AND
    rating >= @min_rating AND
    @rating <= max_rating AND
    @rating >= min_rating
ORDER BY created_at
    FOR NO KEY UPDATE SKIP LOCKED
LIMIT 1;

-- WITH locked_challenge AS (
--
-- )
-- UPDATE challenges
-- SET status = $5 AND matched_with = $6
-- FROM locked_challenge
-- WHERE challenges.id = locked_challenge.id


