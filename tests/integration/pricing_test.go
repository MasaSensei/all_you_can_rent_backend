//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/config"
	pricingmodule "rentos-backend/internal/modules/pricing"
	"rentos-backend/pkg/testutil"
)

func TestPricingRule_CreateAndResolve(t *testing.T) {
	cfg, _ := config.Load()
	container, err := bootstrap.New(cfg)
	if err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	defer container.Close()

	pricing := pricingmodule.New(container)
	app := testApp(pricing.RegisterRoutes)

	tenantID := testutil.NewTenantID()
	assetID  := testutil.NewID()

	// ---- Create a per_day pricing rule ----
	body := map[string]any{
		"asset_id":      assetID,
		"name":          "Standard Day Rate",
		"rule_type":     "per_day",
		"value":         150000,
		"duration_unit": "day",
	}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/pricing-rules", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", tenantID)

	resp, err := app.Test(req, 5000)
	testutil.AssertNoError(t, err)

	if resp.StatusCode != fiber.StatusCreated {
		raw, _ := io.ReadAll(resp.Body)
		t.Fatalf("CreatePricingRule: want 201, got %d: %s", resp.StatusCode, raw)
	}

	var created struct {
		Data struct {
			ID    string  `json:"id"`
			Value float64 `json:"value"`
		} `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&created)
	testutil.AssertEqual(t, float64(150000), created.Data.Value, "rule value")
	t.Logf("Created pricing rule: %s", created.Data.ID)
}
