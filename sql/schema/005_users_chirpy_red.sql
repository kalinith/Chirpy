-- +goose up
ALTER TABLE users
	ADD COLUMN is_chirpy_red bool default false not null;

-- +goose down
ALTER TABLE users
	DROP COLUMN is_chirpy_red;
