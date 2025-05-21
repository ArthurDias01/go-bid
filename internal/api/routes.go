package api

import "github.com/go-chi/chi/v5"

func (api *API) BindRoutes() {
	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Post("/signup", api.handleSignupUser)
				r.Post("/signin", api.handleSigninUser)
				r.Post("/signout", api.handleSignoutUser)
			})
		})
	})
}
