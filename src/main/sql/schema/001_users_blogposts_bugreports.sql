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

create table blog_posts
(
    id             uuid primary key,
    created_at     timestamp not null,
    updated_at     timestamp not null,
    title          text      not null,
    description    text      not null,
    image_filename text      not null,
    user_id        uuid      not null references users (id) on delete cascade
);

create table bug_reports
(
    id             uuid primary key,
    created_at     timestamp not null,
    updated_at     timestamp not null,
    title          text      not null,
    description    text      not null,
    image_filename text      not null,
    user_id        uuid      not null references users (id) on delete cascade
);

create table files
(
    id                  uuid primary key,
    file_name           text not null,
    folder_name         text not null
);

-- +goose Down
drop table blog_posts;
drop table bug_reports;
drop table users;
drop table files;
