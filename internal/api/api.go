package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/arthurdias01/gobid/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type API struct {
	Router          *chi.Mux
	UsersService    *services.UsersService
	ProductsService *services.ProductsService
	Sessions        *scs.SessionManager
	WsUpgrader      websocket.Upgrader
	AuctionLobby    services.AuctionLobby
	BidsService     services.BidsService
}
