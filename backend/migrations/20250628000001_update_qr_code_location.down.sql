-- Revert the changes: add back location_id and remove location text field
ALTER TABLE qr_codes ADD COLUMN location_id UUID;

-- Remove the location text column
ALTER TABLE qr_codes DROP COLUMN IF EXISTS location;