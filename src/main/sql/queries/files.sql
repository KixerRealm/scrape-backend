-- name: CreateFile :one
insert into files (id, file_name, folder_name)
values ($1, $2, $3)
    returning *;