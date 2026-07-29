-- 0001_init.up.sql
-- Full RentOS schema, exactly as approved in the ER diagram. Audit columns
-- (created_by/updated_by/deleted_by) are intentionally stored as plain
-- UUID without a foreign key to users(id): enforcing that FK would create
-- circular dependencies across nearly every table and complicate user
-- deletion; the application layer is the source of truth for that
-- relationship, consistent across every module.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =========================================================================
-- CORE / MULTI-TENANT
-- =========================================================================

CREATE TABLE tenants (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(100) NOT NULL UNIQUE,
    domain      VARCHAR(255),
    plan        VARCHAR(50) NOT NULL DEFAULT 'free',
    status      VARCHAR(30) NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE system_settings (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key         VARCHAR(100) NOT NULL,
    value       TEXT NOT NULL,
    type        VARCHAR(20) NOT NULL DEFAULT 'string',
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_system_settings_key ON system_settings (key) WHERE deleted_at IS NULL;

CREATE TABLE settings (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    key         VARCHAR(100) NOT NULL,
    value       TEXT NOT NULL,
    type        VARCHAR(20) NOT NULL DEFAULT 'string',
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by  UUID,
    updated_by  UUID,
    deleted_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_settings_tenant_key ON settings (tenant_id, key) WHERE deleted_at IS NULL;
CREATE INDEX idx_settings_tenant ON settings (tenant_id);

CREATE TABLE audit_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id     UUID,
    entity_type VARCHAR(100) NOT NULL,
    entity_id   UUID NOT NULL,
    action      VARCHAR(50) NOT NULL,
    old_values  JSONB,
    new_values  JSONB,
    ip_address  VARCHAR(64),
    user_agent  VARCHAR(255),
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_audit_logs_tenant ON audit_logs (tenant_id, created_at DESC);
CREATE INDEX idx_audit_logs_entity ON audit_logs (entity_type, entity_id);

-- =========================================================================
-- AUTHENTICATION
-- =========================================================================

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email         VARCHAR(255) NOT NULL,
    username      VARCHAR(100),
    password_hash VARCHAR(255) NOT NULL,
    first_name    VARCHAR(100),
    last_name     VARCHAR(100),
    phone         VARCHAR(30),
    avatar_url    VARCHAR(500),
    is_active     BOOLEAN NOT NULL DEFAULT true,
    last_login_at TIMESTAMPTZ,
    status        VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by    UUID,
    updated_by    UUID,
    deleted_by    UUID,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ,
    version       INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_users_tenant_email ON users (tenant_id, email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_tenant ON users (tenant_id);

CREATE TABLE user_sessions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token         VARCHAR(500) NOT NULL,
    refresh_token VARCHAR(500) NOT NULL,
    ip_address    VARCHAR(64),
    user_agent    VARCHAR(255),
    expires_at    TIMESTAMPTZ NOT NULL,
    status        VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ,
    version       INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_user_sessions_user ON user_sessions (user_id);
CREATE INDEX idx_user_sessions_refresh_token ON user_sessions (refresh_token);

CREATE TABLE password_resets (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token       VARCHAR(255) NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    used        BOOLEAN NOT NULL DEFAULT false,
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_password_resets_token ON password_resets (token);

CREATE TABLE auth_providers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    type        VARCHAR(50) NOT NULL,
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE user_auth_providers (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    auth_provider_id UUID NOT NULL REFERENCES auth_providers(id) ON DELETE RESTRICT,
    provider_user_id VARCHAR(255) NOT NULL,
    access_token     VARCHAR(1000),
    refresh_token    VARCHAR(1000),
    status           VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at       TIMESTAMPTZ,
    version          INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_uap_provider_user ON user_auth_providers (auth_provider_id, provider_user_id) WHERE deleted_at IS NULL;

-- =========================================================================
-- ROLES & PERMISSIONS
-- =========================================================================

CREATE TABLE roles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    is_system   BOOLEAN NOT NULL DEFAULT false,
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by  UUID,
    updated_by  UUID,
    deleted_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_roles_tenant_name ON roles (tenant_id, name) WHERE deleted_at IS NULL;

CREATE TABLE permissions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(150) NOT NULL UNIQUE,
    module      VARCHAR(100) NOT NULL,
    action      VARCHAR(50) NOT NULL,
    description TEXT,
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE role_permissions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    role_id       UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    status        VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ,
    version       INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_role_permissions ON role_permissions (role_id, permission_id) WHERE deleted_at IS NULL;

CREATE TABLE user_roles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id     UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_user_roles ON user_roles (user_id, role_id) WHERE deleted_at IS NULL;

-- =========================================================================
-- INVENTORY
-- =========================================================================

CREATE TABLE categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    parent_id   UUID REFERENCES categories(id) ON DELETE SET NULL,
    name        VARCHAR(150) NOT NULL,
    slug        VARCHAR(150) NOT NULL,
    description TEXT,
    icon        VARCHAR(255),
    sort_order  INTEGER NOT NULL DEFAULT 0,
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by  UUID,
    updated_by  UUID,
    deleted_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_categories_tenant_slug ON categories (tenant_id, slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_categories_tenant ON categories (tenant_id);
CREATE INDEX idx_categories_parent ON categories (parent_id);

CREATE TABLE asset_templates (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    name        VARCHAR(150) NOT NULL,
    description TEXT,
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by  UUID,
    updated_by  UUID,
    deleted_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_asset_templates_tenant ON asset_templates (tenant_id);

CREATE TABLE template_fields (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    asset_template_id UUID NOT NULL REFERENCES asset_templates(id) ON DELETE CASCADE,
    field_name        VARCHAR(100) NOT NULL,
    field_label       VARCHAR(150) NOT NULL,
    field_type        VARCHAR(30) NOT NULL,
    is_required       BOOLEAN NOT NULL DEFAULT false,
    default_value     TEXT,
    options           JSONB,
    sort_order        INTEGER NOT NULL DEFAULT 0,
    status            VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by        UUID,
    updated_by        UUID,
    deleted_by        UUID,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at        TIMESTAMPTZ,
    version           INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_template_fields_template ON template_fields (asset_template_id);

CREATE TABLE assets (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id          UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    category_id        UUID REFERENCES categories(id) ON DELETE SET NULL,
    asset_template_id  UUID REFERENCES asset_templates(id) ON DELETE SET NULL,
    name               VARCHAR(255) NOT NULL,
    sku                VARCHAR(100),
    serial_number      VARCHAR(100),
    description        TEXT,
    purchase_price     NUMERIC(14,2),
    replacement_value  NUMERIC(14,2),
    purchase_date      DATE,
    condition          VARCHAR(30) NOT NULL DEFAULT 'good',
    location           VARCHAR(255),
    status             VARCHAR(30) NOT NULL DEFAULT 'available',
    created_by         UUID,
    updated_by         UUID,
    deleted_by         UUID,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at         TIMESTAMPTZ,
    version            INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_assets_tenant ON assets (tenant_id);
CREATE INDEX idx_assets_category ON assets (category_id);
CREATE INDEX idx_assets_status ON assets (tenant_id, status);

CREATE TABLE asset_values (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    asset_id          UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    template_field_id UUID NOT NULL REFERENCES template_fields(id) ON DELETE CASCADE,
    value             TEXT,
    status            VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at        TIMESTAMPTZ,
    version           INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_asset_values_field ON asset_values (asset_id, template_field_id) WHERE deleted_at IS NULL;

CREATE TABLE asset_images (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    asset_id    UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    url         VARCHAR(500) NOT NULL,
    alt_text    VARCHAR(255),
    is_primary  BOOLEAN NOT NULL DEFAULT false,
    sort_order  INTEGER NOT NULL DEFAULT 0,
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by  UUID,
    updated_by  UUID,
    deleted_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_asset_images_asset ON asset_images (asset_id);

CREATE TABLE asset_documents (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    asset_id    UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    file_url    VARCHAR(500) NOT NULL,
    file_type   VARCHAR(50),
    file_size   INTEGER,
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by  UUID,
    updated_by  UUID,
    deleted_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_asset_documents_asset ON asset_documents (asset_id);

CREATE TABLE asset_availability (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    asset_id          UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    start_date        TIMESTAMPTZ NOT NULL,
    end_date          TIMESTAMPTZ NOT NULL,
    availability_type VARCHAR(30) NOT NULL DEFAULT 'blocked',
    reason            TEXT,
    status            VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by        UUID,
    updated_by        UUID,
    deleted_by        UUID,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at        TIMESTAMPTZ,
    version           INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_asset_availability_asset_range ON asset_availability (asset_id, start_date, end_date);

-- =========================================================================
-- CRM
-- =========================================================================

CREATE TABLE customers (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    first_name        VARCHAR(100) NOT NULL,
    last_name         VARCHAR(100) NOT NULL,
    email             VARCHAR(255) NOT NULL,
    phone             VARCHAR(30),
    company_name      VARCHAR(255),
    date_of_birth     DATE,
    id_document_type  VARCHAR(50),
    id_document_number VARCHAR(100),
    customer_type     VARCHAR(30) NOT NULL DEFAULT 'individual',
    status            VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by        UUID,
    updated_by        UUID,
    deleted_by        UUID,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at        TIMESTAMPTZ,
    version           INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_customers_tenant ON customers (tenant_id);
CREATE INDEX idx_customers_email ON customers (tenant_id, email);

CREATE TABLE customer_addresses (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id  UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    address_type VARCHAR(30) NOT NULL DEFAULT 'billing',
    line1        VARCHAR(255) NOT NULL,
    line2        VARCHAR(255),
    city         VARCHAR(100) NOT NULL,
    state        VARCHAR(100),
    postal_code  VARCHAR(20),
    country      VARCHAR(100) NOT NULL,
    is_default   BOOLEAN NOT NULL DEFAULT false,
    status       VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by   UUID,
    updated_by   UUID,
    deleted_by   UUID,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ,
    version      INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_customer_addresses_customer ON customer_addresses (customer_id);

CREATE TABLE memberships (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id          UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id        UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    plan_name          VARCHAR(100) NOT NULL,
    tier               VARCHAR(50),
    start_date         DATE NOT NULL,
    end_date           DATE,
    fee                NUMERIC(12,2) NOT NULL DEFAULT 0,
    membership_status  VARCHAR(30) NOT NULL DEFAULT 'active',
    status             VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by         UUID,
    updated_by         UUID,
    deleted_by         UUID,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at         TIMESTAMPTZ,
    version            INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_memberships_customer ON memberships (customer_id);

CREATE TABLE loyalty_programs (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id            UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name                 VARCHAR(150) NOT NULL,
    description          TEXT,
    points_per_currency  NUMERIC(10,4) NOT NULL DEFAULT 1,
    redemption_rate      NUMERIC(10,4) NOT NULL DEFAULT 1,
    status               VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by           UUID,
    updated_by           UUID,
    deleted_by           UUID,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at           TIMESTAMPTZ,
    version              INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE loyalty_transactions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    loyalty_program_id  UUID NOT NULL REFERENCES loyalty_programs(id) ON DELETE CASCADE,
    customer_id         UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    booking_id          UUID,
    points              INTEGER NOT NULL,
    transaction_type    VARCHAR(30) NOT NULL,
    description         TEXT,
    status              VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at          TIMESTAMPTZ,
    version             INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_loyalty_transactions_customer ON loyalty_transactions (customer_id);

-- =========================================================================
-- PRICING
-- =========================================================================

CREATE TABLE pricing_rules (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    category_id   UUID REFERENCES categories(id) ON DELETE SET NULL,
    asset_id      UUID REFERENCES assets(id) ON DELETE SET NULL,
    name          VARCHAR(150) NOT NULL,
    rule_type     VARCHAR(30) NOT NULL,
    value         NUMERIC(14,2) NOT NULL,
    duration_unit VARCHAR(20) NOT NULL DEFAULT 'day',
    min_duration  INTEGER,
    max_duration  INTEGER,
    valid_from    DATE,
    valid_to      DATE,
    status        VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by    UUID,
    updated_by    UUID,
    deleted_by    UUID,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ,
    version       INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_pricing_rules_asset ON pricing_rules (asset_id);
CREATE INDEX idx_pricing_rules_category ON pricing_rules (category_id);

CREATE TABLE coupons (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    code            VARCHAR(50) NOT NULL,
    discount_type   VARCHAR(20) NOT NULL,
    discount_value  NUMERIC(14,2) NOT NULL,
    min_order_value NUMERIC(14,2) NOT NULL DEFAULT 0,
    usage_limit     INTEGER,
    used_count      INTEGER NOT NULL DEFAULT 0,
    valid_from      DATE,
    valid_to        DATE,
    status          VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by      UUID,
    updated_by      UUID,
    deleted_by      UUID,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    version         INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_coupons_tenant_code ON coupons (tenant_id, code) WHERE deleted_at IS NULL;

CREATE TABLE promotions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name            VARCHAR(150) NOT NULL,
    description     TEXT,
    promotion_type  VARCHAR(30) NOT NULL,
    value           NUMERIC(14,2) NOT NULL,
    start_date      DATE,
    end_date        DATE,
    status          VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by      UUID,
    updated_by      UUID,
    deleted_by      UUID,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    version         INTEGER NOT NULL DEFAULT 1
);

-- =========================================================================
-- BOOKING
-- =========================================================================

CREATE TABLE bookings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id     UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    coupon_id       UUID REFERENCES coupons(id) ON DELETE SET NULL,
    booking_number  VARCHAR(50) NOT NULL,
    start_date      TIMESTAMPTZ NOT NULL,
    end_date        TIMESTAMPTZ NOT NULL,
    subtotal        NUMERIC(14,2) NOT NULL DEFAULT 0,
    tax_total       NUMERIC(14,2) NOT NULL DEFAULT 0,
    discount_total  NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount    NUMERIC(14,2) NOT NULL DEFAULT 0,
    booking_status  VARCHAR(30) NOT NULL DEFAULT 'pending',
    payment_status  VARCHAR(30) NOT NULL DEFAULT 'unpaid',
    notes           TEXT,
    status          VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by      UUID,
    updated_by      UUID,
    deleted_by      UUID,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    version         INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_bookings_tenant_number ON bookings (tenant_id, booking_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_bookings_customer ON bookings (customer_id);
CREATE INDEX idx_bookings_tenant_status ON bookings (tenant_id, booking_status);

CREATE TABLE booking_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    booking_id  UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    asset_id    UUID NOT NULL REFERENCES assets(id) ON DELETE RESTRICT,
    quantity    INTEGER NOT NULL DEFAULT 1,
    unit_price  NUMERIC(14,2) NOT NULL,
    line_total  NUMERIC(14,2) NOT NULL,
    start_date  TIMESTAMPTZ NOT NULL,
    end_date    TIMESTAMPTZ NOT NULL,
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by  UUID,
    updated_by  UUID,
    deleted_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_booking_items_booking ON booking_items (booking_id);
CREATE INDEX idx_booking_items_asset_range ON booking_items (asset_id, start_date, end_date);

CREATE TABLE booking_extensions (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    booking_id       UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    booking_item_id  UUID NOT NULL REFERENCES booking_items(id) ON DELETE CASCADE,
    old_end_date     TIMESTAMPTZ NOT NULL,
    new_end_date     TIMESTAMPTZ NOT NULL,
    additional_cost  NUMERIC(14,2) NOT NULL DEFAULT 0,
    reason           TEXT,
    status           VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by       UUID,
    updated_by       UUID,
    deleted_by       UUID,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at       TIMESTAMPTZ,
    version          INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_booking_extensions_item ON booking_extensions (booking_item_id);

CREATE TABLE booking_returns (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    booking_id          UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    booking_item_id     UUID NOT NULL REFERENCES booking_items(id) ON DELETE CASCADE,
    returned_at         TIMESTAMPTZ NOT NULL,
    condition_on_return VARCHAR(30) NOT NULL DEFAULT 'good',
    late_fee            NUMERIC(14,2) NOT NULL DEFAULT 0,
    damage_fee          NUMERIC(14,2) NOT NULL DEFAULT 0,
    notes               TEXT,
    status              VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by          UUID,
    updated_by          UUID,
    deleted_by          UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at          TIMESTAMPTZ,
    version             INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_booking_returns_item ON booking_returns (booking_item_id);

-- =========================================================================
-- FINANCE
-- =========================================================================

CREATE TABLE taxes (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    rate        NUMERIC(6,4) NOT NULL,
    tax_type    VARCHAR(30) NOT NULL DEFAULT 'percentage',
    is_default  BOOLEAN NOT NULL DEFAULT false,
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by  UUID,
    updated_by  UUID,
    deleted_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE invoices (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id     UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    booking_id      UUID REFERENCES bookings(id) ON DELETE SET NULL,
    invoice_number  VARCHAR(50) NOT NULL,
    issue_date      DATE NOT NULL,
    due_date        DATE NOT NULL,
    subtotal        NUMERIC(14,2) NOT NULL DEFAULT 0,
    tax_total       NUMERIC(14,2) NOT NULL DEFAULT 0,
    discount_total  NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_amount    NUMERIC(14,2) NOT NULL DEFAULT 0,
    amount_paid     NUMERIC(14,2) NOT NULL DEFAULT 0,
    amount_due      NUMERIC(14,2) NOT NULL DEFAULT 0,
    invoice_status  VARCHAR(30) NOT NULL DEFAULT 'unpaid',
    status          VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by      UUID,
    updated_by      UUID,
    deleted_by      UUID,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    version         INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_invoices_tenant_number ON invoices (tenant_id, invoice_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_invoices_customer ON invoices (customer_id);
CREATE INDEX idx_invoices_booking ON invoices (booking_id);

CREATE TABLE invoice_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    invoice_id      UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    booking_item_id UUID REFERENCES booking_items(id) ON DELETE SET NULL,
    description     VARCHAR(255) NOT NULL,
    quantity        INTEGER NOT NULL DEFAULT 1,
    unit_price      NUMERIC(14,2) NOT NULL,
    tax_amount      NUMERIC(14,2) NOT NULL DEFAULT 0,
    line_total      NUMERIC(14,2) NOT NULL,
    status          VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    version         INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_invoice_items_invoice ON invoice_items (invoice_id);

CREATE TABLE payments (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    invoice_id             UUID NOT NULL REFERENCES invoices(id) ON DELETE RESTRICT,
    customer_id            UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    payment_method         VARCHAR(50) NOT NULL,
    transaction_reference  VARCHAR(150),
    amount                 NUMERIC(14,2) NOT NULL,
    currency               VARCHAR(10) NOT NULL DEFAULT 'USD',
    paid_at                TIMESTAMPTZ,
    payment_status         VARCHAR(30) NOT NULL DEFAULT 'pending',
    status                 VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by             UUID,
    updated_by             UUID,
    deleted_by             UUID,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at             TIMESTAMPTZ,
    version                INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_payments_invoice ON payments (invoice_id);

CREATE TABLE refunds (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    payment_id      UUID NOT NULL REFERENCES payments(id) ON DELETE RESTRICT,
    amount          NUMERIC(14,2) NOT NULL,
    reason          TEXT,
    refund_status   VARCHAR(30) NOT NULL DEFAULT 'pending',
    processed_at    TIMESTAMPTZ,
    status          VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by      UUID,
    updated_by      UUID,
    deleted_by      UUID,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    version         INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_refunds_payment ON refunds (payment_id);

CREATE TABLE expenses (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    asset_id     UUID REFERENCES assets(id) ON DELETE SET NULL,
    category     VARCHAR(100) NOT NULL,
    amount       NUMERIC(14,2) NOT NULL,
    expense_date DATE NOT NULL,
    description  TEXT,
    vendor       VARCHAR(255),
    status       VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by   UUID,
    updated_by   UUID,
    deleted_by   UUID,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ,
    version      INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_expenses_asset ON expenses (asset_id);

-- =========================================================================
-- MAINTENANCE
-- =========================================================================

CREATE TABLE maintenance_records (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    asset_id            UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    maintenance_type    VARCHAR(50) NOT NULL,
    description         TEXT,
    cost                NUMERIC(14,2) NOT NULL DEFAULT 0,
    scheduled_date      DATE,
    completed_date      DATE,
    performed_by        VARCHAR(150),
    maintenance_status  VARCHAR(30) NOT NULL DEFAULT 'scheduled',
    status              VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by          UUID,
    updated_by          UUID,
    deleted_by          UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at          TIMESTAMPTZ,
    version             INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_maintenance_records_asset ON maintenance_records (asset_id);
CREATE INDEX idx_maintenance_records_due ON maintenance_records (scheduled_date) WHERE maintenance_status = 'scheduled';

CREATE TABLE inspections (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    asset_id         UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    booking_item_id  UUID REFERENCES booking_items(id) ON DELETE SET NULL,
    inspection_type  VARCHAR(50) NOT NULL,
    inspected_at     TIMESTAMPTZ NOT NULL,
    inspector_name   VARCHAR(150),
    findings         TEXT,
    result           VARCHAR(30) NOT NULL DEFAULT 'pass',
    status           VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by       UUID,
    updated_by       UUID,
    deleted_by       UUID,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at       TIMESTAMPTZ,
    version          INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_inspections_asset ON inspections (asset_id);

CREATE TABLE damage_reports (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    asset_id        UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    booking_id      UUID REFERENCES bookings(id) ON DELETE SET NULL,
    inspection_id   UUID REFERENCES inspections(id) ON DELETE SET NULL,
    description     TEXT NOT NULL,
    severity        VARCHAR(30) NOT NULL DEFAULT 'minor',
    repair_cost     NUMERIC(14,2) NOT NULL DEFAULT 0,
    charged_amount  NUMERIC(14,2) NOT NULL DEFAULT 0,
    report_status   VARCHAR(30) NOT NULL DEFAULT 'open',
    status          VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by      UUID,
    updated_by      UUID,
    deleted_by      UUID,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    version         INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_damage_reports_asset ON damage_reports (asset_id);

-- =========================================================================
-- NOTIFICATIONS
-- =========================================================================

CREATE TABLE notification_templates (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name           VARCHAR(150) NOT NULL,
    channel        VARCHAR(30) NOT NULL,
    subject        VARCHAR(255),
    body           TEXT NOT NULL,
    event_trigger  VARCHAR(100) NOT NULL,
    status         VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by     UUID,
    updated_by     UUID,
    deleted_by     UUID,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at     TIMESTAMPTZ,
    version        INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_notification_templates_trigger ON notification_templates (tenant_id, event_trigger);

CREATE TABLE notifications (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id      UUID REFERENCES users(id) ON DELETE CASCADE,
    customer_id  UUID REFERENCES customers(id) ON DELETE CASCADE,
    channel      VARCHAR(30) NOT NULL DEFAULT 'in_app',
    title        VARCHAR(255) NOT NULL,
    message      TEXT NOT NULL,
    is_read      BOOLEAN NOT NULL DEFAULT false,
    read_at      TIMESTAMPTZ,
    status       VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ,
    version      INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_notifications_user ON notifications (user_id, is_read);

CREATE TABLE notification_logs (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    notification_template_id UUID REFERENCES notification_templates(id) ON DELETE SET NULL,
    recipient_id             UUID NOT NULL,
    recipient_type           VARCHAR(30) NOT NULL,
    channel                  VARCHAR(30) NOT NULL,
    delivery_status          VARCHAR(30) NOT NULL DEFAULT 'pending',
    error_message            TEXT,
    sent_at                  TIMESTAMPTZ,
    status                   VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at               TIMESTAMPTZ,
    version                  INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_notification_logs_recipient ON notification_logs (recipient_id, recipient_type);

-- =========================================================================
-- CMS
-- =========================================================================

CREATE TABLE websites (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    domain      VARCHAR(255) NOT NULL,
    title       VARCHAR(255) NOT NULL,
    theme       VARCHAR(100),
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by  UUID,
    updated_by  UUID,
    deleted_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_websites_domain ON websites (domain) WHERE deleted_at IS NULL;

CREATE TABLE pages (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    website_id  UUID NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL,
    content     TEXT,
    template    VARCHAR(100),
    status      VARCHAR(30) NOT NULL DEFAULT 'draft',
    created_by  UUID,
    updated_by  UUID,
    deleted_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_pages_website_slug ON pages (website_id, slug) WHERE deleted_at IS NULL;

CREATE TABLE menus (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    website_id  UUID NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    name        VARCHAR(150) NOT NULL,
    location    VARCHAR(50),
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by  UUID,
    updated_by  UUID,
    deleted_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE menu_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    menu_id     UUID NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
    parent_id   UUID REFERENCES menu_items(id) ON DELETE SET NULL,
    label       VARCHAR(150) NOT NULL,
    url         VARCHAR(500),
    sort_order  INTEGER NOT NULL DEFAULT 0,
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by  UUID,
    updated_by  UUID,
    deleted_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_menu_items_menu ON menu_items (menu_id);
CREATE INDEX idx_menu_items_parent ON menu_items (parent_id);

CREATE TABLE blog_categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name        VARCHAR(150) NOT NULL,
    slug        VARCHAR(150) NOT NULL,
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by  UUID,
    updated_by  UUID,
    deleted_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_blog_categories_tenant_slug ON blog_categories (tenant_id, slug) WHERE deleted_at IS NULL;

CREATE TABLE blogs (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    website_id        UUID NOT NULL REFERENCES websites(id) ON DELETE CASCADE,
    author_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    blog_category_id  UUID REFERENCES blog_categories(id) ON DELETE SET NULL,
    title             VARCHAR(255) NOT NULL,
    slug              VARCHAR(255) NOT NULL,
    content           TEXT,
    featured_image    VARCHAR(500),
    published_at      TIMESTAMPTZ,
    status            VARCHAR(30) NOT NULL DEFAULT 'draft',
    created_by        UUID,
    updated_by        UUID,
    deleted_by        UUID,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at        TIMESTAMPTZ,
    version           INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_blogs_website_slug ON blogs (website_id, slug) WHERE deleted_at IS NULL;

CREATE TABLE seo_meta (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    entity_type       VARCHAR(50) NOT NULL,
    entity_id         UUID NOT NULL,
    meta_title        VARCHAR(255),
    meta_description  TEXT,
    meta_keywords     VARCHAR(500),
    og_image          VARCHAR(500),
    canonical_url     VARCHAR(500),
    status            VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at        TIMESTAMPTZ,
    version           INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_seo_meta_entity ON seo_meta (entity_type, entity_id) WHERE deleted_at IS NULL;

-- =========================================================================
-- REPORTS / ANALYTICS
-- =========================================================================

CREATE TABLE reports (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name              VARCHAR(255) NOT NULL,
    report_type       VARCHAR(50) NOT NULL,
    parameters        JSONB,
    generated_format  VARCHAR(20) NOT NULL DEFAULT 'pdf',
    file_url          VARCHAR(500),
    status            VARCHAR(30) NOT NULL DEFAULT 'queued',
    created_by        UUID,
    updated_by        UUID,
    deleted_by        UUID,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at        TIMESTAMPTZ,
    version           INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_reports_tenant ON reports (tenant_id);

CREATE TABLE analytics_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    customer_id     UUID REFERENCES customers(id) ON DELETE SET NULL,
    event_name      VARCHAR(100) NOT NULL,
    event_category  VARCHAR(100),
    event_data      JSONB,
    source          VARCHAR(50),
    occurred_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    status          VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    version         INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_analytics_events_tenant_time ON analytics_events (tenant_id, occurred_at DESC);

-- =========================================================================
-- API KEYS / WEBHOOKS
-- =========================================================================

CREATE TABLE api_keys (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name        VARCHAR(150) NOT NULL,
    key_prefix  VARCHAR(20) NOT NULL,
    key_hash    VARCHAR(255) NOT NULL,
    scopes      JSONB,
    last_used_at TIMESTAMPTZ,
    expires_at  TIMESTAMPTZ,
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by  UUID,
    updated_by  UUID,
    deleted_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE UNIQUE INDEX uq_api_keys_hash ON api_keys (key_hash) WHERE deleted_at IS NULL;

CREATE TABLE webhooks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    url         VARCHAR(500) NOT NULL,
    events      JSONB NOT NULL,
    secret      VARCHAR(255) NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    status      VARCHAR(30) NOT NULL DEFAULT 'active',
    created_by  UUID,
    updated_by  UUID,
    deleted_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    version     INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_webhooks_tenant ON webhooks (tenant_id) WHERE is_active = true;

CREATE TABLE webhook_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    webhook_id      UUID NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
    event_type      VARCHAR(100) NOT NULL,
    payload         JSONB,
    response_code   INTEGER,
    response_body   TEXT,
    triggered_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    status          VARCHAR(30) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    version         INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_webhook_logs_webhook ON webhook_logs (webhook_id);

-- =========================================================================
-- updated_at AUTO-MAINTENANCE
-- A single trigger function applied to every table that has an
-- updated_at column, attached generically via a DO block so adding a new
-- table never requires a new trigger definition by hand.
-- =========================================================================

CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
DECLARE
    t text;
BEGIN
    FOR t IN
        SELECT table_name FROM information_schema.columns
        WHERE table_schema = 'public' AND column_name = 'updated_at'
    LOOP
        EXECUTE format(
            'CREATE TRIGGER trg_set_updated_at BEFORE UPDATE ON %I FOR EACH ROW EXECUTE FUNCTION set_updated_at();',
            t
        );
    END LOOP;
END $$;
