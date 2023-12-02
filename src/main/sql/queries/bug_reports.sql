-- name: CreateBugReport :one
INSERT INTO bug_reports (id, created_at, updated_at, title, description, image_filename, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetBugReports :many
SELECT * FROM bug_reports;