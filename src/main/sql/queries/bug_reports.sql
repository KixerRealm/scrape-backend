-- name: CreateBugReport :one
INSERT INTO bug_reports (id, created_at, updated_at, title, description, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetBugReports :many
select
    bug_reports.id as bug_report_id,
    bug_reports.created_at as bug_report_created_at,
    bug_reports.updated_at as bug_report_updated_at,
    bug_reports.title as bug_report_title,
    bug_reports.description as bug_report_description,
    files.id as file_id,
    files.file_name,
    files.folder_name
from
    bug_reports
        join
    bug_report_files on bug_reports.id = bug_report_files.bug_report_id
        join
    files on bug_report_files.file_id = files.id;

-- name: GetBugReportsByUserWithFiles :many
select
    bug_reports.id as bug_report_id,
    bug_reports.created_at as bug_report_created_at,
    bug_reports.updated_at as bug_report_updated_at,
    bug_reports.title as bug_report_title,
    bug_reports.description as bug_report_description,
    files.id as file_id,
    files.file_name,
    files.folder_name
from
    bug_reports
        join
    bug_report_files on bug_reports.id = bug_report_files.bug_report_id
        join
    files on bug_report_files.file_id = files.id
where
        bug_reports.user_id = $1;

-- name: CreateBugReportFile :one
insert into bug_report_files (bug_report_id, file_id)
values ($1, $2)
returning *;