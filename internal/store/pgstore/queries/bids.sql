-- name: CreateBid :one
INSERT INTO bids (product_id, bidder_id, bid_amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBidByID :one
SELECT * FROM bids WHERE id = $1 LIMIT 1;

-- name: GetBidsByProductID :many
SELECT * FROM bids WHERE product_id = $1;

-- name: GetBidsByUserID :many
SELECT * FROM bids WHERE bidder_id = $1;

-- name: GetHighestBidByProductID :one
SELECT * FROM bids WHERE product_id = $1 ORDER BY bid_amount DESC LIMIT 1;