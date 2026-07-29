package notification

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	notifSvc "rentos-backend/internal/modules/notification/service"
)

// DeliveryWorker processes pending notification log entries and
// dispatches them to the appropriate external channel (email, whatsapp, sms).
// In-app notifications are written synchronously by NotificationService.Send()
// and do not go through this worker.
type DeliveryWorker struct {
	logSvc   notifSvc.LogService
	interval time.Duration
	log      zerolog.Logger
}

func NewDeliveryWorker(logSvc notifSvc.LogService, interval time.Duration, log zerolog.Logger) *DeliveryWorker {
	return &DeliveryWorker{logSvc: logSvc, interval: interval, log: log}
}

// Run blocks until ctx is cancelled.
func (w *DeliveryWorker) Run(ctx context.Context) {
	w.log.Info().Dur("interval", w.interval).Msg("notification delivery worker started")
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("notification delivery worker stopped")
			return
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}

func (w *DeliveryWorker) tick(ctx context.Context) {
	// NOTE: In production this would SELECT notification_logs WHERE
	// delivery_status='pending' ORDER BY created_at ASC LIMIT 50.
	// Keeping as a stub — the pattern is identical to the webhook dispatcher.
	// Wire up a PendingLogsRepository and plug in email/whatsapp/sms SDK clients.
	w.log.Debug().Msg("notification worker: tick (pending log poll not yet wired)")
}

// DispatchEmail is the email delivery stub. Replace body with real SMTP / SES call.
func (w *DeliveryWorker) DispatchEmail(_ context.Context, logID, to, subject, body string) error {
	// TODO: plug in email SDK (AWS SES, SendGrid, Mailgun, etc.)
	w.log.Info().Str("log_id", logID).Str("to", to).Str("subject", subject).Msg("notification worker: email dispatch (stub)")
	return w.logSvc.UpdateStatus(context.Background(), logID, "sent", nil)
}

// DispatchWhatsApp is the WhatsApp delivery stub. Replace body with real API call.
func (w *DeliveryWorker) DispatchWhatsApp(_ context.Context, logID, to, message string) error {
	// TODO: plug in WhatsApp Business API / Twilio / Zenziva, etc.
	w.log.Info().Str("log_id", logID).Str("to", to).Msg("notification worker: whatsapp dispatch (stub)")
	return w.logSvc.UpdateStatus(context.Background(), logID, "sent", nil)
}

// DispatchSMS is the SMS delivery stub.
func (w *DeliveryWorker) DispatchSMS(_ context.Context, logID, to, message string) error {
	// TODO: plug in SMS gateway (Twilio, Vonage, Zenziva, etc.)
	w.log.Info().Str("log_id", logID).Str("to", to).Msg("notification worker: sms dispatch (stub)")
	return w.logSvc.UpdateStatus(context.Background(), logID, "sent", nil)
}
