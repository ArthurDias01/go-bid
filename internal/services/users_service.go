package services

import (
	"context"
	"errors"

	"github.com/arthurdias01/gobid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

var ErrDuplicatedEmailOrUsername = errors.New("invalid username or email")

func NewUsersService(pool *pgxpool.Pool) *UsersService {
	return &UsersService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (us *UsersService) CreateUser(ctx context.Context, userName, email, password, bio string) (uuid.UUID, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return uuid.UUID{}, err
	}

	args := pgstore.CreateUserParams{
		UserName:     userName,
		Email:        email,
		PasswordHash: hash,
		Bio:          bio,
	}

	id, err := us.queries.CreateUser(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return uuid.UUID{}, ErrDuplicatedEmailOrUsername
			}
		}
	}

	return id, nil
}
