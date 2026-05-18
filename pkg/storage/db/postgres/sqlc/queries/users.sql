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

-- name: IsAdminByUseID :one
SELECT admin
FROM users
WHERE ID = $1
  AND admin = true;

-- name: IsAdminByLogin :one
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
       reviews.created_at,
       outbox_reviews.last_error,
       outbox_reviews.result_out
from users
         left join reviews on users.id = reviews.userid
         left join tasks on reviews.task = tasks.id
         left join outbox_reviews on reviews.id = outbox_reviews.reviewid
    and outbox_reviews.userid = users.id
where users.id = $1;

-- name: GetReviewsAll :many
select users.id,
       users.login,
       users.name,
       tasks.id   as taskid,
       tasks.name as taskname,
       reviews.id as reviewid,
       reviews.status,
       reviews.created_at,
       reviews.updated_at
from users
         left join reviews on users.id = reviews.userid
         left join tasks on reviews.task = tasks.id
where users.admin = false;

-- name: UpdateReviewStatusByID :exec
UPDATE reviews
SET status = $1
WHERE id = $2;

-- name: CreateOutboxReview :exec
INSERT INTO outbox_reviews (event_id,
                            userid,
                            reviewid,
                            payload)
VALUES ($1, $2, $3, $4);

-- name: GetPendingOutboxReviews :many
SELECT *
FROM outbox_reviews
WHERE status = 'pending'
  AND (next_retry_at IS NULL OR next_retry_at <= NOW())
ORDER BY created_at
LIMIT 50 FOR UPDATE SKIP LOCKED;

-- name: GetOutboxReviewsBatch :many
SELECT *
FROM outbox_reviews
WHERE status = 'pending'
  AND attempts < max_attempts
  AND (next_retry_at IS NULL OR next_retry_at <= NOW())
ORDER BY created_at
LIMIT 50 FOR UPDATE SKIP LOCKED;

-- name: MarkOutboxReviewProcessingOne :exec
UPDATE outbox_reviews
SET status       = 'processing',
    processed_at = NOW()
WHERE id = $1;

-- name: MarkOutboxReviewFailed :exec
UPDATE outbox_reviews
SET attempts      = attempts + 1,
    next_retry_at = NOW() + ($2 * INTERVAL '1 second'),
    last_error    = $3
WHERE id = $1;

-- name: MarkOutboxReviewsProcessingMany :exec
UPDATE outbox_reviews
SET status = 'processing'
WHERE id = ANY ($1::bigint[]);

-- name: MarkOutboxReviewFailedFinal :exec
UPDATE outbox_reviews
SET status     = 'failed',
    last_error = $2
WHERE id = $1;

-- name: MarkOutboxReviewRetry :exec
UPDATE outbox_reviews
SET attempts      = attempts + 1,
    next_retry_at = NOW() + ($2 * INTERVAL '1 second'),
    last_error    = $3,
    status        = 'pending'
WHERE id = $1;

-- name: MarkOutboxReviewStatus :exec
UPDATE outbox_reviews
SET status     = $1,
    last_error = $2,
    result_out = $3
WHERE reviewid = $4
  AND userid = $5;

-- name: MarkReviewStatus :exec
UPDATE reviews
SET status = $1
WHERE id = $2
  AND userid = $3
  AND task = $4;

-- name: GetTaskNumberByID :one
SELECT id, name
from tasks
where id = $1;

-- name: GetTaskByID :one
SELECT id, name, target_file_name, target_file_validation, target_file_validation_language
from tasks
where id = $1;