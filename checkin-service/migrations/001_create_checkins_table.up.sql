CREATE SCHEMA IF NOT EXISTS checkin;
SET search_path TO checkin;

-- Checkin records: one row per successful (or attempted) ticket scan on event day.
-- ticket_id has a UNIQUE constraint — idempotent check-ins at the DB level.

CREATE TABLE checkins (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id     UUID         NOT NULL UNIQUE,
    event_id      UUID         NOT NULL,
    user_id       UUID         NOT NULL,
    staff_id      UUID,
    qr_code       VARCHAR(512) NOT NULL,
    status        VARCHAR(50)  NOT NULL DEFAULT 'success',
    check_in_time TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    device_id     VARCHAR(255),
    gate          VARCHAR(100),
    notes         TEXT,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE checkins IS 'Records each ticket scan on event day';
COMMENT ON COLUMN checkins.status IS 'success | invalid | already_used | cancelled';

CREATE INDEX idx_checkins_event_id   ON checkins (event_id);
CREATE INDEX idx_checkins_user_id    ON checkins (user_id);
CREATE INDEX idx_checkins_event_time ON checkins (event_id, check_in_time DESC);
CREATE INDEX idx_checkins_gate       ON checkins (event_id, gate);
