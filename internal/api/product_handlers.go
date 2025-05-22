package api

import (
	"net/http"

	"github.com/arthurdias01/gobid/internal/jsonutils"
	"github.com/arthurdias01/gobid/internal/usecase/product"
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

	id, err := api.ProductsService.CreateProduct(
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

	jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"message":    "Product created successfully",
		"product_id": id,
	})
}
