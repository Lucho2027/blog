.PHONY: start sqlc types gen

start:
	set -a; . ./.env; set +a; go run ./cmd/api

sqlc:
	sqlc generate

types:
	go run ./tools/extract_ts_types -in web/types/generated/posts_sql.ts -out web/types/generated/posts.types.ts

gen: sqlc types

