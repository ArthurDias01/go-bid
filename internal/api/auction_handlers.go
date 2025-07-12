package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/arthurdias01/gobid/internal/jsonutils"
	"github.com/arthurdias01/gobid/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *API) handleSubscribeUserToAuction(w http.ResponseWriter, r *http.Request) {
	rawProductID := chi.URLParam(r, "product_id")
	slog.Info("[handleSubscribeUserToAuction]", "product_id", rawProductID)
	productID, err := uuid.Parse(rawProductID)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{"message": "Invalid product ID, must be a valid UUID"})
		return
	}

	_, err = api.ProductsService.GetProductByID(r.Context(), productID)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
				"message": "Product not found",
			})
			return
		}
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"message": "Failed to get product",
		})
		return
	}

	userId, ok := api.Sessions.Get(r.Context(), "AuthenticatedUserID").(uuid.UUID)
	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"message": "unauthorized",
		})
		return
	}

	api.AuctionLobby.Lock()
	room, ok := api.AuctionLobby.Rooms[productID]
	api.AuctionLobby.Unlock()

	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"message": "Auction has ended",
		})
		return
	}

	conn, err := api.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"message": "Could not upgrade connection to a websocket protocol",
		})
		return
	}

	client := services.NewClient(room, conn, userId)

	room.Register <- client
	// go client.ReadEventLoop()
	// go client.WriteEventLoop()
	for {

	}
}
