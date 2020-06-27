
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE IF NOT EXISTS pronostics (
  id                                   uuid              PRIMARY KEY DEFAULT uuid_generate_v4(),
  pronostics_group_id                  varchar(255)      NOT NULL,
  first_team                           varchar(255)      NOT NULL,
  second_team                          varchar(255)      NOT NULL,
  pronostic                            varchar(255)      NOT NULL,
  competition_id                       varchar(255)      NOT NULL,
  sport_id                             varchar(255)      NOT NULL,
  status                               varchar(255)      NOT NULL,
  odd                                  numeric(9,6)      NOT NULL,
  event_date                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc(),
  created_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc(),
  updated_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc()
);

CREATE TABLE IF NOT EXISTS tickets (
  id                                   uuid              PRIMARY KEY DEFAULT uuid_generate_v4(),
  title                                varchar(255)      NOT NULL,
  stake                                numeric(9,6)      NOT NULL,
  public                               boolean           NOT NULL DEFAULT false,
  risk                                 varchar(255)      NOT NULL,
  live                                 boolean           NOT NULL DEFAULT false,
  created_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc(),
  updated_at                           timestamp         WITH TIME ZONE NOT NULL DEFAULT now_utc()
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE tickets;
DROP TABLE pronostics;
