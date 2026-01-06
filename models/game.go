package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameStatus string

const (
	StatusBacklog   GameStatus = "backlog"
	StatusPlaying   GameStatus = "playing"
	StatusCompleted GameStatus = "completed"
	StatusDropped   GameStatus = "dropped"
)

type Game struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Title        string     `gorm:"not null" json:"title"`
	CoverURL     string     `json:"cover_url"`
	Genre        string     `json:"genre"`
	Status       GameStatus `gorm:"type:varchar(20);default:'backlog'" json:"status"`
	Platform     string     `gorm:"type:varchar(50)" json:"platform"` // Steam, Xbox, Epic, Switch, GOG, Pirated
	Platinum     bool       `gorm:"default:false" json:"platinum"`
	Score        *float64   `json:"score"` // Pointer to allow null (0.0 is a valid score)
	HoursPlayed  int        `json:"hours_played"`
	HLTBEstimate int        `json:"hltb_estimate"`
	ReleaseYear  int        `json:"release_year"`
	DateFinished *time.Time `json:"date_finished"`
	LastPlayedAt *time.Time `json:"last_played_at"`
	ReviewText   string     `gorm:"type:text" json:"review_text"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	Tags []*GameTag `gorm:"many2many:game_related_tags;" json:"tags"`
}

func (g *Game) BeforeCreate(tx *gorm.DB) (err error) {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return
}
