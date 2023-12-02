-- +goose Up

create table users
(
    id         uuid primary key,
    created_at timestamp          not null,
    updated_at timestamp          not null,
    email      text               not null,
    password   text               not null,
    username   text               not null,
    api_key    varchar(64) unique not null default (
        encode(sha256(random()::text::bytea), 'hex')
        )
);

CREATE table blog_posts
(
    id             uuid primary key,
    created_at     timestamp not null,
    updated_at     timestamp not null,
    title          text      not null,
    description    text,
    image_filename text,
    user_id        uuid      not null references users (id) on delete cascade
);

CREATE table bug_reports
(
    id             uuid primary key,
    created_at     timestamp not null,
    updated_at     timestamp not null,
    title          text      not null,
    description    text,
    image_filename text,
    user_id        uuid      not null references users (id) on delete cascade
);

-- +goose Down
drop table blog_posts;
drop table bug_reports;
drop table users;
