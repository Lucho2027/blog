-- +goose Up

CREATE TABLE
	refresh_tokens (
		token text PRIMARY KEY NOT Null,
		created_at TIMESTAMP DEFAULT now() not null,
		updated_at TIMESTAMP DEFAULT now() not null,
		user_id uuid not null references users(id) on delete cascade,
		expires_at TIMESTAMP,
		revoked_at TIMESTAMP
	);
create index idx_refresh_tokens_user_id on refresh_tokens(user_id);

-- +goose Down
DROP TABLE refresh_tokens;
