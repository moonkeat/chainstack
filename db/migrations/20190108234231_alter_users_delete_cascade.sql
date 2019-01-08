-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE resources DROP CONSTRAINT resources_user_id_fkey;
ALTER TABLE access_tokens DROP CONSTRAINT access_tokens_user_id_fkey;
ALTER TABLE resources ADD FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE access_tokens ADD FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE resources DROP CONSTRAINT resources_user_id_fkey;
ALTER TABLE access_tokens DROP CONSTRAINT access_tokens_user_id_fkey;
ALTER TABLE resources ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE access_tokens ADD FOREIGN KEY (user_id) REFERENCES users (id);
