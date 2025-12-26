-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
ALTER TABLE board ALTER COLUMN id TYPE UUID USING id::UUID;
ALTER TABLE board ALTER COLUMN id SET DEFAULT uuid_generate_v4();
ALTER TABLE board ALTER COLUMN id SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE board ALTER COLUMN id DROP DEFAULT;
ALTER TABLE board ALTER COLUMN id TYPE VARCHAR(255) USING id::VARCHAR(255);
ALTER TABLE board ALTER COLUMN id SET NOT NULL;
-- +goose StatementEnd
