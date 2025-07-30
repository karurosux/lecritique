CREATE TABLE IF NOT EXISTS seed_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seed_name VARCHAR(255) NOT NULL UNIQUE,
    executed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version VARCHAR(50),
    metadata JSONB
);

CREATE INDEX idx_seed_runs_name ON seed_runs(seed_name);