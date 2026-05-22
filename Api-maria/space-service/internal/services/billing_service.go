package services

import (
	"errors"
	"time"

	"coworking/space-service/internal/models"
	"coworking/space-service/internal/repository"
)

// BillingService encapsulates business logic for invoice management.
type BillingService struct {
	repo repository.InvoiceRepository
}

// NewBillingService constructs a BillingService with the provided repository.
func NewBillingService(repo repository.InvoiceRepository) *BillingService {
	return &BillingService{repo: repo}
}

// CreateInvoiceInput holds the validated fields required to create an invoice.
type CreateInvoiceInput struct {
	UserID        string    `json:"user_id"         binding:"required"`
	ReservationID string    `json:"reservation_id"  binding:"required"`
	Amount        float64   `json:"amount"          binding:"required,gt=0"`
	IssuedAt      time.Time `json:"issued_at"`
	DueAt         time.Time `json:"due_at"          binding:"required"`
}

// CreateInvoice validates and persists a new Invoice in PENDING status.
func (s *BillingService) CreateInvoice(input CreateInvoiceInput) (*models.Invoice, error) {
	if input.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	issuedAt := input.IssuedAt
	if issuedAt.IsZero() {
		issuedAt = time.Now().UTC()
	}

	if !input.DueAt.IsZero() && input.DueAt.Before(issuedAt) {
		return nil, errors.New("due_at must be after issued_at")
	}

	invoice := &models.Invoice{
		UserID:        input.UserID,
		ReservationID: input.ReservationID,
		Amount:        input.Amount,
		Status:        models.InvoiceStatusPending,
		IssuedAt:      issuedAt,
		DueAt:         input.DueAt,
	}

	if err := s.repo.Create(invoice); err != nil {
		return nil, err
	}
	return invoice, nil
}

// GetInvoices returns a paginated list of all invoices along with total count.
func (s *BillingService) GetInvoices(page, limit int) ([]models.Invoice, int64, error) {
	pagination := repository.Pagination{Page: page, Limit: limit}
	return s.repo.GetAll(pagination)
}

// GetInvoice fetches a single invoice by ID, returning an error if not found.
func (s *BillingService) GetInvoice(id string) (*models.Invoice, error) {
	invoice, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, errors.New("invoice not found")
	}
	return invoice, nil
}

// GetUserInvoices returns all invoices for a specific user.
func (s *BillingService) GetUserInvoices(userID string) ([]models.Invoice, error) {
	return s.repo.GetByUserID(userID)
}

// MarkAsPaid transitions an invoice to PAID status and sets PaidAt to now.
func (s *BillingService) MarkAsPaid(id string) (*models.Invoice, error) {
	invoice, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, errors.New("invoice not found")
	}
	if invoice.Status == models.InvoiceStatusPaid {
		return nil, errors.New("invoice is already paid")
	}
	if invoice.Status == models.InvoiceStatusCancelled {
		return nil, errors.New("cannot pay a cancelled invoice")
	}

	now := time.Now().UTC()
	if err := s.repo.UpdateStatus(id, models.InvoiceStatusPaid, &now); err != nil {
		return nil, err
	}

	invoice.Status = models.InvoiceStatusPaid
	invoice.PaidAt = &now
	return invoice, nil
}

// CancelInvoice transitions an invoice to CANCELLED status.
func (s *BillingService) CancelInvoice(id string) (*models.Invoice, error) {
	invoice, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, errors.New("invoice not found")
	}
	if invoice.Status == models.InvoiceStatusCancelled {
		return nil, errors.New("invoice is already cancelled")
	}
	if invoice.Status == models.InvoiceStatusPaid {
		return nil, errors.New("cannot cancel a paid invoice")
	}

	if err := s.repo.UpdateStatus(id, models.InvoiceStatusCancelled, nil); err != nil {
		return nil, err
	}

	invoice.Status = models.InvoiceStatusCancelled
	return invoice, nil
}
