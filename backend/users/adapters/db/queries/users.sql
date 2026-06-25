-- name: CreateUser :exec
INSERT INTO users.users (
	user_uuid,
	email,
	name,
	role,
	password_hash,
	created_at,
	updated_at
)
VALUES
	($1, $2, $3, $4, $5, $6, $7)
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
