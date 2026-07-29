package report

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	reportSvc "rentos-backend/internal/modules/reports/service"
	"rentos-backend/internal/modules/reports/entity"
)

// Generator polls for queued report jobs, processes them, then updates
// status to completed (or failed). Real file generation (PDF/XLSX/CSV)
// is plugged in per report_type; this skeleton covers the orchestration.
type Generator struct {
	svc      reportSvc.ReportService
	interval time.Duration
	log      zerolog.Logger
}

func NewGenerator(svc reportSvc.ReportService, interval time.Duration, log zerolog.Logger) *Generator {
	return &Generator{svc: svc, interval: interval, log: log}
}

// Run blocks until ctx is cancelled.
func (g *Generator) Run(ctx context.Context) {
	g.log.Info().Dur("interval", g.interval).Msg("report generator started")
	ticker := time.NewTicker(g.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			g.log.Info().Msg("report generator stopped")
			return
		case <-ticker.C:
			g.tick(ctx)
		}
	}
}

func (g *Generator) tick(ctx context.Context) {
	// List queued reports (page 1, up to 10 at a time).
	reports, err := g.svc.List(ctx, "", 1, 10) // empty tenantID = all tenants
	if err != nil {
		g.log.Error().Err(err).Msg("report generator: List failed")
		return
	}

	for _, rep := range reports {
		if rep.Status != entity.ReportStatusQueued {
			continue
		}
		g.process(ctx, rep.ID, rep.TenantID, rep.ReportType, rep.GeneratedFormat)
	}
}

func (g *Generator) process(ctx context.Context, id, tenantID, reportType, format string) {
	g.log.Info().Str("id", id).Str("type", reportType).Msg("report generator: processing")

	// Mark as processing.
	_ = g.svc.UpdateStatus(ctx, id, tenantID, entity.ReportStatusProcessing, nil)

	// TODO: plug in real generators per report_type + format.
	// For now, set a placeholder file URL to mark completion.
	fileURL := fmt.Sprintf("/storage/reports/%s.%s", id, format)

	if err := g.svc.UpdateStatus(ctx, id, tenantID, entity.ReportStatusCompleted, &fileURL); err != nil {
		g.log.Error().Err(err).Str("id", id).Msg("report generator: UpdateStatus completed failed")
		failed := entity.ReportStatusFailed
		_ = g.svc.UpdateStatus(ctx, id, tenantID, failed, nil)
		return
	}

	g.log.Info().Str("id", id).Str("file_url", fileURL).Msg("report generator: completed")
}
