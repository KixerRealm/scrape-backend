-- name: CreateBlogPost :one
INSERT INTO blog_posts (id, created_at, updated_at, title, description, image_filename, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetBlogPosts :many
SELECT * FROM blog_posts;