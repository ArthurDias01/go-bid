package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/arthurdias01/gobid/internal/services"
	"github.com/go-chi/chi/v5"
)

type API struct {
	Router          *chi.Mux
	UsersService    *services.UsersService
	ProductsService *services.ProductsService
	Sessions        *scs.SessionManager
}

// func (api *API) handleCreateUser(w http.ResponseWriter, r *http.Request) {

// }

// func (api *API) handleGetUser(w http.ResponseWriter, r *http.Request) {

// }
