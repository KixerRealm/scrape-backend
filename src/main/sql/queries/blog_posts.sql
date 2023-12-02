-- name: CreateBlogPost :one
insert into blog_posts (id, created_at, updated_at, title, description, image_filename, user_id)
values ($1, $2, $3, $4, $5, $6, $7)
returning *;

-- name: GetBlogPosts :many
select * from blog_posts;

-- name: GetBlogPostsByUser :many
select * from blog_posts where user_id = $1;

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
