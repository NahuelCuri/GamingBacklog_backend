package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameTag struct {
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;index" json:"user_id"` // Optional: if tags are per user
	Name   string    `gorm:"uniqueIndex:idx_name_user;not null" json:"name"`
	Games  []*Game   `gorm:"many2many:game_related_tags;" json:"games,omitempty"`
}

func (t *GameTag) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}
