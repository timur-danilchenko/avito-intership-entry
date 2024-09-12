-- CREATE TYPE organization_type AS ENUM (
--     'IE',
--     'LLC',
--     'JSC'
-- );

CREATE TABLE IF NOT EXISTS organization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE,
    user_id INT REFERENCES employee(id) ON DELETE CASCADE
);