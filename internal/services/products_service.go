package services

import (
	"context"
	"errors"
	"time"

	"github.com/arthurdias01/gobid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (ps *ProductsService) GetProducts(ctx context.Context) ([]pgstore.Product, error) {
	products, err := ps.queries.GetAllProducts(ctx)
	if err != nil {
		return []pgstore.Product{}, err
	}
	return products, nil
}

var ErrProductNotFound = errors.New("product not found")

func (ps *ProductsService) GetProductByID(ctx context.Context, id uuid.UUID) (pgstore.Product, error) {
	product, err := ps.queries.GetProductByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return pgstore.Product{}, ErrProductNotFound
		}
		return pgstore.Product{}, err
	}
	return product, nil
}

func (ps *ProductsService) UpdateProduct(ctx context.Context, id uuid.UUID, productName string, description string, basePrice float64, auctionEnd time.Time) (pgstore.Product, error) {
	product, err := ps.queries.UpdateProduct(ctx, pgstore.UpdateProductParams{
		ID:          id,
		ProductName: productName,
		Description: description,
		BasePrice:   basePrice,
		AuctionEnd:  auctionEnd,
	})
	if err != nil {
		return pgstore.Product{}, err
	}
	return product, nil
}

func (ps *ProductsService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	err := ps.queries.DeleteProduct(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
