package model

import (
	"time"

	"gorm.io/gorm"
)

type SourceRepository struct {
	ID        uint32     `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Link      string     `json:"link,omitempty"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (sr *SourceRepository) BeforeUpdate(tx *gorm.DB) (err error) {
	current := time.Now()
	sr.UpdatedAt = &current
	return nil
}
