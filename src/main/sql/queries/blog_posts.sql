-- name: CreateBlogPost :one
insert into blog_posts (id, created_at, updated_at, title, description, user_id)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetBlogPosts :many
select
    blog_posts.id as blog_post_id,
    blog_posts.created_at as blog_post_created_at,
    blog_posts.updated_at as blog_post_updated_at,
    blog_posts.title as blog_post_title,
    blog_posts.description as blog_post_description,
    files.id as file_id,
    files.file_name,
    files.folder_name
from
    blog_posts
        join
    blog_post_files on blog_posts.id = blog_post_files.blog_post_id
        join
    files on blog_post_files.file_id = files.id;

-- name: GetBlogPostsByUserWithFiles :many
select
    blog_posts.id as blog_post_id,
    blog_posts.created_at as blog_post_created_at,
    blog_posts.updated_at as blog_post_updated_at,
    blog_posts.title as blog_post_title,
    blog_posts.description as blog_post_description,
    files.id as file_id,
    files.file_name,
    files.folder_name
from
    blog_posts
        join
    blog_post_files on blog_posts.id = blog_post_files.blog_post_id
        join
    files on blog_post_files.file_id = files.id
where
        blog_posts.user_id = $1;

-- name: GetBlogPostsByCreatedAt :many
with WeeklyAggregation as (
    select
        extract(week from created_at) as week_number,
        extract(year from created_at) as year,
        count(*) as post_count
    from
        blog_posts
    group by
        week_number, year
)

select
    bp.*,
    wa.week_number,
    wa.year,
    wa.post_count
from
    blog_posts bp
        join
    WeeklyAggregation wa on extract(week from bp.created_at) = wa.week_number and extract(year from bp.created_at) = wa.year
order by
    wa.week_number;

-- name: CreateBlogPostFile :one
insert into blog_post_files (blog_post_id, file_id)
values ($1, $2)
returning *;

-- name: GetFilesByBlogPostID :many
select files.id, files.file_name, files.folder_name
from blog_post_files
join files on blog_post_files.file_id = files.id
where blog_post_files.blog_post_id = $1;
