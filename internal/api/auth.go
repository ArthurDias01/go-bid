package api

import (
	"net/http"

	"github.com/arthurdias01/gobid/internal/jsonutils"
	"github.com/gorilla/csrf"
)

// func (api *Api) HandleGetCSRFToken(w http.ResponseWriter, r *http.Request) {
// 	token := csrf.Token(r)
// 	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
// 		"csrf_token": token,
// 	})
// }

func (api *API) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), "AuthenticatedUserID") {
			jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{
				"message": "must be logged in",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (api *API) HangleGetCSRFToken(w http.ResponseWriter, r *http.Request) {
	token := csrf.Token(r)
	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"csrf_token": token,
	})
}
