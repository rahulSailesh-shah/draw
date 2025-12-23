-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
	id VARCHAR(255) PRIMARY KEY NOT NULL,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	email_verified BOOLEAN DEFAULT false NOT NULL,
	image VARCHAR(255),
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT user_email_unique UNIQUE("email")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE "user"
-- +goose StatementEnd
