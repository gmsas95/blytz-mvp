-- Seed data for payment service testing and development
-- This file contains sample data for testing the payment system

-- Insert sample payment methods with comprehensive configuration
INSERT INTO payment_methods (code, name, type, description, icon_url, enabled, supported_currencies, min_amount, max_amount, processing_time_seconds, requires_3ds, requires_redirect, country_codes, metadata) VALUES
-- Malaysian payment methods
('FPX', 'FPX Online Banking', 'bank_transfer', 'Pay directly from your bank account via FPX', 'https://example.com/icons/fpx.png', true, ARRAY['MYR'], 1.00, 30000.00, 5, false, false, ARRAY['MY'], '{"banks": ["Maybank", "CIMB", "Public Bank", "RHB", "Hong Leong"], "available_hours": "06:00-23:00"}'),
('FPX_BB', 'FPX Business Banking', 'bank_transfer', 'Business banking transfers for high-value transactions', 'https://example.com/icons/fpx-bb.png', true, ARRAY['MYR'], 1000.00, 1000000.00, 10, false, false, ARRAY['MY'], '{"transaction_type": "B2B", "requires_business_registration": true}'),
('GRABPAY', 'GrabPay', 'ewallet', 'Pay using your GrabPay e-wallet balance', 'https://example.com/icons/grabpay.png', true, ARRAY['MYR'], 1.00, 5000.00, 3, false, true, ARRAY['MY'], '{"app_required": true, "instant_confirmation": true}'),
('TNG', 'Touch ''n Go eWallet', 'ewallet', 'Pay using Touch ''n Go e-wallet', 'https://example.com/icons/tng.png', true, ARRAY['MYR'], 1.00, 10000.00, 3, false, true, ARRAY['MY'], '{"app_required": true, "qr_support": true}'),
('SHOPEEPAY', 'ShopeePay', 'ewallet', 'Pay using ShopeePay e-wallet', 'https://example.com/icons/shopeepay.png', true, ARRAY['MYR'], 1.00, 5000.00, 3, false, true, ARRAY['MY'], '{"integration_type": "seamless", "cashback_available": true}'),
('BOOST', 'Boost', 'ewallet', 'Pay using Boost e-wallet', 'https://example.com/icons/boost.png', true, ARRAY['MYR'], 1.00, 3000.00, 3, false, true, ARRAY['MY'], '{"app_required": true, " loyalty_points": true}'),
('MAYBANKQR', 'MAE QR', 'qr_code', 'Scan QR code with Maybank MAE app', 'https://example.com/icons/mae-qr.png', true, ARRAY['MYR'], 0.10, 5000.00, 2, false, false, ARRAY['MY'], '{"bank": "Maybank", "qr_type": "static"}'),
('DUITNOWQR', 'DuitNow QR', 'qr_code', 'Universal QR code for any Malaysian bank app', 'https://example.com/icons/duitnow-qr.png', true, ARRAY['MYR'], 0.10, 10000.00, 2, false, false, ARRAY['MY'], '{"network": "DuitNow", "bank_support": 40}'),
('ATOME', 'Atome', 'bnpl', 'Buy now, pay later in 3 installments', 'https://example.com/icons/atome.png', true, ARRAY['MYR'], 10.00, 1000.00, 15, false, true, ARRAY['MY'], '{"installments": 3, "interest_free": true, "credit_check": true}'),
('RELY', 'Rely', 'bnpl', 'Buy now, pay later with flexible terms', 'https://example.com/icons/rely.png', true, ARRAY['MYR'], 50.00, 5000.00, 20, false, true, ARRAY['MY'], '{"installments": "3-12", "interest_rate": "0-8%"}'),

-- Credit cards
('VISA', 'Visa Credit/Debit', 'credit_card', 'Pay with Visa credit or debit cards', 'https://example.com/icons/visa.png', true, ARRAY['MYR', 'SGD', 'USD'], 1.00, 50000.00, 30, true, true, ARRAY['MY', 'SG', 'US'], '{"card_types": ["Credit", "Debit", "Prepaid"], "3ds_version": "2.0"}'),
('MASTERCARD', 'Mastercard Credit/Debit', 'credit_card', 'Pay with Mastercard credit or debit cards', 'https://example.com/icons/mastercard.png', true, ARRAY['MYR', 'SGD', 'USD'], 1.00, 50000.00, 30, true, true, ARRAY['MY', 'SG', 'US'], '{"card_types": ["Credit", "Debit", "Prepaid"], "3ds_version": "2.0"}'),

-- Singapore payment methods (for future expansion)
('PAYNOW', 'PayNow', 'qr_code', 'Singapore''s QR payment system', 'https://example.com/icons/paynow.png', false, ARRAY['SGD'], 0.01, 10000.00, 2, false, false, ARRAY['SG'], '{"country": "Singapore", "network": "PayNow"}'),
('ENETS', 'eNETS', 'bank_transfer', 'Singapore electronic funds transfer', 'https://example.com/icons/enets.png', false, ARRAY['SGD'], 1.00, 10000.00, 5, false, false, ARRAY['SG'], '{"country": "Singapore", "banks": ["DBS", "OCBC", "UOB"]}'),

-- Philippines payment methods (for future expansion)
('GCASH', 'GCash', 'ewallet', 'Philippines most popular e-wallet', 'https://example.com/icons/gcash.png', false, ARRAY['PHP'], 1.00, 50000.00, 3, false, true, ARRAY['PH'], '{"country": "Philippines", "qr_support": true}'),
('MAYA', 'Maya', 'ewallet', 'Pay with Maya e-wallet', 'https://example.com/icons/maya.png', false, ARRAY['PHP'], 1.00, 30000.00, 3, false, true, ARRAY['PH'], '{"country": "Philippines", "former_name": "PayMaya"}')

ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    description = EXCLUDED.description,
    icon_url = EXCLUDED.icon_url,
    enabled = EXCLUDED.enabled,
    supported_currencies = EXCLUDED.supported_currencies,
    min_amount = EXCLUDED.min_amount,
    max_amount = EXCLUDED.max_amount,
    processing_time_seconds = EXCLUDED.processing_time_seconds,
    requires_3ds = EXCLUDED.requires_3ds,
    requires_redirect = EXCLUDED.requires_redirect,
    country_codes = EXCLUDED.country_codes,
    metadata = EXCLUDED.metadata;

-- Insert sample payment records for testing
INSERT INTO payments (
    id, user_id, order_id, transaction_id, amount, currency, payment_method, payment_channel,
    status, gateway_status, bill_name, bill_email, bill_mobile, bill_description,
    return_url, notify_url, fiuu_merchant_id, fiuu_reference_no, fiuu_approval_code,
    ip_address, user_agent, metadata
) VALUES
-- Successful FPX payment
('550e8400-e29b-41d4-a716-446655440001', 'user123', 'ORD202401001', 'TXN123456789', 150.75, 'MYR', 'FPX', 'Maybank',
 'completed', 'SUCCESS', 'John Doe', 'john.doe@example.com', '01234567890', 'Auction item payment - Vintage Camera',
 'https://example.com/return', 'https://example.com/webhook', 'MERCHANT123', 'REF123456789', 'APP987654',
 '203.106.85.123', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36', '{"auction_id": "AUCT123", "item_name": "Vintage Camera"}'),

-- Pending GrabPay payment
('550e8400-e29b-41d4-a716-446655440002', 'user456', 'ORD202401002', NULL, 89.99, 'MYR', 'GRABPAY', NULL,
 'pending', 'PENDING', 'Jane Smith', 'jane.smith@example.com', '01123456789', 'Online shopping - Electronics',
 'https://example.com/return', 'https://example.com/webhook', 'MERCHANT123', NULL, NULL,
 '210.195.38.45', 'Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X)', '{"store_id": "STORE456", "category": "Electronics"}'),

-- Failed credit card payment
('550e8400-e29b-41d4-a716-446655440003', 'user789', 'ORD202401003', 'TXN987654321', 299.00, 'MYR', 'VISA', 'Credit Card',
 'failed', 'DECLINED', 'Bob Johnson', 'bob.johnson@example.com', '01623456789', 'Flight tickets',
 'https://example.com/return', 'https://example.com/webhook', 'MERCHANT123', 'REF987654321', NULL,
 '115.132.74.15', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36', '{"booking_ref": "BK789", "travel_date": "2024-02-15"}', 'INSUFFICIENT_FUNDS'),

-- Refunded payment
('550e8400-e29b-41d4-a716-446655440004', 'user234', 'ORD202401004', 'TXN456789123', 199.50, 'MYR', 'TNG', 'Touch n Go',
 'refunded', 'SUCCESS', 'Alice Brown', 'alice.brown@example.com', '01987654321', 'Fashion items',
 'https://example.com/return', 'https://example.com/webhook', 'MERCHANT123', 'REF456789123', 'APP123456789',
 '60.50.12.34', 'Mozilla/5.0 (Android 11; Mobile; rv:68.0) Gecko/68.0 Firefox/88.0', '{"order_items": 3, "store": "Fashion Store"}'),

-- Expired payment
('550e8400-e29b-41d4-a716-446655440005', 'user567', 'ORD202401005', NULL, 75.00, 'MYR', 'FPX', 'Public Bank',
 'expired', 'TIMEOUT', 'Charlie Wilson', 'charlie.wilson@example.com', '01456789012', 'Food delivery',
 'https://example.com/return', 'https://example.com/webhook', 'MERCHANT123', NULL, NULL,
 '175.136.96.78', 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36', '{"restaurant": "Pizza Place", "delivery_time": "30 mins"}')

ON CONFLICT (id) DO NOTHING;

-- Update timestamps for realistic data
UPDATE payments SET
    created_at = CASE id
        WHEN '550e8400-e29b-41d4-a716-446655440001' THEN '2024-01-15 10:30:00+08:00'
        WHEN '550e8400-e29b-41d4-a716-446655440002' THEN '2024-01-15 14:45:00+08:00'
        WHEN '550e8400-e29b-41d4-a716-446655440003' THEN '2024-01-15 16:20:00+08:00'
        WHEN '550e8400-e29b-41d4-a716-446655440004' THEN '2024-01-14 09:15:00+08:00'
        WHEN '550e8400-e29b-41d4-a716-446655440005' THEN '2024-01-13 20:10:00+08:00'
    END,
    updated_at = CASE id
        WHEN '550e8400-e29b-41d4-a716-446655440001' THEN '2024-01-15 10:35:00+08:00'
        WHEN '550e8400-e29b-41d4-a716-446655440002' THEN '2024-01-15 14:45:00+08:00'
        WHEN '550e8400-e29b-41d4-a716-446655440003' THEN '2024-01-15 16:22:00+08:00'
        WHEN '550e8400-e29b-41d4-a716-446655440004' THEN '2024-01-16 11:30:00+08:00'
        WHEN '550e8400-e29b-41d4-a716-446655440005' THEN '2024-01-13 22:10:00+08:00'
    END,
    completed_at = CASE id
        WHEN '550e8400-e29b-41d4-a716-446655440001' THEN '2024-01-15 10:35:00+08:00'
        WHEN '550e8400-e29b-41d4-a716-446655440004' THEN '2024-01-14 09:20:00+08:00'
        WHEN '550e8400-e29b-41d4-a716-446655440005' THEN '2024-01-13 22:10:00+08:00'
    END,
    expired_at = CASE id
        WHEN '550e8400-e29b-41d4-a716-446655440005' THEN '2024-01-13 22:10:00+08:00'
    END;

-- Insert sample refunds
INSERT INTO refunds (
    id, payment_id, refund_id, amount, currency, reason, status, gateway_status,
    fiuu_refund_id, processed_by, metadata
) VALUES
('660e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440004', 'REF789123456', 199.50, 'MYR', 'Customer requested refund - wrong item', 'completed', 'SUCCESS',
 'RF789123456', 'admin1', '{"reason_code": "WRONG_ITEM", "customer_contacted": true}'),

('660e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440001', 'REF456789123', 50.00, 'MYR', 'Partial refund - item damaged', 'processing', 'PENDING',
 NULL, 'admin2', '{"reason_code": "DAMAGED_ITEM", "partial_refund": true}')

ON CONFLICT (id) DO NOTHING;

-- Update refund timestamps
UPDATE refunds SET
    created_at = CASE id
        WHEN '660e8400-e29b-41d4-a716-446655440001' THEN '2024-01-16 11:00:00+08:00'
        WHEN '660e8400-e29b-41d4-a716-446655440002' THEN '2024-01-16 15:30:00+08:00'
    END,
    updated_at = CASE id
        WHEN '660e8400-e29b-41d4-a716-446655440001' THEN '2024-01-16 11:35:00+08:00'
        WHEN '660e8400-e29b-41d4-a716-446655440002' THEN '2024-01-16 15:30:00+08:00'
    END,
    completed_at = CASE id
        WHEN '660e8400-e29b-41d4-a716-446655440001' THEN '2024-01-16 11:35:00+08:00'
    END;

-- Insert sample payment attempts
INSERT INTO payment_attempts (
    id, payment_id, attempt_number, status, gateway_request, gateway_response,
    started_at, completed_at, duration_ms, error_code, error_message
) VALUES
-- Successful FPX payment attempts
('770e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 1, 'completed',
 '{"merchant_id": "MERCHANT123", "amount": "150.75", "channel": "FPX"}',
 '{"transaction_id": "TXN123456789", "status": "SUCCESS"}',
 '2024-01-15 10:30:00+08:00', '2024-01-15 10:32:45+08:00', 165000, NULL, NULL),

-- Failed credit card payment attempts
('770e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440003', 1, 'failed',
 '{"merchant_id": "MERCHANT123", "amount": "299.00", "channel": "VISA"}',
 '{"error_code": "DECLINED", "error_desc": "Insufficient funds"}',
 '2024-01-15 16:20:00+08:00', '2024-01-15 16:21:30+08:00', 90000, 'DECLINED', 'Insufficient funds'),

-- Pending GrabPay payment attempt
('770e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440002', 1, 'processing',
 '{"merchant_id": "MERCHANT123", "amount": "89.99", "channel": "GRABPAY"}',
 NULL,
 '2024-01-15 14:45:00+08:00', NULL, NULL, NULL, NULL)

ON CONFLICT (id) DO NOTHING;

-- Insert sample webhook events
INSERT INTO webhooks (
    id, event_type, source, source_id, payload, signature, status,
    created_at, processed_at, retry_count
) VALUES
-- Payment success webhook
('880e8400-e29b-41d4-a716-446655440001', 'payment.success', 'fiuu', 'TXN123456789',
 '{"transaction_id": "TXN123456789", "order_id": "ORD202401001", "amount": "150.75", "status": "1", "signature": "abc123"}',
 'sha256=def456', 'completed',
 '2024-01-15 10:35:00+08:00', '2024-01-15 10:35:05+08:00', 0),

-- Payment failure webhook
('880e8400-e29b-41d4-a716-446655440002', 'payment.failed', 'fiuu', 'TXN987654321',
 '{"transaction_id": "TXN987654321", "order_id": "ORD202401003", "amount": "299.00", "status": "0", "error": "Insufficient funds", "signature": "ghi789"}',
 'sha256=jkl012', 'completed',
 '2024-01-15 16:22:00+08:00', '2024-01-15 16:22:03+08:00', 0),

-- Pending webhook for retry
('880e8400-e29b-41d4-a716-446655440003', 'payment.success', 'fiuu', 'TXN456789123',
 '{"transaction_id": "TXN456789123", "order_id": "ORD202401004", "amount": "199.50", "status": "1", "signature": "mno345"}',
 'sha256=pqr678', 'retrying',
 '2024-01-16 11:35:00+08:00', NULL, 2)

ON CONFLICT (id) DO NOTHING;

-- Insert sample settlements
INSERT INTO payment_settlements (
    id, payment_id, settlement_id, amount, currency, fee_amount, net_amount,
    settlement_date, batch_id, status, settled_at
) VALUES
-- Settled FPX payment
('990e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 'SETT20240115001', 150.75, 'MYR', 1.50, 149.25,
 '2024-01-15', 'BATCH20240115001', 'completed', '2024-01-16 10:00:00+08:00'),

-- Pending settlement
('990e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440004', 'SETT20240116001', 199.50, 'MYR', 2.00, 197.50,
 '2024-01-16', 'BATCH20240116001', 'pending', NULL)

ON CONFLICT (id) DO NOTHING;

-- Create database user for payment service (for development)
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'payment_service') THEN
        CREATE ROLE payment_service WITH LOGIN PASSWORD 'payment_service_password';
    END IF;
END
$$;

-- Grant permissions to payment service user
GRANT CONNECT ON DATABASE postgres TO payment_service;
GRANT USAGE ON SCHEMA public TO payment_service;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO payment_service;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO payment_service;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO payment_service;

-- Create read-only user for analytics
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'payment_analytics') THEN
        CREATE ROLE payment_analytics WITH LOGIN PASSWORD 'analytics_password';
    END IF;
END
$$;

-- Grant read-only permissions
GRANT CONNECT ON DATABASE postgres TO payment_analytics;
GRANT USAGE ON SCHEMA public TO payment_analytics;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO payment_analytics;

-- Set default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO payment_service;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO payment_analytics;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO payment_service;

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_payments_composite_lookup ON payments(user_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_payments_gateway_composite ON payments(gateway_status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_payments_settlement_lookup ON payments(status, created_at) WHERE status IN ('completed', 'refunded');
CREATE INDEX IF NOT EXISTS idx_webhooks_retry_processing ON webhooks(next_retry_at, status) WHERE status IN ('retrying', 'pending');

-- Create partial indexes for common queries
CREATE INDEX IF NOT EXISTS idx_pending_payments ON payments(created_at) WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_completed_payments ON payments(created_at) WHERE status = 'completed';
CREATE INDEX IF NOT EXISTS idx_failed_payments ON payments(created_at) WHERE status = 'failed';

-- Analyze tables for better query planning
ANALYZE payments;
ANALYZE payment_methods;
ANALYZE refunds;
ANALYZE payment_attempts;
ANALYZE webhooks;
ANALYZE payment_settlements;

-- Add helpful comments for database administrators
COMMENT ON ROLE payment_service IS 'Application user for payment service with full CRUD permissions';
COMMENT ON ROLE payment_analytics IS 'Read-only user for payment analytics and reporting';
COMMENT ON INDEX idx_payments_composite_lookup IS 'Optimized for user payment history lookups';
COMMENT ON INDEX idx_pending_payments IS 'Optimized for processing pending payments queue';
COMMENT ON INDEX idx_webhooks_retry_processing IS 'Optimized for webhook retry processing';

-- Display summary of seed data
DO $$
DECLARE
    payment_count INTEGER;
    refund_count INTEGER;
    webhook_count INTEGER;
    settlement_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO payment_count FROM payments;
    SELECT COUNT(*) INTO refund_count FROM refunds;
    SELECT COUNT(*) INTO webhook_count FROM webhooks;
    SELECT COUNT(*) INTO settlement_count FROM payment_settlements;

    RAISE NOTICE 'Seed data loaded successfully:';
    RAISE NOTICE '- % payments inserted', payment_count;
    RAISE NOTICE '- % refunds inserted', refund_count;
    RAISE NOTICE '- % webhooks inserted', webhook_count;
    RAISE NOTICE '- % settlements inserted', settlement_count;
    RAISE NOTICE '- % payment methods configured', (SELECT COUNT(*) FROM payment_methods WHERE enabled = true);
END $$;