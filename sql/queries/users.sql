-- name: CreateUser :one
insert into users(email, hash_pw)
values (@email, @hash_pw)
returning id;
