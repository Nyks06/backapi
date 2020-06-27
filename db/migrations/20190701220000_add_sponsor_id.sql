
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE users ADD COLUMN sponsor_id varchar(255) DEFAULT '';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE users DROP COLUMN sponsor_id;