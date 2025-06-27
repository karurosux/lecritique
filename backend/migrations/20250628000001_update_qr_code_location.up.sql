-- Add location column as text field
ALTER TABLE qr_codes ADD COLUMN location TEXT;

-- Drop the old location_id foreign key constraint if it exists
ALTER TABLE qr_codes DROP CONSTRAINT IF EXISTS fk_qr_codes_location;

-- Remove location_id column
ALTER TABLE qr_codes DROP COLUMN IF EXISTS location_id;