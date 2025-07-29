-- Add address column to organizations table
ALTER TABLE organizations ADD COLUMN address VARCHAR(500);

-- Drop locations table as we're moving to single address field
DROP TABLE IF EXISTS locations;