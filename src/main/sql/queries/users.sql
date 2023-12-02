-- name: CreateUser :one
insert into users (id, created_at, updated_at, email, password, username, api_key)
values ($1, $2, $3, $4, $5, $6,
        encode(sha256(random()::text::bytea), 'hex'))
returning *;
