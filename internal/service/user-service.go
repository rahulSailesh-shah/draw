package service

import (
	"context"

	"draw/internal/db/repo"
	"draw/internal/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error)
}

type userService struct {
	queries *repo.Queries
	db      *pgxpool.Pool
}

func NewUserService(
	db *pgxpool.Pool,
	queries *repo.Queries,
) UserService {
	return &userService{
		db:      db,
		queries: queries,
	}
}
func (s *userService) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.UserResponse{
		ID: user.ID,
		Name: user.Name,
	}, nil
}