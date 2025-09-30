package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *API) BindRoutes() {
	api.Router.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger, api.Sessions.LoadAndSave)

	// csrfMiddleware := csrf.Protect([]byte(os.Getenv("GOBID_CSRF_KEY")),
	// 	csrf.Path("/api/v1"),
	// 	csrf.Secure(false), // DEV ONLY
	// 	csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		log.Printf("CSRF Token Received: %s", r.Header.Get("X-CSRF-Token"))
	// 		log.Printf("Cookies Received: %v", r.Cookies())
	// 		http.Error(w, "Forbidden - CSRF token invalid", http.StatusForbidden)
	// 	})),
	// )
	// api.Router.Use(csrfMiddleware)

	// csrf Secure false for dev only**
	// GOBID_CSRF_KEY -> generate with
	// csrfMiddleware := csrf.Protect([]byte(os.Getenv("GOBID_CSRF_KEY")), csrf.Path("/api/v1"), csrf.Secure(false))

	// api.Router.Use(csrfMiddleware)

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// r.Get("/csrftoken", api.HangleGetCSRFToken)
			r.Route("/users", func(r chi.Router) {
				r.Post("/signup", api.handleSignupUser)
				r.Post("/signin", api.handleSigninUser)
				r.Group(func(r chi.Router) {
					r.Use(api.AuthMiddleware)
					r.Post("/signout", api.handleSignoutUser)
				})
			})
			r.Route("/products", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(api.AuthMiddleware)
					r.Post("/", api.handleCreateProduct)
					r.Get("/", api.handleGetProducts)
					// r.Get("/{id}", api.handleGetProduct)
					// r.Put("/{id}", api.handleUpdateProduct)
					// r.Delete("/{id}", api.handleDeleteProduct)
					r.Get("/ws/subscribe/{product_id}", api.handleSubscribeUserToAuction)
				})
			})
		})
	})
}
