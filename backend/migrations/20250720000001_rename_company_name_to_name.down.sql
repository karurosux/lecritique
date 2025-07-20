-- Revert name column back to company_name in accounts table
ALTER TABLE accounts RENAME COLUMN name TO company_name;