package user

import (
	"context"

	"github.com/arthurdias01/gobid/internal/validator"
)

type SignInUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInUserResponse struct {
	Token string `json:"token"`
}

func (req SignInUserRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator
	eval.CheckField(validator.NotBlank(req.Email), "email", "Email is required")
	eval.CheckField(validator.Matches(req.Email, validator.EmailRX), "email", "Email must be a valid email address")
	eval.CheckField(validator.NotBlank(req.Password), "password", "Password is required")
	return eval
}
