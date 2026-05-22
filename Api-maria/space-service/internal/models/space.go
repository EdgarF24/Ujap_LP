package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SpaceType represents the kind of co-working space.
type SpaceType string

const (
	SpaceTypeDesk        SpaceType = "DESK"
	SpaceTypeOffice      SpaceType = "OFFICE"
	SpaceTypeMeetingRoom SpaceType = "MEETING_ROOM"
	SpaceTypeLounge      SpaceType = "LOUNGE"
)

// Space represents a physical co-working space available for reservation.
type Space struct {
	ID           string    `gorm:"type:uuid;primaryKey"             json:"id"`
	Name         string    `gorm:"type:varchar(255);not null"       json:"name"`
	Type         SpaceType `gorm:"type:varchar(50);not null"        json:"type"`
	Capacity     int       `gorm:"not null"                         json:"capacity"`
	PricePerHour float64   `gorm:"type:decimal(10,2);not null"      json:"price_per_hour"`
	Amenities    string    `gorm:"type:text"                        json:"amenities"` // comma-separated
	IsActive     bool      `gorm:"default:true;not null"            json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime"                   json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"                   json:"updated_at"`
}

// BeforeCreate generates a UUID primary key before inserting a new Space record.
func (s *Space) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	return nil
}

// TableName overrides the default GORM table name.
func (Space) TableName() string {
	return "spaces"
}
