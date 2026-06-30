BEGIN;
ALTER TABLE users.users ADD COLUMN avatar_url varchar(1024) NULL;
COMMIT;
