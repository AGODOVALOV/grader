-- name: CreateUser :one
INSERT INTO users (login, name, password)
VALUES ($1, $2, $3)
RETURNING id, login, name, password;

-- name: GetUserByID :one
SELECT id, login, name, password
FROM users
WHERE id = $1;

-- name: GetUserByLogin :one
SELECT id, login, name, password
FROM users
WHERE login = $1;

-- name: ListUsers :many
SELECT id, login, name, password
FROM users
ORDER BY id;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;

-- name: IsAdmin :one
SELECT admin
FROM users
WHERE login = $1
  AND admin = true;

-- name: GetReviewByID :one
SELECT *
FROM reviews
WHERE id = $1;

-- name: GetReviewsByUser :many
SELECT *
FROM reviews
WHERE userid = $1
ORDER BY created_at DESC;

-- name: GetReviewsByTask :many
SELECT *
FROM reviews
WHERE task = $1
ORDER BY created_at DESC;

-- name: GetPendingReviews :many
SELECT *
FROM reviews
WHERE status = 'pending'
ORDER BY created_at ASC;

-- name: UpdateReviewStatus :exec
UPDATE reviews
SET status   = $2,
    attempts = attempts + 1
WHERE id = $1;

-- name: GetReviewfileID :one
SELECT fileID
FROM reviews
WHERE id = $1;

-- name: CreateReview :one
INSERT INTO reviews (userid, task, fileID)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetReviewsByUserID :many
select users.id,
       users.login,
       users.name,
       tasks.id   as taskid,
       tasks.name as taskname,
       reviews.id as reviewid,
       reviews.status,
       reviews.created_at
from users
         left join reviews on users.id = reviews.userid
         left join tasks on reviews.task = tasks.id
where users.id = $1;
