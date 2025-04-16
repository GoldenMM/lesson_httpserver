-- +goose Up
ALTER TABLE users
ADD COLUMN is_chirpy_red BOOLEAN NOT NULL DEFAULT FALSE;
UPDATE users
SET is_chirpy_red = FALSE;


-- +goose Down
ALTER TABLE users
DROP COLUMN is_chirpy_red;
