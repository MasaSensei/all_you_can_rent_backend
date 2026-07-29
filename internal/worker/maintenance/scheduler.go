package maintenance

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	maintreq "rentos-backend/internal/modules/maintenance/dto/request"
	maintSvc "rentos-backend/internal/modules/maintenance/service"
)

// Scheduler polls for maintenance records whose scheduled_date has arrived
// and transitions them from scheduled → in_progress.
type Scheduler struct {
	svc      maintSvc.MaintenanceService
	interval time.Duration
	log      zerolog.Logger
}

func NewScheduler(svc maintSvc.MaintenanceService, interval time.Duration, log zerolog.Logger) *Scheduler {
	return &Scheduler{svc: svc, interval: interval, log: log}
}

// Run blocks until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	s.log.Info().Dur("interval", s.interval).Msg("maintenance scheduler started")
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.log.Info().Msg("maintenance scheduler stopped")
			return
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	records, err := s.svc.ListDue(ctx, time.Now())
	if err != nil {
		s.log.Error().Err(err).Msg("maintenance scheduler: ListDue failed")
		return
	}
	if len(records) == 0 {
		return
	}

	s.log.Info().Int("count", len(records)).Msg("maintenance scheduler: processing due records")

	for _, rec := range records {
		req := maintreq.UpdateMaintenanceStatus{MaintenanceStatus: "in_progress"}
		if _, err := s.svc.UpdateStatus(ctx, rec.ID, rec.TenantID, "system", req); err != nil {
			s.log.Error().Err(err).Str("record_id", rec.ID).Msg("maintenance scheduler: UpdateStatus failed")
		}
	}
}
