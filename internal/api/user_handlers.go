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
			"error": "invalid username or email",
		})
		return
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"user_id": id,
	})

}

func (api *API) handleSigninUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (api *API) handleSignoutUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
