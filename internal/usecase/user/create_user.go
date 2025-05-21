package user

import (
	"context"

	"github.com/arthurdias01/gobid/internal/validator"
)

type CreateUserRequest struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func (req CreateUserRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator
	eval.CheckField(validator.NotBlank(req.UserName), "user_name", "User name cannot be blank")
	eval.CheckField(validator.NotBlank(req.Email), "email", "Email cannot be blank")
	eval.CheckField(validator.Matches(req.Email, validator.EmailRX), "email", "Email must be a valid email address")
	eval.CheckField(validator.NotBlank(req.Bio), "bio", "Bio cannot be blank")
	eval.CheckField(validator.MinChars(req.Bio, 10), "bio", "Bio must be at least 10 characters long")
	eval.CheckField(validator.MaxChars(req.Bio, 255), "bio", "Bio must be at most 100 characters long")
	eval.CheckField(validator.NotBlank(req.Password), "password", "Password cannot be blank")
	eval.CheckField(validator.MinChars(req.Password, 8), "password", "Password must be at least 8 characters long")
	return eval
}
