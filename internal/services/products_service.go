package services

import (
	"context"
	"time"

	"github.com/arthurdias01/gobid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductsService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

type CreateProductParams struct {
	SellerID    uuid.UUID `json:"seller_id"`
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"base_price"`
	AuctionEnd  time.Time `json:"auction_end"`
}

func NewProductsService(productsRepository *pgxpool.Pool) *ProductsService {
	return &ProductsService{
		pool:    productsRepository,
		queries: pgstore.New(productsRepository),
	}
}

func (ps *ProductsService) CreateProduct(ctx context.Context,
	sellerID uuid.UUID,
	productName string,
	description string,
	basePrice float64,
	auctionEnd time.Time,
) (uuid.UUID, error) {
	id, err := ps.queries.CreateProduct(ctx, pgstore.CreateProductParams{
		SellerID:    sellerID,
		ProductName: productName,
		Description: description,
		BasePrice:   basePrice,
		AuctionEnd:  auctionEnd,
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}
