package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InvoiceStatus represents the payment state of an invoice.
type InvoiceStatus string

const (
	InvoiceStatusPending   InvoiceStatus = "PENDING"
	InvoiceStatusPaid      InvoiceStatus = "PAID"
	InvoiceStatusCancelled InvoiceStatus = "CANCELLED"
)

// Invoice represents a billing record linked to a user and reservation.
type Invoice struct {
	ID            string        `gorm:"type:uuid;primaryKey"        json:"id"`
	UserID        string        `gorm:"type:varchar(255);not null"  json:"user_id"`
	ReservationID string        `gorm:"type:varchar(255);not null"  json:"reservation_id"`
	Amount        float64       `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status        InvoiceStatus `gorm:"type:varchar(50);not null;default:'PENDING'" json:"status"`
	IssuedAt      time.Time     `gorm:"not null"                    json:"issued_at"`
	DueAt         time.Time     `gorm:"not null"                    json:"due_at"`
	PaidAt        *time.Time    `gorm:"default:null"                json:"paid_at,omitempty"`
	CreatedAt     time.Time     `gorm:"autoCreateTime"              json:"created_at"`
	UpdatedAt     time.Time     `gorm:"autoUpdateTime"              json:"updated_at"`
}

// BeforeCreate generates a UUID primary key before inserting a new Invoice record.
func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.NewString()
	}
	return nil
}

// TableName overrides the default GORM table name.
func (Invoice) TableName() string {
	return "invoices"
}
