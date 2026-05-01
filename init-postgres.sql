-- Create Status Enum
CREATE TYPE work_item_status AS ENUM ('OPEN', 'INVESTIGATING', 'RESOLVED', 'CLOSED');

-- Create Work Items Table
CREATE TABLE IF NOT EXISTS work_items (
    id SERIAL PRIMARY KEY,
    component_id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    severity INT NOT NULL DEFAULT 3, -- 0 for P0, 1 for P1, etc.
    status work_item_status NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create RCA Records Table
CREATE TABLE IF NOT EXISTS rca_records (
    id SERIAL PRIMARY KEY,
    work_item_id INTEGER REFERENCES work_items(id) ON DELETE CASCADE UNIQUE,
    root_cause_category VARCHAR(100) NOT NULL,
    fix_applied TEXT NOT NULL,
    prevention_steps TEXT,
    mttr_minutes INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
