-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE resources (
  id SERIAL,
  key TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  user_id SERIAL REFERENCES users,
  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX resources_unique_key_idx ON resources(key);
CREATE INDEX resources_user_id_idx ON resources(user_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE resources;
