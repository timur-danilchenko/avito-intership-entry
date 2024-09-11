-- CREATE TYPE tender_status_type AS ENUM (
--     'Created',
--     'Published',
--     'Closed'
-- );

-- CREATE TYPE tender_service_type AS ENUM (
--     'Construction',
--     'Delivery',
--     'Manufacture'
-- );

CREATE TABLE IF NOT EXISTS tenders (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    service_type tender_service_type,
    status tender_status_type,
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
