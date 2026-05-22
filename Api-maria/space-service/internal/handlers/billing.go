package handlers

import (
	"net/http"
	"strconv"

	"coworking/space-service/internal/services"

	"github.com/gin-gonic/gin"
)

// BillingHandler groups HTTP handlers for the /billing resource.
type BillingHandler struct {
	svc *services.BillingService
}

// NewBillingHandler constructs a BillingHandler.
func NewBillingHandler(svc *services.BillingService) *BillingHandler {
	return &BillingHandler{svc: svc}
}

// ListInvoices handles GET /billing (admin only)
// @Summary List invoices
// @Description List invoices with pagination (admin only)
// @Tags Billing
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /billing [get]
func (h *BillingHandler) ListInvoices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	invoices, total, err := h.svc.GetInvoices(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve invoices: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"invoices": invoices,
			"total":    total,
			"page":     page,
			"limit":    limit,
		},
		"message": "Invoices retrieved successfully",
	})
}

// GetMyInvoices handles GET /billing/mine (auth required)
// @Summary Get my invoices
// @Description Returns invoices for the authenticated user
// @Tags Billing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /billing/mine [get]
func (h *BillingHandler) GetMyInvoices(c *gin.Context) {
	userID, _ := c.Get("user_id")

	invoices, err := h.svc.GetUserInvoices(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve your invoices: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    invoices,
		"message": "Your invoices retrieved successfully",
	})
}

// GetInvoice handles GET /billing/:id (auth required)
// @Summary Get invoice
// @Description Get an invoice by ID
// @Tags Billing
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /billing/{id} [get]
func (h *BillingHandler) GetInvoice(c *gin.Context) {
	id := c.Param("id")

	invoice, err := h.svc.GetInvoice(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "invoice not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Non-admin users may only view their own invoices.
	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")
	if role != "ADMIN" && invoice.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Access denied: you can only view your own invoices",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    invoice,
		"message": "Invoice retrieved successfully",
	})
}

// CreateInvoice handles POST /billing (internal use)
// @Summary Create invoice
// @Description Create a new invoice
// @Tags Billing
// @Accept json
// @Produce json
// @Param invoice body services.CreateInvoiceInput true "Invoice details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /billing [post]
func (h *BillingHandler) CreateInvoice(c *gin.Context) {
	var input services.CreateInvoiceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body: " + err.Error(),
		})
		return
	}

	invoice, err := h.svc.CreateInvoice(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    invoice,
		"message": "Invoice created successfully",
	})
}

// MarkAsPaid handles PATCH /billing/:id/pay (admin only)
// @Summary Mark invoice as paid
// @Description Mark an invoice as paid (admin only)
// @Tags Billing
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /billing/{id}/pay [patch]
func (h *BillingHandler) MarkAsPaid(c *gin.Context) {
	id := c.Param("id")

	invoice, err := h.svc.MarkAsPaid(id)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "invoice not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    invoice,
		"message": "Invoice marked as paid",
	})
}

// CancelInvoice handles PATCH /billing/:id/cancel (admin only)
// @Summary Cancel invoice
// @Description Cancel an invoice (admin only)
// @Tags Billing
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /billing/{id}/cancel [patch]
func (h *BillingHandler) CancelInvoice(c *gin.Context) {
	id := c.Param("id")

	invoice, err := h.svc.CancelInvoice(id)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "invoice not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    invoice,
		"message": "Invoice cancelled successfully",
	})
}
