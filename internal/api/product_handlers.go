package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/arthurdias01/gobid/internal/jsonutils"
	"github.com/arthurdias01/gobid/internal/services"
	"github.com/arthurdias01/gobid/internal/usecase/product"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *API) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[product.CreateProductRequest](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}
	userID, ok := api.Sessions.Get(r.Context(), "AuthenticatedUserID").(uuid.UUID)
	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	productId, err := api.ProductsService.CreateProduct(
		r.Context(),
		userID,
		data.ProductName,
		data.Description,
		data.BasePrice,
		data.AuctionEnd,
	)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "Failed to create product",
		})
		return
	}

	ctx, _ := context.WithDeadline(context.Background(), data.AuctionEnd)

	auctionRoom := services.NewActionRoom(ctx, productId, api.BidsService)

	go auctionRoom.Run()

	api.AuctionLobby.Lock()
	api.AuctionLobby.Rooms[productId] = auctionRoom
	api.AuctionLobby.Unlock()

	successMessage := fmt.Sprintf("Auction for product %s started successfuly", productId)

	jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"message":    successMessage,
		"product_id": productId,
	})
}

func (api *API) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := api.ProductsService.GetProducts(r.Context())
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "Failed to get products",
		})
		return
	}
	jsonutils.EncodeJson(w, r, http.StatusOK, products)
}

func (api *API) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	product, err := api.ProductsService.GetProductByID(r.Context(), uuid.MustParse(id))
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "Failed to get product",
		})
		return
	}
	jsonutils.EncodeJson(w, r, http.StatusOK, product)
}

func (api *API) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	data, problems, err := jsonutils.DecodeValidJson[product.UpdateProductRequest](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}
	product, err := api.ProductsService.UpdateProduct(r.Context(), uuid.MustParse(id), data.ProductName, data.Description, data.BasePrice, data.AuctionEnd)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "Failed to update product",
		})
		return
	}
	jsonutils.EncodeJson(w, r, http.StatusOK, product)
}

func (api *API) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := api.ProductsService.DeleteProduct(r.Context(), uuid.MustParse(id))
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "Failed to delete product",
		})
		return
	}
	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message": "Product deleted successfully",
	})
}
