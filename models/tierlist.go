package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TierList struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Rows      []TierRow `gorm:"foreignKey:TierListID;constraint:OnDelete:CASCADE;" json:"rows"`
}

type TierRow struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TierListID uuid.UUID  `gorm:"type:uuid;not null;index" json:"tier_list_id"`
	Label      string     `gorm:"not null" json:"label"`
	Color      string     `gorm:"type:varchar(7);default:'#FFFFFF'" json:"color"` // Hex code
	SortOrder  int        `gorm:"not null" json:"sort_order"`
	Items      []TierItem `gorm:"foreignKey:TierRowID;constraint:OnDelete:CASCADE;" json:"items"`
}

type TierItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TierRowID uuid.UUID `gorm:"type:uuid;not null;index" json:"tier_row_id"`
	GameID    uuid.UUID `gorm:"type:uuid;not null;index" json:"game_id"`
	SortOrder int       `gorm:"not null" json:"sort_order"`

	Game Game `gorm:"foreignKey:GameID" json:"game"`
}

func (t *TierList) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}

func (t *TierRow) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}

func (t *TierItem) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}
