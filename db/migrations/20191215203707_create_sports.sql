
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE IF NOT EXISTS sports (
  id                                   uuid              PRIMARY KEY DEFAULT uuid_generate_v4(),
  name                                 varchar(255)      NOT NULL,
  created_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc(),
  updated_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc()
);

CREATE TABLE IF NOT EXISTS competitions (
  id                                   uuid              PRIMARY KEY DEFAULT uuid_generate_v4(),
  sport_id                             varchar(255)      NOT NULL,
  name                                 varchar(255)      NOT NULL,
  start_at                             timestamp         WITH TIME ZONE NOT NULL,
  end_at                               timestamp         WITH TIME ZONE NOT NULL,
  created_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc(),
  updated_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc()
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE sports;
DROP TABLE competitions;
