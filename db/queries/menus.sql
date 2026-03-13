-- name: GetMenu :one
SELECT id, restaurant_id, name, description, active, created_at
FROM menus
WHERE id = $1;

-- name: ListMenus :many
SELECT id, restaurant_id, name, description, active, created_at
FROM menus
ORDER BY created_at;

-- name: ListMenusByRestaurant :many
SELECT id, restaurant_id, name, description, active, created_at
FROM menus
WHERE restaurant_id = $1
ORDER BY created_at;

-- name: CreateMenu :one
INSERT INTO menus (id, restaurant_id, name, description, active, created_at)
VALUES ($1, $2, $3, $4, $5, NOW())
RETURNING id, restaurant_id, name, description, active, created_at;

-- name: UpdateMenu :one
UPDATE menus
SET name        = $2,
    description = $3,
    active      = $4
WHERE id = $1
RETURNING id, restaurant_id, name, description, active, created_at;

-- name: DeleteMenu :exec
DELETE FROM menus WHERE id = $1;
