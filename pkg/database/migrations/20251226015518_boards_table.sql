-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS "board" (
	id VARCHAR(255) PRIMARY KEY NOT NULL,
	name VARCHAR(255) NOT NULL,
    owner_id VARCHAR(255) NOT NULL,
    elements JSONB,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT board_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES "user"(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE "board";
-- +goose StatementEnd
