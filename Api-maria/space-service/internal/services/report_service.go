package services

import (
	"time"

	"gorm.io/gorm"
)

// ReportService provides analytical reports via SQL aggregations.
type ReportService struct {
	db *gorm.DB
}

// NewReportService constructs a ReportService with the given GORM connection.
func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{db: db}
}

// OccupancyResult holds aggregated booking data for a single space in the
// requested time window.
type OccupancyResult struct {
	SpaceID       string  `json:"space_id"`
	SpaceName     string  `json:"space_name"`
	SpaceType     string  `json:"space_type"`
	TotalInvoices int64   `json:"total_invoices"`
	PaidInvoices  int64   `json:"paid_invoices"`
	OccupancyRate float64 `json:"occupancy_rate"` // paid / total * 100
}

// OccupancyReport returns per-space booking statistics for the given period.
// It joins invoices with reservations via reservation_id which embeds the
// space_id as a prefix (format: "<space_id>:<reservation_uuid>").
// When the reservation service is unavailable, the aggregation falls back to
// counting all invoices grouped by the space portion of reservation_id.
func (s *ReportService) OccupancyReport(from, to time.Time) ([]OccupancyResult, error) {
	var results []OccupancyResult

	query := `
		SELECT
			sp.id                                                    AS space_id,
			sp.name                                                  AS space_name,
			sp.type                                                  AS space_type,
			COUNT(inv.id)                                            AS total_invoices,
			COUNT(CASE WHEN inv.status = 'PAID' THEN 1 END)         AS paid_invoices,
			CASE
				WHEN COUNT(inv.id) = 0 THEN 0
				ELSE ROUND(
					COUNT(CASE WHEN inv.status = 'PAID' THEN 1 END)::numeric
					/ COUNT(inv.id)::numeric * 100, 2)
			END                                                      AS occupancy_rate
		FROM spaces sp
		LEFT JOIN invoices inv
			ON inv.reservation_id LIKE sp.id || '%'
			AND inv.issued_at BETWEEN ? AND ?
		GROUP BY sp.id, sp.name, sp.type
		ORDER BY total_invoices DESC
	`

	if err := s.db.Raw(query, from.UTC(), to.UTC()).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// RevenueResult carries revenue totals for a single space.
type RevenueResult struct {
	SpaceID      string  `json:"space_id"`
	SpaceName    string  `json:"space_name"`
	SpaceType    string  `json:"space_type"`
	TotalRevenue float64 `json:"total_revenue"`
	PaidRevenue  float64 `json:"paid_revenue"`
	PendingRevenue float64 `json:"pending_revenue"`
}

// RevenueSummary wraps the per-space breakdown with an overall total.
type RevenueSummary struct {
	TotalRevenue   float64         `json:"total_revenue"`
	PaidRevenue    float64         `json:"paid_revenue"`
	PendingRevenue float64         `json:"pending_revenue"`
	BySpace        []RevenueResult `json:"by_space"`
}

// RevenueReport returns a revenue summary for all spaces in the given period.
func (s *ReportService) RevenueReport(from, to time.Time) (*RevenueSummary, error) {
	var bySpace []RevenueResult

	query := `
		SELECT
			sp.id                                                           AS space_id,
			sp.name                                                         AS space_name,
			sp.type                                                         AS space_type,
			COALESCE(SUM(inv.amount), 0)                                    AS total_revenue,
			COALESCE(SUM(CASE WHEN inv.status = 'PAID' THEN inv.amount END), 0)    AS paid_revenue,
			COALESCE(SUM(CASE WHEN inv.status = 'PENDING' THEN inv.amount END), 0) AS pending_revenue
		FROM spaces sp
		LEFT JOIN invoices inv
			ON inv.reservation_id LIKE sp.id || '%'
			AND inv.issued_at BETWEEN ? AND ?
			AND inv.status != 'CANCELLED'
		GROUP BY sp.id, sp.name, sp.type
		ORDER BY paid_revenue DESC
	`

	if err := s.db.Raw(query, from.UTC(), to.UTC()).Scan(&bySpace).Error; err != nil {
		return nil, err
	}

	summary := &RevenueSummary{BySpace: bySpace}
	for _, r := range bySpace {
		summary.TotalRevenue += r.TotalRevenue
		summary.PaidRevenue += r.PaidRevenue
		summary.PendingRevenue += r.PendingRevenue
	}

	return summary, nil
}

// SpaceStats holds lifetime booking and revenue figures for a single space.
type SpaceStats struct {
	SpaceID       string  `json:"space_id"`
	SpaceName     string  `json:"space_name"`
	SpaceType     string  `json:"space_type"`
	Capacity      int     `json:"capacity"`
	PricePerHour  float64 `json:"price_per_hour"`
	IsActive      bool    `json:"is_active"`
	BookingCount  int64   `json:"booking_count"`
	TotalRevenue  float64 `json:"total_revenue"`
	PaidRevenue   float64 `json:"paid_revenue"`
}

// SpacesReport returns lifetime booking counts and revenue for every space.
func (s *ReportService) SpacesReport() ([]SpaceStats, error) {
	var stats []SpaceStats

	query := `
		SELECT
			sp.id                                                                  AS space_id,
			sp.name                                                                AS space_name,
			sp.type                                                                AS space_type,
			sp.capacity,
			sp.price_per_hour,
			sp.is_active,
			COUNT(inv.id)                                                          AS booking_count,
			COALESCE(SUM(inv.amount), 0)                                           AS total_revenue,
			COALESCE(SUM(CASE WHEN inv.status = 'PAID' THEN inv.amount END), 0)   AS paid_revenue
		FROM spaces sp
		LEFT JOIN invoices inv
			ON inv.reservation_id LIKE sp.id || '%'
			AND inv.status != 'CANCELLED'
		GROUP BY sp.id, sp.name, sp.type, sp.capacity, sp.price_per_hour, sp.is_active
		ORDER BY booking_count DESC
	`

	if err := s.db.Raw(query).Scan(&stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}
