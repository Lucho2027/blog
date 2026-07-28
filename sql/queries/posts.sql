-- name: CreateDraftPost :one
insert into posts ( title, excerpt, content_md, status, user_id, published_at )
values ( @title, @excerpt, @content_md, 'draft', @user_id, null)
returning *;

-- name: GetPostById :one
select * from posts where id = @id and status = 'published';

-- name: GetAllPosts :many
select * from posts where status = 'published'
order by published_at desc
LIMIT @limit_val OFFSET @offset_val;
