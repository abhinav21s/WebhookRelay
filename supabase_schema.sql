-- WebhookRelay Database Schema for Supabase
-- Run this in Supabase SQL Editor

-- Drop tables if they exist (for clean setup)
DROP TABLE IF EXISTS delivery_attempts CASCADE;
DROP TABLE IF EXISTS events CASCADE;
DROP TABLE IF EXISTS endpoints CASCADE;

-- Endpoints table
CREATE TABLE endpoints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    secret TEXT NOT NULL,
    event_types TEXT[] NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Events table
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Delivery attempts table
CREATE TABLE delivery_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    endpoint_id UUID NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
    attempt_number INT NOT NULL,
    status_code INT,
    latency_ms INT,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_endpoints_active ON endpoints(is_active);
CREATE INDEX idx_endpoints_created_at ON endpoints(created_at DESC);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_event_type ON events(event_type);
CREATE INDEX idx_events_created_at ON events(created_at DESC);
CREATE INDEX idx_delivery_attempts_event_id ON delivery_attempts(event_id);
CREATE INDEX idx_delivery_attempts_endpoint_id ON delivery_attempts(endpoint_id);
CREATE INDEX idx_delivery_attempts_created_at ON delivery_attempts(created_at DESC);

-- Optional: Add sample data for testing
-- INSERT INTO endpoints (name, url, secret, event_types, is_active)
-- VALUES 
--   ('Test Endpoint', 'http://localhost:9090/webhook', 'test-secret-123', ARRAY['test.event', 'payment.completed'], true),
--   ('Demo Endpoint', 'http://localhost:9090/webhook-a', 'demo-secret-456', ARRAY['user.created'], true);

-- Verify tables created
SELECT 
    'endpoints' as table_name, 
    COUNT(*) as row_count 
FROM endpoints
UNION ALL
SELECT 
    'events' as table_name, 
    COUNT(*) as row_count 
FROM events
UNION ALL
SELECT 
    'delivery_attempts' as table_name, 
    COUNT(*) as row_count 
FROM delivery_attempts;
