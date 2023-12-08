-- name: CreateBugReport :one
INSERT INTO bug_reports (id, created_at, updated_at, title, description, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetBugReports :many
SELECT * FROM bug_reports;

-- name: GetBugReportsByUser :many
select * from bug_reports where user_id = $1;

-- name: CreateBugReportFile :one
insert into bug_report_files (bug_report_id, file_id)
values ($1, $2)
returning *;