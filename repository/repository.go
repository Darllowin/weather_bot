package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Darllowin/weather_bot/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetUserCity(ctx context.Context, userID int64) (string, error) {
	var city string
	row := r.db.QueryRow(ctx, "SELECT coalesce(city, '') FROM users WHERE id=$1", userID)
	err := row.Scan(&city)
	if err != nil {
		return "", fmt.Errorf("error row.Scan: %w", err)
	}

	return city, nil
}

func (r *Repository) CreateUser(ctx context.Context, userID int64) error {
	_, err := r.db.Exec(ctx, "INSERT INTO users (id) VALUES ($1)", userID)
	if err != nil {
		return fmt.Errorf("error db.Exec: %w", err)
	}

	return nil
}

func (r *Repository) UpdateCity(ctx context.Context, userID int64, city string) error {
	_, err := r.db.Exec(ctx, "UPDATE users SET city = $1 WHERE id = $2", city, userID)
	if err != nil {
		return fmt.Errorf("error db.Exec: %w", err)
	}

	return nil
}

func (r *Repository) GetUser(ctx context.Context, userID int64) (*models.User, error) {
	user := models.User{}
	row := r.db.QueryRow(ctx, "SELECT id, coalesce(city, ''), created_at FROM users WHERE id = $1", userID)
	err := row.Scan(&user.ID, &user.City, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error row.Scan: %w", err)
	}

	return &user, nil
}
