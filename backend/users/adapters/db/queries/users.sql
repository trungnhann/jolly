-- name: CreateUser :exec
INSERT INTO users.users (
	user_uuid,
	email,
	name,
	role,
	created_at,
	updated_at
)
VALUES
	($1, $2, $3, $4, $5, $6)
;

-- name: GetUser :one
SELECT
	*
FROM
	users.users
WHERE
	user_uuid = $1
LIMIT 1;
