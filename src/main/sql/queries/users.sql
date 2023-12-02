-- name: CreateUser :one
insert into users (id, created_at, updated_at, email, password, username, token)
values ($1, $2, $3, $4, $5, $6,
        '')
returning *;

-- name: GetUserByEmail :one
select * from users where email = $1;

-- name: UpdateUserToken :exec
update users set token = $2 where id = $1;

-- name: GetUserByToken :one
SELECT * FROM users WHERE token = $1;
