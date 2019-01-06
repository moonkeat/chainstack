-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE users (
  id SERIAL,
  email TEXT NOT NULL,
  password TEXT NOT NULL,
  admin BOOLEAN NOT NULL DEFAULT FALSE,
  PRIMARY KEY (id)
);

CREATE UNIQUE INDEX users_unique_lower_email_idx ON users(lower(email));

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE users;
