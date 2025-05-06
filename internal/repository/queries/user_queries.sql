-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = ? LIMIT 1;

-- name: GetUserById :one
SELECT * FROM users 
WHERE id = ? LIMIT 1;

-- name: UserExistsByEmail :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE email = ?
) AS user_exists;

-- name: CreateUser :one
INSERT INTO users (id, username, email, password_hash, user_role)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET username = ?, email = ?, password_hash = ?
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: DropAllUsers :exec
DELETE FROM users;
