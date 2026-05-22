package handlers

import (
	"net/http"
	"time"

	"coworking/space-service/internal/services"

	"github.com/gin-gonic/gin"
)

// ReportHandler groups HTTP handlers for the /reports resource.
type ReportHandler struct {
	svc *services.ReportService
}

// NewReportHandler constructs a ReportHandler.
func NewReportHandler(svc *services.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

// parseDateRange extracts and validates the "from" and "to" query parameters.
// Both are expected in RFC3339 or 2006-01-02 formats.
func parseDateRange(c *gin.Context) (time.Time, time.Time, bool) {
	fromStr := c.DefaultQuery("from", "")
	toStr := c.DefaultQuery("to", "")

	layouts := []string{time.RFC3339, "2006-01-02"}

	parseDate := func(s, field string) (time.Time, bool) {
		if s == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": field + " query parameter is required (RFC3339 or YYYY-MM-DD)",
			})
			return time.Time{}, false
		}
		for _, l := range layouts {
			if t, err := time.Parse(l, s); err == nil {
				return t.UTC(), true
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid " + field + " format. Use RFC3339 or YYYY-MM-DD",
		})
		return time.Time{}, false
	}

	from, ok := parseDate(fromStr, "from")
	if !ok {
		return time.Time{}, time.Time{}, false
	}
	to, ok := parseDate(toStr, "to")
	if !ok {
		return time.Time{}, time.Time{}, false
	}

	if to.Before(from) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "'to' must be after 'from'",
		})
		return time.Time{}, time.Time{}, false
	}

	return from, to, true
}

// OccupancyReport handles GET /reports/occupancy (admin only)
// @Summary Occupancy report
// @Description Get occupancy report (admin only)
// @Tags Reports
// @Accept json
// @Produce json
// @Param from query string true "From date (RFC3339 or YYYY-MM-DD)"
// @Param to query string true "To date (RFC3339 or YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /reports/occupancy [get]
func (h *ReportHandler) OccupancyReport(c *gin.Context) {
	from, to, ok := parseDateRange(c)
	if !ok {
		return
	}

	report, err := h.svc.OccupancyReport(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate occupancy report: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"from":   from.Format(time.RFC3339),
			"to":     to.Format(time.RFC3339),
			"spaces": report,
		},
		"message": "Occupancy report generated successfully",
	})
}

// RevenueReport handles GET /reports/revenue (admin only)
// @Summary Revenue report
// @Description Get revenue report (admin only)
// @Tags Reports
// @Accept json
// @Produce json
// @Param from query string true "From date (RFC3339 or YYYY-MM-DD)"
// @Param to query string true "To date (RFC3339 or YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /reports/revenue [get]
func (h *ReportHandler) RevenueReport(c *gin.Context) {
	from, to, ok := parseDateRange(c)
	if !ok {
		return
	}

	report, err := h.svc.RevenueReport(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate revenue report: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"from":    from.Format(time.RFC3339),
			"to":      to.Format(time.RFC3339),
			"summary": report,
		},
		"message": "Revenue report generated successfully",
	})
}

// SpacesReport handles GET /reports/spaces (admin only)
// @Summary Spaces report
// @Description Get spaces report (admin only)
// @Tags Reports
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /reports/spaces [get]
func (h *ReportHandler) SpacesReport(c *gin.Context) {
	stats, err := h.svc.SpacesReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate spaces report: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"message": "Spaces report generated successfully",
	})
}
