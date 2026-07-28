-- +goose Up

create table
  posts (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title text not null,
    excerpt text default null,
    content_md text not null,
    status text not null,
    user_id uuid not null references users(id),
    published_at timestamp,
    created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
  );

-- +goose Down
  drop table posts;
