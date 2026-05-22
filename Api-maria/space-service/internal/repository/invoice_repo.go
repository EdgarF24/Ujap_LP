package repository

import (
	"errors"

	"coworking/space-service/internal/models"

	"gorm.io/gorm"
)

// Pagination holds limit/offset parameters for list queries.
type Pagination struct {
	Page  int
	Limit int
}

// Offset computes the SQL OFFSET from the page number and limit.
func (p *Pagination) Offset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.Limit
}

// InvoiceRepository defines the data-access contract for Invoice records.
type InvoiceRepository interface {
	Create(invoice *models.Invoice) error
	GetAll(pagination Pagination) ([]models.Invoice, int64, error)
	GetByID(id string) (*models.Invoice, error)
	GetByUserID(userID string) ([]models.Invoice, error)
	UpdateStatus(id string, status models.InvoiceStatus, paidAt interface{}) error
}

type invoiceRepo struct {
	db *gorm.DB
}

// NewInvoiceRepository returns a concrete InvoiceRepository backed by the
// given GORM database connection.
func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepo{db: db}
}

// Create inserts a new Invoice record.
func (r *invoiceRepo) Create(invoice *models.Invoice) error {
	return r.db.Create(invoice).Error
}

// GetAll returns a paginated list of all invoices and the total row count.
func (r *invoiceRepo) GetAll(pagination Pagination) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	query := r.db.Model(&models.Invoice{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if pagination.Limit <= 0 {
		pagination.Limit = 20
	}

	if err := query.Order("created_at DESC").
		Limit(pagination.Limit).
		Offset(pagination.Offset()).
		Find(&invoices).Error; err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

// GetByID fetches a single Invoice by its UUID primary key.
func (r *invoiceRepo) GetByID(id string) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := r.db.Where("id = ?", id).First(&invoice).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &invoice, nil
}

// GetByUserID returns all invoices belonging to a specific user.
func (r *invoiceRepo) GetByUserID(userID string) ([]models.Invoice, error) {
	var invoices []models.Invoice
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}

// UpdateStatus changes the status (and optionally PaidAt) of an invoice.
// paidAt should be a *time.Time or nil.
func (r *invoiceRepo) UpdateStatus(id string, status models.InvoiceStatus, paidAt interface{}) error {
	updates := map[string]interface{}{
		"status":  status,
		"paid_at": paidAt,
	}
	result := r.db.Model(&models.Invoice{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
