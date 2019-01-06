-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE access_tokens (
  id SERIAL,
  token TEXT NOT NULL,
  expires TIMESTAMP NOT NULL,
  scope TEXT NOT NULL,
  user_id SERIAL REFERENCES users,
  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX access_token_unique_token_idx ON access_tokens(token);
CREATE INDEX access_token_expires_idx ON access_tokens(expires ASC);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE access_tokens;
