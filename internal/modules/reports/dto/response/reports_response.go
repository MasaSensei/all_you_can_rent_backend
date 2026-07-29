package response

import "time"

// Report is the API-facing shape of entity.Report.
type Report struct {
	ID              string     `json:"id"`
	TenantID        string     `json:"tenant_id"`
	Name            string     `json:"name"`
	ReportType      string     `json:"report_type"`
	GeneratedFormat string     `json:"generated_format"`
	FileURL         *string    `json:"file_url,omitempty"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// RevenueDataPoint is one period bucket in a revenue time-series.
type RevenueDataPoint struct {
	Period  string  `json:"period"`
	Revenue float64 `json:"revenue"`
	Count   int     `json:"count"`
}

// AssetUtilizationItem shows booking rate for one asset.
type AssetUtilizationItem struct {
	AssetID   string  `json:"asset_id"`
	AssetName string  `json:"asset_name"`
	BookedDays int    `json:"booked_days"`
	TotalDays  int    `json:"total_days"`
	Utilization float64 `json:"utilization_pct"`
}

// Dashboard aggregates the key metrics shown on the dashboard.
type Dashboard struct {
	TotalRevenue    float64                `json:"total_revenue"`
	TotalBookings   int                    `json:"total_bookings"`
	ActiveCustomers int                    `json:"active_customers"`
	Revenue         []RevenueDataPoint     `json:"revenue"`
	TopAssets       []AssetUtilizationItem `json:"top_assets"`
}
