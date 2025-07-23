-- Create time_series_metrics table
CREATE TABLE time_series_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    product_id UUID,
    question_id UUID,
    metric_type VARCHAR(255) NOT NULL,
    metric_name VARCHAR(255) NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    count BIGINT NOT NULL DEFAULT 0,
    timestamp TIMESTAMPTZ NOT NULL,
    granularity VARCHAR(50) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_time_series_account_id ON time_series_metrics(account_id);
CREATE INDEX idx_time_series_org_id ON time_series_metrics(organization_id);
CREATE INDEX idx_time_series_product_id ON time_series_metrics(product_id) WHERE product_id IS NOT NULL;
CREATE INDEX idx_time_series_question_id ON time_series_metrics(question_id) WHERE question_id IS NOT NULL;
CREATE INDEX idx_time_series_metric_type ON time_series_metrics(metric_type);
CREATE INDEX idx_time_series_timestamp ON time_series_metrics(timestamp);
CREATE INDEX idx_time_series_granularity ON time_series_metrics(granularity);

-- Composite indexes for common query patterns
CREATE INDEX idx_time_series_org_metric_time ON time_series_metrics(organization_id, metric_type, timestamp);
CREATE INDEX idx_time_series_product_metric_time ON time_series_metrics(product_id, metric_type, timestamp) WHERE product_id IS NOT NULL;
CREATE INDEX idx_time_series_question_metric_time ON time_series_metrics(question_id, metric_type, timestamp) WHERE question_id IS NOT NULL;

-- Add foreign key constraints
ALTER TABLE time_series_metrics 
    ADD CONSTRAINT fk_time_series_organization 
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;

-- Note: We don't add FK constraints for product_id and question_id as they might be optional
-- and could reference different tables depending on the metric type