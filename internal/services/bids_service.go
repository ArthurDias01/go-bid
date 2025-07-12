package services

import (
	"context"
	"errors"
	"time"

	"github.com/arthurdias01/gobid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BidsService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewBidsService(bidsRepository *pgxpool.Pool) BidsService {
	return BidsService{
		pool:    bidsRepository,
		queries: pgstore.New(bidsRepository),
	}
}

var ErrAuctionEnded = errors.New("auction has ended")
var ErrBidIsTooLow = errors.New("bid is too low")

func (bs *BidsService) PlaceBid(ctx context.Context, productID uuid.UUID, bidderID uuid.UUID, bidAmount float64) (pgstore.Bid, error) {
	product, err := bs.queries.GetProductByID(ctx, productID)
	if err != nil {
		return pgstore.Bid{}, err
	}

	if product.AuctionEnd.Before(time.Now()) {
		return pgstore.Bid{}, ErrAuctionEnded
	}

	highestBid, err := bs.queries.GetHighestBidByProductID(ctx, productID)
	if err != nil {
		return pgstore.Bid{}, err
	}

	if product.BasePrice >= bidAmount || highestBid.BidAmount >= bidAmount {
		return pgstore.Bid{}, ErrBidIsTooLow
	}

	bid, err := bs.queries.CreateBid(ctx, pgstore.CreateBidParams{
		ProductID: productID,
		BidderID:  bidderID,
		BidAmount: bidAmount,
	})

	if err != nil {
		return pgstore.Bid{}, err
	}

	return bid, nil
}
