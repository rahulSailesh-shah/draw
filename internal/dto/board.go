package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)


type Board struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	OwnerID string `json:"ownerId"`
	Elements json.RawMessage `json:"elements"`
}

// Request

type GetBoardRequest struct {
	BoardID string `json:"-"`
	UserID string `json:"-"`
}

type CreateBoardRequest struct {
	UserID string `json:"-"`
	Name string `json:"name" binding:"required"`
}

type UpdateBoardRequest struct {
	BoardID string `json:"-"`
	UserID string `json:"-"`
	Name string `json:"name,omitempty"`
	Elements json.RawMessage `json:"elements,omitempty"`
}

// Response
type CreateBoardResponse struct {
	BoardID uuid.UUID `json:"boardId"`
	Token string `json:"token"`
}

type GetBoardResponse struct {
	Board Board `json:"board"`
}