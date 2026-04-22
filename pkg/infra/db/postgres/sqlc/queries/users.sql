-- name: CreateUser :one
INSERT INTO users (login, name, password)
VALUES ($1, $2, $3) RETURNING id, login, name, password;

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