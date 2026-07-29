package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"rentos-backend/internal/modules/integration/repository"
	"rentos-backend/pkg/database"
)

// Dispatcher polls pending webhook_log entries and delivers them via HTTP POST.
// On success it updates delivery_status=sent; on failure it marks failed.
type Dispatcher struct {
	webhookRepo repository.WebhookRepository
	logRepo     repository.WebhookLogRepository
	db          database.Querier
	httpClient  *http.Client
	interval    time.Duration
	log         zerolog.Logger
}

func NewDispatcher(
	webhookRepo repository.WebhookRepository,
	logRepo repository.WebhookLogRepository,
	db database.Querier,
	interval time.Duration,
	log zerolog.Logger,
) *Dispatcher {
	return &Dispatcher{
		webhookRepo: webhookRepo,
		logRepo:     logRepo,
		db:          db,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		interval:    interval,
		log:         log,
	}
}

// Run blocks until ctx is cancelled.
func (d *Dispatcher) Run(ctx context.Context) {
	d.log.Info().Dur("interval", d.interval).Msg("webhook dispatcher started")
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			d.log.Info().Msg("webhook dispatcher stopped")
			return
		case <-ticker.C:
			d.tick(ctx)
		}
	}
}

func (d *Dispatcher) tick(ctx context.Context) {
	// NOTE: In a production system, this would SELECT pending logs from
	// webhook_logs WHERE delivery_status = 'pending' ORDER BY created_at ASC LIMIT 50.
	// Wired via a dedicated PendingLogsRepository. Kept as a stub here to
	// show the delivery pattern without adding another SQL file.
	d.log.Debug().Msg("webhook dispatcher: tick (pending log poll not yet wired)")
}

// Deliver performs the signed HTTP POST for one log entry.
// Called directly for testing or by a future queue consumer.
func (d *Dispatcher) Deliver(ctx context.Context, webhookID, logID, targetURL, secret, payload string) error {
	sig := signPayload(secret, payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewBufferString(payload))
	if err != nil {
		return fmt.Errorf("dispatcher.Deliver: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-RentOS-Signature", "sha256="+sig)
	req.Header.Set("X-RentOS-Event", logID)

	resp, err := d.httpClient.Do(req)
	if err != nil {
		errMsg := err.Error()
		_ = d.logRepo.UpdateStatus(ctx, d.db, logID, "failed", &errMsg)
		return fmt.Errorf("dispatcher.Deliver: HTTP POST: %w", err)
	}
	defer resp.Body.Close()

	code := resp.StatusCode
	if code >= 200 && code < 300 {
		_ = d.logRepo.UpdateStatus(ctx, d.db, logID, "sent", nil)
	} else {
		msg := fmt.Sprintf("non-2xx response: %d", code)
		_ = d.logRepo.UpdateStatus(ctx, d.db, logID, "failed", &msg)
	}
	return nil
}

// signPayload creates an HMAC-SHA256 signature for the payload.
// Receivers verify: sha256=<hex> against X-RentOS-Signature header.
func signPayload(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
