CREATE TABLE generation_jobs (
    id UUID PRIMARY KEY,
    org_id UUID NOT NULL REFERENCES orgs(id),
    template_id UUID NOT NULL REFERENCES templates(id),
    status VARCHAR(50) NOT NULL, -- pending, processing, completed, failed
    output_asset_id UUID REFERENCES assets(id),
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_jobs_org_id ON generation_jobs(org_id);
