-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION now_utc() returns timestamp as $$
  select now() at time zone 'utc';
$$ language sql;

CREATE OR REPLACE FUNCTION set_updated_at_column() RETURNS TRIGGER AS $$
  BEGIN
   NEW.updated_at = now_utc();
   RETURN NEW;
  END;
$$ language 'plpgsql';
-- +goose StatementEnd


CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE IF NOT EXISTS users (
  id                                   uuid              PRIMARY KEY DEFAULT uuid_generate_v4(),
  customer_id                          varchar(255)      NOT NULL,
  email                                varchar(255)      NOT NULL,
  password                             varchar(255)      NOT NULL,
  firstname                            varchar(255)      NOT NULL,
  lastname                             varchar(255)      NOT NULL,
  username                             varchar(255)      NOT NULL,
  phone_number                         varchar(255)      NOT NULL,
  admin                                boolean           NOT NULL DEFAULT FALSE,
  created_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc(),
  updated_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc(),
  confirmed_at                         timestamp         WITH TIME ZONE NOT NULL,
  deleted_at                           timestamp         WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
  id                                   uuid              PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id                              varchar(255)      NOT NULL,
  created_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc(),
  updated_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc(),
  expires_at                           timestamp         WITH TIME ZONE NOT NULL
);


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS sessions;

DROP FUNCTION now_utc();
