-- name: CreateUser :exec
INSERT INTO users.users (
	user_uuid,
	email,
	name,
	role,
	password_hash,
	avatar_url,
	created_at,
	updated_at
)
VALUES
	($1, $2, $3, $4, $5, $6, $7, $8)
;

-- name: GetUser :one
SELECT
	*
FROM
	users.users
WHERE
	user_uuid = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT
	*
FROM
	users.users
WHERE
	email = $1
LIMIT 1;

-- name: UpdateUserAvatar :exec
UPDATE users.users
SET
	avatar_url = $2,
	updated_at = $3
WHERE
	user_uuid = $1;
