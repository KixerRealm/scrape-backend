-- +goose Up

create table users
(
    id         uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    email      text      not null,
    password   text      not null,
    username   text      not null,
    token      text unique not null
);

create table files
(
    id                  uuid primary key,
    file_name           text not null,
    folder_name         text not null
);

create table blog_posts
(
    id             uuid primary key,
    created_at     timestamp not null,
    updated_at     timestamp not null,
    title          text      not null,
    description    text      not null,
    user_id        uuid      not null references users (id) on delete cascade
);

create table bug_reports
(
    id             uuid primary key,
    created_at     timestamp not null,
    updated_at     timestamp not null,
    title          text      not null,
    description    text      not null,
    user_id        uuid      not null references users (id) on delete cascade
);

create table blog_post_files
(
    blog_post_id uuid not null references blog_posts (id) on delete cascade,
    file_id      uuid not null references files (id) on delete cascade,
    primary key (blog_post_id, file_id)
);

create table bug_report_files
(
    bug_report_id uuid not null references bug_reports (id) on delete cascade,
    file_id       uuid not null references files (id) on delete cascade,
    primary key (bug_report_id, file_id)
);

-- +goose Down
drop table blog_post_files;
drop table blog_posts;
drop table bug_report_files;
drop table bug_reports;
drop table files;
drop table users;
