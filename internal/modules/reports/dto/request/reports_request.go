package request

import "time"

// GenerateReport queues a new report generation job.
type GenerateReport struct {
	Name            string            `json:"name" validate:"required,max=255"`
	ReportType      string            `json:"report_type" validate:"required,oneof=revenue asset_usage bookings customers expenses"`
	GeneratedFormat string            `json:"generated_format" validate:"required,oneof=pdf csv xlsx"`
	Parameters      map[string]string `json:"parameters"`
}

// TrackEvent records a single analytics event.
type TrackEvent struct {
	EventName     string            `json:"event_name" validate:"required,max=100"`
	EventCategory *string           `json:"event_category" validate:"omitempty,max=100"`
	EventData     map[string]string `json:"event_data"`
	Source        *string           `json:"source" validate:"omitempty,max=50"`
	CustomerID    *string           `json:"customer_id" validate:"omitempty,uuid"`
	OccurredAt    *time.Time        `json:"occurred_at"`
}

// DashboardFilter holds date range for aggregated dashboard queries.
type DashboardFilter struct {
	From    time.Time `json:"from" validate:"required"`
	To      time.Time `json:"to" validate:"required,gtfield=From"`
	GroupBy string    `json:"group_by" validate:"required,oneof=day week month"`
}
