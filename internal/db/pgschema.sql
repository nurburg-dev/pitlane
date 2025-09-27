CREATE TABLE IF NOT EXISTS workflows (
    name VARCHAR(255) PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS workflow_runs (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    input JSONB NOT NULL,
    workflow_name VARCHAR(255) REFERENCES workflows(name) NOT NULL,
    status VARCHAR(255) NOT NULL,
    scheduled_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS activity_runs (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    activity_name VARCHAR(255) NOT NULL,
    workflow_run_id VARCHAR(255) REFERENCES workflow_runs(id) NOT NULL,
    errorMessage TEXT,
    input JSONB NOT NULL,
    output JSONB,
    status VARCHAR(255) NOT NULL,
    retry_status JSONB,
    scheduled_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- Indexes for optimal pending task fetching (latest scheduled first)
CREATE INDEX IF NOT EXISTS idx_workflow_runs_pending ON workflow_runs (status, scheduled_at DESC) WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_activity_runs_pending ON activity_runs (status, scheduled_at DESC) WHERE status = 'pending';

-- Index for activity run history by workflow run ID
CREATE INDEX IF NOT EXISTS idx_activity_runs_workflow_history ON activity_runs (workflow_run_id, created_at ASC);
