-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS users
(
  id                             BIGSERIAL                              NOT NULL PRIMARY KEY,
  name                           VARCHAR                                NOT NULL,
  username                       VARCHAR(160)                           NOT NULL,
  password                       VARCHAR(60)                            NOT NULL,
  created_at                     TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
  updated_at                     TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);


CREATE UNIQUE INDEX unique_users_username_index ON users(username);
