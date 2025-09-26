-- Workflow Engine Database Schema

-- Workflows table
CREATE TABLE IF NOT EXISTS workflows (
    name VARCHAR(255) PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workflow_runs (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    input JSONB NOT NULL,
    workflow_name VARCHAR(255) REFERENCES workflows(name) NOT NULL,
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
)