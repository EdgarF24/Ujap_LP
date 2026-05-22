-- ============================================================
-- Migration 001: Create spaces and invoices tables
-- ============================================================

-- Enable pgcrypto extension for UUID generation (if not already enabled)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ── Space Type Enum ──────────────────────────────────────────
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'space_type') THEN
        CREATE TYPE space_type AS ENUM (
            'DESK',
            'OFFICE',
            'MEETING_ROOM',
            'LOUNGE'
        );
    END IF;
END$$;

-- ── Invoice Status Enum ──────────────────────────────────────
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'invoice_status') THEN
        CREATE TYPE invoice_status AS ENUM (
            'PENDING',
            'PAID',
            'CANCELLED'
        );
    END IF;
END$$;

-- ── Spaces Table ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS spaces (
    id             UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    name           VARCHAR(255)    NOT NULL,
    type           VARCHAR(50)     NOT NULL
                       CHECK (type IN ('DESK', 'OFFICE', 'MEETING_ROOM', 'LOUNGE')),
    capacity       INTEGER         NOT NULL CHECK (capacity > 0),
    price_per_hour NUMERIC(10, 2)  NOT NULL CHECK (price_per_hour > 0),
    amenities      TEXT            NOT NULL DEFAULT '',
    is_active      BOOLEAN         NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- Indexes for spaces
CREATE INDEX IF NOT EXISTS idx_spaces_type      ON spaces (type);
CREATE INDEX IF NOT EXISTS idx_spaces_is_active ON spaces (is_active);
CREATE INDEX IF NOT EXISTS idx_spaces_capacity  ON spaces (capacity);

-- ── Invoices Table ───────────────────────────────────────────
CREATE TABLE IF NOT EXISTS invoices (
    id             UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        VARCHAR(255)    NOT NULL,
    reservation_id VARCHAR(255)    NOT NULL,
    amount         NUMERIC(10, 2)  NOT NULL CHECK (amount > 0),
    status         VARCHAR(50)     NOT NULL DEFAULT 'PENDING'
                       CHECK (status IN ('PENDING', 'PAID', 'CANCELLED')),
    issued_at      TIMESTAMPTZ     NOT NULL,
    due_at         TIMESTAMPTZ     NOT NULL,
    paid_at        TIMESTAMPTZ     NULL,
    created_at     TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- Indexes for invoices
CREATE INDEX IF NOT EXISTS idx_invoices_user_id        ON invoices (user_id);
CREATE INDEX IF NOT EXISTS idx_invoices_reservation_id ON invoices (reservation_id);
CREATE INDEX IF NOT EXISTS idx_invoices_status         ON invoices (status);
CREATE INDEX IF NOT EXISTS idx_invoices_issued_at      ON invoices (issued_at);

-- ── Trigger: auto-update updated_at on spaces ────────────────
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = 'set_spaces_updated_at'
    ) THEN
        CREATE TRIGGER set_spaces_updated_at
            BEFORE UPDATE ON spaces
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = 'set_invoices_updated_at'
    ) THEN
        CREATE TRIGGER set_invoices_updated_at
            BEFORE UPDATE ON invoices
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
    END IF;
END$$;
