
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

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
  sponsorship_id                       varchar(255)      NOT NULL,
  firstname                            varchar(255)      NOT NULL,
  lastname                             varchar(255)      NOT NULL,
  username                             varchar(255)      NOT NULL,
  email                                varchar(255)      NOT NULL,
  phone_number                         varchar(255)      NOT NULL,
  telegram                             varchar(255)      NOT NULL,
  password                             varchar(255)      NOT NULL,
  confirmed                            boolean           NOT NULL DEFAULT FALSE,
  created_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc(),
  updated_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc()
);

CREATE TABLE IF NOT EXISTS sessions (
  id                                   uuid              PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id                              varchar(255)      NOT NULL,
  created_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc(),
  updated_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc(),
  expires_at                           timestamp         WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS account_recover (
  id                                   uuid              PRIMARY KEY DEFAULT uuid_generate_v4(),
  token                                text              NOT NULL,
  created_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc(),
  used_at                              timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc(),
  active                               boolean           DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS account_confirmation (
  id                                   uuid              PRIMARY KEY DEFAULT uuid_generate_v4(),
  token                                text              NOT NULL,
  created_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc(),
  used_at                              timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc(),
  active                               boolean           DEFAULT TRUE
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE users;
DROP TABLE sessions;
DROP TABLE account_recover;
DROP TABLE account_confirmation;

DROP FUNCTION now_utc();
