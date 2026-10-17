package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(
	ctx context.Context,
	username, email, password, API_key string,
) (*User, error) {
	query := `
		INSERT INTO users (username, email, password, API_key)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, API_key, created_at;
	`
	var user User
	err := r.db.QueryRow(
		ctx, query, username, email, password, API_key,
	).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.API_key,
		&user.CreatedAt,
	)
	return &user, err
}

func (r *Repository) Get(
	ctx context.Context,
	email string,
) (*User, error) {
	query := `
		SELECT id, username, email, password, created_at
		FROM users
		WHERE email = $1
	`
	var user User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	return &user, err
}
