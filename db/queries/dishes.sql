-- name: GetDish :one
SELECT id, menu_id, name, compound, price, category,
       calories, proteins, fats, carbohydrates, weight_grams, created_at
FROM dishes
WHERE id = $1;

-- name: ListDishes :many
SELECT id, menu_id, name, compound, price, category,
       calories, proteins, fats, carbohydrates, weight_grams, created_at
FROM dishes
ORDER BY created_at;

-- name: ListDishesByMenu :many
SELECT id, menu_id, name, compound, price, category,
       calories, proteins, fats, carbohydrates, weight_grams, created_at
FROM dishes
WHERE menu_id = $1
ORDER BY created_at;

-- name: CreateDish :one
INSERT INTO dishes (menu_id, name, compound, price, category,
                    calories, proteins, fats, carbohydrates, weight_grams, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
RETURNING id, menu_id, name, compound, price, category,
          calories, proteins, fats, carbohydrates, weight_grams, created_at;

-- name: UpdateDish :one
UPDATE dishes
SET name          = $2,
    compound      = $3,
    price         = $4,
    category      = $5,
    calories      = $6,
    proteins      = $7,
    fats          = $8,
    carbohydrates = $9,
    weight_grams  = $10
WHERE id = $1
RETURNING id, menu_id, name, compound, price, category,
          calories, proteins, fats, carbohydrates, weight_grams, created_at;

-- name: DeleteDish :exec
DELETE FROM dishes WHERE id = $1;
