-- name: CreateProduct :one
INSERT INTO products (seller_id, product_name, description, base_price, auction_end)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1 LIMIT 1;

-- name: GetAllProducts :many
SELECT * FROM products;

-- name: UpdateProduct :one
UPDATE products SET product_name = $2, description = $3, base_price = $4, auction_end = $5 WHERE id = $1 RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: GetProductsBySellerID :many
SELECT * FROM products WHERE seller_id = $1;

-- name: GetProductsByIsSold :many
SELECT * FROM products WHERE is_sold = $1;
