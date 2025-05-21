package api

import (
	"errors"
	"net/http"

	"github.com/arthurdias01/gobid/internal/jsonutils"
	"github.com/arthurdias01/gobid/internal/services"
	"github.com/arthurdias01/gobid/internal/usecase/user"
)

func (api *API) handleSignupUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[user.CreateUserRequest](r)
	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UsersService.CreateUser(r.Context(), data.UserName, data.Email, data.Password, data.Bio)
	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	if errors.Is(err, services.ErrDuplicatedEmailOrUsername) {
		_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
			"error": "Email or username already taken",
		})
		return
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"user_id": id,
	})

}

func (api *API) handleSigninUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[user.SignInUserRequest](r)
	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UsersService.AuthenticateUser(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": "Invalid email or password",
			})
			return
		}
		_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
			"error": "Internal server error",
		})
		return
	}

	err = api.Sessions.RenewToken(r.Context())
	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "Internal server error",
		})
		return
	}

	api.Sessions.Put(r.Context(), "AuthenticatedUserID", id)
	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message": "User signed in successfully",
	})
}

func (api *API) handleSignoutUser(w http.ResponseWriter, r *http.Request) {
	err := api.Sessions.RenewToken(r.Context())
	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "Internal server error",
		})
		return
	}
	api.Sessions.Remove(r.Context(), "AuthenticatedUserID")

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message": "User signed out successfully",
	})
}
