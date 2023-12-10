-- name: CreateUser :one
INSERT INTO users (
  username,
  password,
  email
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: SearchUsers :many
SELECT profile_picture, username FROM users
WHERE user_id <> $1
ORDER BY RANDOM()
LIMIT 10;

-- name: UpdateUser :one
UPDATE users 
SET profile_picture = $1, biography = $2
WHERE user_id = $3
RETURNING 1;