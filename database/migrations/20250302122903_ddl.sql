CREATE TABLE workflows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL,
    latest_version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_workflows_owner ON workflows(owner_id);

CREATE TABLE workflow_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL,
    version INT NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE,
    UNIQUE (workflow_id, version)  -- Ensures uniqueness for FK references
);

CREATE TABLE workflow_nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL,
    version INT NOT NULL,
    type TEXT CHECK (type IN ('event_trigger', 'api_call', 'condition', 'loop', 'delay', 'slack', 'gmail', 'outlook', 'calendar', 's3', 'mysql', 'mongodb', 'neo4j')),
    credential_id UUID NULL,  -- For Secure API Calls
    metadata JSONB NOT NULL,  -- Stores dynamic execution details
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (workflow_id, version) REFERENCES workflow_versions(workflow_id, version) ON DELETE CASCADE
);

CREATE INDEX idx_nodes_workflow ON workflow_nodes(workflow_id, version);


CREATE TABLE node_connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL,
    version INT NOT NULL,
    source_node UUID NOT NULL,
    target_node UUID NOT NULL,
    condition JSONB NULL,  -- {"rule": "$.status == 'success'"}
    parallel_group TEXT NULL,  -- Grouping for parallel execution
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (source_node) REFERENCES workflow_nodes(id) ON DELETE CASCADE,
    FOREIGN KEY (target_node) REFERENCES workflow_nodes(id) ON DELETE CASCADE
);

CREATE INDEX idx_connections_workflow ON node_connections(workflow_id, version);

CREATE TABLE execution_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL,
    version INT NOT NULL,
    status TEXT CHECK (status IN ('running', 'completed', 'failed', 'paused')),
    started_at TIMESTAMP DEFAULT NOW(),
    finished_at TIMESTAMP NULL,
    FOREIGN KEY (workflow_id, version) REFERENCES workflow_versions(workflow_id, version) ON DELETE CASCADE
);

CREATE INDEX idx_execution_workflow ON execution_instances(workflow_id, version);

CREATE TABLE execution_nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    execution_id UUID NOT NULL,
    node_id UUID NOT NULL,
    status TEXT CHECK (status IN ('running', 'completed', 'failed', 'skipped')),
    output_data JSONB NULL,  -- Stores API responses, conditions met
    started_at TIMESTAMP DEFAULT NOW(),
    finished_at TIMESTAMP NULL,
    FOREIGN KEY (execution_id) REFERENCES execution_instances(id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES workflow_nodes(id) ON DELETE CASCADE
);

CREATE INDEX idx_execution_nodes_exec ON execution_nodes(execution_id);

CREATE TABLE workflow_credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,
    name TEXT NOT NULL,
    type TEXT CHECK (type IN ('oauth', 'api_key', 'basic_auth')),
    credentials JSONB NOT NULL,  -- Securely stored credentials
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL, -- e.g., "Stripe API Key"
    type TEXT NOT NULL CHECK (type IN ('api_key', 'oauth', 'jwt', 'basic_auth')),
    encrypted_data BYTEA NOT NULL, -- Encrypted credential storage
    refresh_token BYTEA, -- (Optional) OAuth refresh token
    expires_at TIMESTAMP, -- Expiry time for OAuth tokens
    metadata JSONB DEFAULT '{}'::JSONB, -- API-specific settings
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE credential_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    credential_id UUID REFERENCES credentials(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    permission TEXT CHECK (permission IN ('read', 'write', 'execute')),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE workflow_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL,
    version INT NOT NULL,
    action TEXT CHECK (action IN ('created', 'updated', 'deleted', 'executed', 'failed')),
    details JSONB NOT NULL,  -- Logs what changed
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (workflow_id, version) REFERENCES workflow_versions(workflow_id, version) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS workflow_logs;
DROP TABLE IF EXISTS credential_permissions;
DROP TABLE IF EXISTS credentials;
DROP TABLE IF EXISTS workflow_credentials;
DROP TABLE IF EXISTS execution_nodes;
DROP TABLE IF EXISTS execution_instances;
DROP TABLE IF EXISTS node_connections;
DROP TABLE IF EXISTS workflow_nodes;
DROP TABLE IF EXISTS workflow_versions;
DROP TABLE IF EXISTS workflows;
-- +goose StatementEnd
