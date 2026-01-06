package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateGameRequest struct {
	UserID       uuid.UUID   `json:"user_id"`
	Title        string      `json:"title"`
	CoverURL     string      `json:"cover_url"`
	Genre        string      `json:"genre"`
	Status       string      `json:"status"` // "backlog", "playing", etc.
	Platform     string      `json:"platform"`
	Platinum     bool        `json:"platinum"`
	Score        *float64    `json:"score"`
	HoursPlayed  int         `json:"hours_played"`
	HLTBEstimate int         `json:"hltb_estimate"`
	ReleaseYear  int         `json:"release_year"`
	DateFinished *time.Time  `json:"date_finished"`
	ReviewText   string      `json:"review_text"`
	TagIDs       []uuid.UUID `json:"tag_ids"`
}

type UpdateGameRequest struct {
	Title        string      `json:"title"`
	CoverURL     string      `json:"cover_url"`
	Genre        string      `json:"genre"`
	Status       string      `json:"status"`
	Platform     string      `json:"platform"`
	Platinum     *bool       `json:"platinum"`
	Score        *float64    `json:"score"`
	HoursPlayed  int         `json:"hours_played"`
	HLTBEstimate int         `json:"hltb_estimate"`
	ReleaseYear  int         `json:"release_year"`
	DateFinished *time.Time  `json:"date_finished"`
	ReviewText   string      `json:"review_text"`
	TagIDs       []uuid.UUID `json:"tag_ids"` // Replace tags
}

type GameResponse struct {
	ID           uuid.UUID     `json:"id"`
	UserID       uuid.UUID     `json:"user_id"`
	Title        string        `json:"title"`
	CoverURL     string        `json:"cover_url"`
	Genre        string        `json:"genre"`
	Status       string        `json:"status"`
	Platform     string        `json:"platform"`
	Platinum     bool          `json:"platinum"`
	Score        *float64      `json:"score"`
	HoursPlayed  int           `json:"hours_played"`
	HLTBEstimate int           `json:"hltb_estimate"`
	ReleaseYear  int           `json:"release_year"`
	DateFinished *time.Time    `json:"date_finished"`
	LastPlayedAt *time.Time    `json:"last_played_at"`
	ReviewText   string        `json:"review_text"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Tags         []TagResponse `json:"tags"`
}
