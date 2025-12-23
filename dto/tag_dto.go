package dto

import "github.com/google/uuid"

type CreateTagRequest struct {
	Name   string    `json:"name"`
	UserID uuid.UUID `json:"user_id"` // Optional
}

type UpdateTagRequest struct {
	Name string `json:"name"`
}

type TagResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
