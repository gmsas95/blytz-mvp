-- Migration: Create enhanced payments table with Fiuu integration
-- Version: 001
-- Created: 2024-01-01

-- Create payments table with comprehensive fields for Fiuu integration
CREATE TABLE IF NOT EXISTS payments (
    -- Primary fields
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    order_id VARCHAR(100) NOT NULL,
    transaction_id VARCHAR(100),

    -- Amount and currency
    amount DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'MYR',

    -- Payment method details
    payment_method VARCHAR(50) NOT NULL,
    payment_channel VARCHAR(50),
    payment_subchannel VARCHAR(50),

    -- Status tracking
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    gateway_status VARCHAR(50),
    gateway_response JSONB,

    -- Billing information
    bill_name VARCHAR(100) NOT NULL,
    bill_email VARCHAR(255) NOT NULL,
    bill_mobile VARCHAR(20) NOT NULL,
    bill_description TEXT NOT NULL,

    -- URLs
    return_url VARCHAR(500),
    notify_url VARCHAR(500),

    -- Fiuu specific fields
    fiuu_merchant_id VARCHAR(50),
    fiuu_reference_no VARCHAR(100),
    fiuu_approval_code VARCHAR(50),
    fiuu_error_code VARCHAR(50),
    fiuu_error_description TEXT,
    fiuu_signature VARCHAR(255),

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    expired_at TIMESTAMP WITH TIME ZONE,

    -- Metadata
    metadata JSONB,
    created_by VARCHAR(100),
    ip_address INET,
    user_agent TEXT,

    -- Constraints
    CONSTRAINT payments_order_id_user_id_unique UNIQUE(order_id, user_id),
    CONSTRAINT payments_transaction_id_unique UNIQUE(transaction_id),
    CONSTRAINT payments_status_check CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'refunded', 'cancelled', 'expired')),
    CONSTRAINT payments_currency_check CHECK (currency IN ('MYR', 'SGD', 'PHP', 'THB', 'USD', 'IDR'))
);

-- Create indexes for performance
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_transaction_id ON payments(transaction_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_payment_method ON payments(payment_method);
CREATE INDEX idx_payments_created_at ON payments(created_at);
CREATE INDEX idx_payments_user_id_created_at ON payments(user_id, created_at);
CREATE INDEX idx_payments_status_created_at ON payments(status, created_at);
CREATE INDEX idx_payments_fiuu_reference ON payments(fiuu_reference_no);
CREATE INDEX idx_payments_gateway_status ON payments(gateway_status);

-- Create payment methods table
CREATE TABLE IF NOT EXISTS payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    description TEXT,
    icon_url VARCHAR(500),

    -- Configuration
    enabled BOOLEAN DEFAULT true,
    supported_currencies TEXT[] DEFAULT ARRAY['MYR'],
    min_amount DECIMAL(12,2),
    max_amount DECIMAL(12,2),

    -- Processing configuration
    processing_time_seconds INTEGER DEFAULT 30,
    requires_3ds BOOLEAN DEFAULT false,
    requires_redirect BOOLEAN DEFAULT false,

    -- Metadata
    country_codes TEXT[] DEFAULT ARRAY['MY'],
    metadata JSONB,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT payment_methods_type_check CHECK (type IN ('bank_transfer', 'ewallet', 'credit_card', 'qr_code', 'bnpl', 'crypto'))
);

-- Insert default payment methods for Malaysian market
INSERT INTO payment_methods (code, name, type, description, supported_currencies, country_codes) VALUES
('FPX', 'FPX Online Banking', 'bank_transfer', 'Pay directly from your bank account', ARRAY['MYR'], ARRAY['MY']),
('FPX_BB', 'FPX Business Banking', 'bank_transfer', 'Business banking transfers via FPX', ARRAY['MYR'], ARRAY['MY']),
('GRABPAY', 'GrabPay', 'ewallet', 'Pay using your GrabPay e-wallet', ARRAY['MYR'], ARRAY['MY']),
('TNG', 'Touch ''n Go eWallet', 'ewallet', 'Pay using Touch ''n Go e-wallet', ARRAY['MYR'], ARRAY['MY']),
('SHOPEEPAY', 'ShopeePay', 'ewallet', 'Pay using ShopeePay e-wallet', ARRAY['MYR'], ARRAY['MY']),
('BOOST', 'Boost', 'ewallet', 'Pay using Boost e-wallet', ARRAY['MYR'], ARRAY['MY']),
('MAYBANKQR', 'MAE QR', 'qr_code', 'Scan to pay with Maybank MAE', ARRAY['MYR'], ARRAY['MY']),
('DUITNOWQR', 'DuitNow QR', 'qr_code', 'Universal QR code payment', ARRAY['MYR'], ARRAY['MY']),
('ATOME', 'Atome', 'bnpl', 'Buy now, pay later with Atome', ARRAY['MYR'], ARRAY['MY']),
('RELY', 'Rely', 'bnpl', 'Buy now, pay later with Rely', ARRAY['MYR'], ARRAY['MY']),
('VISA', 'Visa Credit/Debit', 'credit_card', 'Pay with Visa cards', ARRAY['MYR'], ARRAY['MY']),
('MASTERCARD', 'Mastercard Credit/Debit', 'credit_card', 'Pay with Mastercard cards', ARRAY['MYR'], ARRAY['MY'])
ON CONFLICT (code) DO NOTHING;

-- Create refunds table
CREATE TABLE IF NOT EXISTS refunds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    refund_id VARCHAR(100) NOT NULL UNIQUE,

    -- Refund details
    amount DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL,
    reason TEXT NOT NULL,

    -- Status tracking
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    gateway_status VARCHAR(50),
    gateway_response JSONB,

    -- Fiuu specific fields
    fiuu_refund_id VARCHAR(100),
    fiuu_error_code VARCHAR(50),
    fiuu_error_description TEXT,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,

    -- Metadata
    processed_by VARCHAR(100),
    metadata JSONB,

    -- Constraints
    CONSTRAINT refunds_status_check CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled'))
);

-- Create indexes for refunds
CREATE INDEX idx_refunds_payment_id ON refunds(payment_id);
CREATE INDEX idx_refunds_refund_id ON refunds(refund_id);
CREATE INDEX idx_refunds_status ON refunds(status);
CREATE INDEX idx_refunds_created_at ON refunds(created_at);

-- Create payment_attempts table for tracking retry attempts
CREATE TABLE IF NOT EXISTS payment_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    attempt_number INTEGER NOT NULL,

    -- Attempt details
    status VARCHAR(20) NOT NULL,
    gateway_request JSONB,
    gateway_response JSONB,
    error_code VARCHAR(50),
    error_message TEXT,

    -- Timing
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    duration_ms INTEGER,

    -- Metadata
    ip_address INET,
    user_agent TEXT,

    -- Constraints
    CONSTRAINT payment_attempts_status_check CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'timeout'))
);

-- Create indexes for payment attempts
CREATE INDEX idx_payment_attempts_payment_id ON payment_attempts(payment_id);
CREATE INDEX idx_payment_attempts_status ON payment_attempts(status);
CREATE INDEX idx_payment_attempts_started_at ON payment_attempts(started_at);

-- Create webhooks table for tracking webhook events
CREATE TABLE IF NOT EXISTS webhooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(50) NOT NULL,

    -- Source information
    source VARCHAR(50) NOT NULL,
    source_id VARCHAR(100),

    -- Payload
    payload JSONB NOT NULL,
    signature VARCHAR(255),

    -- Processing status
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    processed_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    next_retry_at TIMESTAMP WITH TIME ZONE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT webhooks_status_check CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'retrying'))
);

-- Create indexes for webhooks
CREATE INDEX idx_webhooks_event_type ON webhooks(event_type);
CREATE INDEX idx_webhooks_status ON webhooks(status);
CREATE INDEX idx_webhooks_source ON webhooks(source, source_id);
CREATE INDEX idx_webhooks_created_at ON webhooks(created_at);
CREATE INDEX idx_webhooks_next_retry ON webhooks(next_retry_at) WHERE status = 'retrying';

-- Create payment_settlements table for tracking settlements
CREATE TABLE IF NOT EXISTS payment_settlements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,

    -- Settlement details
    settlement_id VARCHAR(100) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    fee_amount DECIMAL(12,2) DEFAULT 0,
    net_amount DECIMAL(12,2) NOT NULL,

    -- Settlement information
    settlement_date DATE NOT NULL,
    batch_id VARCHAR(100),

    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'pending',

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    settled_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT payment_settlements_status_check CHECK (status IN ('pending', 'processing', 'completed', 'failed'))
);

-- Create indexes for settlements
CREATE INDEX idx_payment_settlements_payment_id ON payment_settlements(payment_id);
CREATE INDEX idx_payment_settlements_settlement_date ON payment_settlements(settlement_date);
CREATE INDEX idx_payment_settlements_status ON payment_settlements(status);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update updated_at
CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payment_methods_updated_at BEFORE UPDATE ON payment_methods
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_refunds_updated_at BEFORE UPDATE ON refunds
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_webhooks_updated_at BEFORE UPDATE ON webhooks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create view for payment analytics
CREATE OR REPLACE VIEW payment_analytics AS
SELECT
    DATE_TRUNC('day', created_at) as date,
    payment_method,
    status,
    currency,
    COUNT(*) as total_payments,
    SUM(amount) as total_amount,
    AVG(amount) as avg_amount,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as successful_payments,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_payments,
    ROUND(
        (COUNT(CASE WHEN status = 'completed' THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0)), 2
    ) as success_rate_percent
FROM payments
WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY DATE_TRUNC('day', created_at), payment_method, status, currency
ORDER BY date DESC, payment_method;

-- Create view for daily settlement summary
CREATE OR REPLACE VIEW daily_settlement_summary AS
SELECT
    settlement_date,
    currency,
    COUNT(*) as total_transactions,
    SUM(amount) as gross_amount,
    SUM(fee_amount) as total_fees,
    SUM(net_amount) as net_amount,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as settled_transactions
FROM payment_settlements
WHERE settlement_date >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY settlement_date, currency
ORDER BY settlement_date DESC;

-- Add comments for documentation
COMMENT ON TABLE payments IS 'Main payments table with Fiuu integration support';
COMMENT ON TABLE payment_methods IS 'Available payment methods and their configurations';
COMMENT ON TABLE refunds IS 'Refund transactions linked to payments';
COMMENT ON TABLE payment_attempts IS 'Track retry attempts for payment processing';
COMMENT ON TABLE webhooks IS 'Incoming webhook events from payment gateways';
COMMENT ON TABLE payment_settlements IS 'Payment settlement information for reconciliation';
COMMENT ON VIEW payment_analytics IS 'Daily payment analytics and success rates';
COMMENT ON VIEW daily_settlement_summary IS 'Daily settlement summary for finance reconciliation';

-- Create user-friendly payment status enum type
CREATE TYPE payment_status AS ENUM (
    'pending',
    'processing',
    'completed',
    'failed',
    'refunded',
    'cancelled',
    'expired'
);

-- Create user-friendly refund status enum type
CREATE TYPE refund_status AS ENUM (
    'pending',
    'processing',
    'completed',
    'failed',
    'cancelled'
);

-- Note: To use these enums, you would need to modify the status columns:
-- ALTER TABLE payments ALTER COLUMN status TYPE payment_status USING status::payment_status;
-- ALTER TABLE refunds ALTER COLUMN status TYPE refund_status USING status::refund_status;