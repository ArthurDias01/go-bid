package services

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageKind int

const (
	// Requests
	PlaceBid MessageKind = iota
	// Success message
	SuccessfullyPlacedBid

	//Info
	NewBidPlaced

	// Errors
	FailedToPlaceBid

	// Auction Done
	AuctionFinished
)

type Message struct {
	Message string      `json:"action"`
	Kind    MessageKind `json:"data"`
	UserID  uuid.UUID   `json:"userId"`
	Amount  float64     `json:"amount"`
}

type AuctionLobby struct {
	sync.Mutex
	Rooms map[uuid.UUID]*AuctionRoom
}

type AuctionRoom struct {
	ID          uuid.UUID
	Context     context.Context
	Broadcast   chan Message
	Register    chan *Client
	Unregister  chan *Client
	Clients     map[uuid.UUID]*Client
	BidsService BidsService
}

type Client struct {
	Room   *AuctionRoom
	Conn   *websocket.Conn
	Send   chan Message
	UserId uuid.UUID
}

func NewActionRoom(ctx context.Context, id uuid.UUID, BidsService BidsService) *AuctionRoom {
	return &AuctionRoom{
		ID:          id,
		Broadcast:   make(chan Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     map[uuid.UUID]*Client{},
		Context:     ctx,
		BidsService: BidsService,
	}
}

func NewClient(room *AuctionRoom, conn *websocket.Conn, userId uuid.UUID) *Client {
	return &Client{
		Room:   room,
		Conn:   conn,
		Send:   make(chan Message, 512),
		UserId: userId,
	}
}

func (r *AuctionRoom) registerClient(c *Client) {
	slog.Info("New user Connected", "Client", c)
	r.Clients[c.UserId] = c
}

func (r *AuctionRoom) unregisterClient(c *Client) {
	slog.Info("New user Disconnected", "Client", c)
	delete(r.Clients, c.UserId)
}

func (r *AuctionRoom) broadcastMessage(m Message) {
	slog.Info("New message Received", "RoomID", r.ID, "message", m.Message)
	switch m.Kind {
	case PlaceBid:
		bid, err := r.BidsService.PlaceBid(r.Context, r.ID, m.UserID, m.Amount)
		if err != nil {
			if errors.Is(err, ErrBidIsTooLow) {
				if client, ok := r.Clients[m.UserID]; ok {
					client.Send <- Message{
						Kind:    FailedToPlaceBid,
						Message: ErrBidIsTooLow.Error(),
					}
				}
			}
			return
		}
		if client, ok := r.Clients[m.UserID]; ok {
			client.Send <- Message{
				Kind:    SuccessfullyPlacedBid,
				Message: "Successfully Placed Bid",
			}
		}
		for id, client := range r.Clients {
			newBidMessage := Message{
				Kind:    NewBidPlaced,
				Message: "A new bid was placed", Amount: bid.BidAmount,
			}
			if id == m.UserID {
				continue
			}
			client.Send <- newBidMessage
		}
	}
}

func (r *AuctionRoom) Run() {
	slog.Info("Auction started.", "AuctionID", r.ID)
	defer func() {
		close(r.Broadcast)
		close(r.Register)
		close(r.Unregister)
	}()
	for {
		select {
		case client := <-r.Register:
			r.registerClient(client)
		case client := <-r.Unregister:
			r.unregisterClient(client)
		case message := <-r.Broadcast:
			r.broadcastMessage(message)
		case <-r.Context.Done():
			slog.Info("Auction has ended.", "auctionID", r.ID)
			for _, client := range r.Clients {
				client.Send <- Message{
					Kind:    AuctionFinished,
					Message: "Auction has been finalized",
				}
			}
			return
		}
	}
}
