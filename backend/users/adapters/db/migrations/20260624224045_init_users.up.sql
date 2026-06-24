BEGIN;

CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE users.users
(
	user_uuid  uuid         NOT NULL,
	email      varchar(255) NOT NULL,
	name       varchar(255) NOT NULL,
	role       varchar(32)  NOT NULL,
	created_at TIMESTAMPTZ  NOT NULL,
	updated_at TIMESTAMPTZ  NOT NULL,
	PRIMARY KEY (user_uuid),
	CONSTRAINT users_users_email_unique UNIQUE (email)
);

COMMIT;
