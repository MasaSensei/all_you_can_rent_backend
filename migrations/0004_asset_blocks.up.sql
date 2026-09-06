CREATE TABLE IF NOT EXISTS asset_blocks (
  id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id   UUID        NOT NULL,
  asset_id    UUID        NOT NULL REFERENCES assets(id),
  start_date  DATE        NOT NULL,
  end_date    DATE        NOT NULL,
  reason      TEXT,
  created_by  UUID,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at  TIMESTAMPTZ,
  CONSTRAINT asset_blocks_date_check CHECK (end_date >= start_date)
);
CREATE INDEX asset_blocks_asset_date_idx ON asset_blocks (asset_id, start_date, end_date) WHERE deleted_at IS NULL;
