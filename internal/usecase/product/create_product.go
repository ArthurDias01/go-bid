package product

import (
	"context"
	"strconv"
	"time"

	"github.com/arthurdias01/gobid/internal/validator"
)

type CreateProductRequest struct {
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"base_price"`
	AuctionEnd  time.Time `json:"auction_end"`
}

const minAuctionDuration = 2 * time.Hour

func (req CreateProductRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator
	eval.CheckField(validator.NotBlank(req.ProductName), "product_name", "Product name cannot be blank")
	eval.CheckField(validator.NotBlank(req.Description), "description", "Description cannot be blank")
	eval.CheckField(validator.MinChars(req.Description, 10) &&
		validator.MaxChars(req.Description, 255),
		"description", "Description must be between 10 and 255 characters")
	eval.CheckField(validator.NotBlank(strconv.FormatFloat(req.BasePrice, 'f', -1, 64)), "base_price", "Base price cannot be blank")
	eval.CheckField(req.BasePrice > 0, "base_price", "Base price must be greater than 0")
	eval.CheckField(validator.NotBlank(req.AuctionEnd.Format(time.RFC3339)), "auction_end", "Auction end cannot be blank")
	eval.CheckField(req.AuctionEnd.Sub(time.Now()) >= minAuctionDuration, "auction_end", "Auction end must be at least 2h in the future")
	return eval
}
