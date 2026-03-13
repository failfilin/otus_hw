-- name: GetRestaurant :one
SELECT id, name, logo, active, phone, address, created_at
FROM restaurants
WHERE id = $1;

-- name: ListRestaurants :many
SELECT id, name, logo, active, phone, address, created_at
FROM restaurants
ORDER BY created_at;

-- name: CreateRestaurant :one
INSERT INTO restaurants (id, name, logo, active, phone, address, created_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW())
RETURNING id, name, logo, active, phone, address, created_at;

-- name: UpdateRestaurant :one
UPDATE restaurants
SET name    = $2,
    logo    = $3,
    active  = $4,
    phone   = $5,
    address = $6
WHERE id = $1
RETURNING id, name, logo, active, phone, address, created_at;

-- name: DeleteRestaurant :exec
DELETE FROM restaurants WHERE id = $1;
