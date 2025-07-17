-- Migration: Fix owner team member records with invalid member_id

-- Update team member records where member_id is null or zero UUID
-- For OWNER role, member_id should be the same as account_id
UPDATE team_members 
SET member_id = account_id
WHERE role = 'OWNER' 
AND (member_id IS NULL OR member_id = '00000000-0000-0000-0000-000000000000');

-- Also update any other team members with zero member_id to match invited_by
-- (assuming invited_by is a valid account ID)
UPDATE team_members 
SET member_id = invited_by
WHERE role != 'OWNER' 
AND (member_id IS NULL OR member_id = '00000000-0000-0000-0000-000000000000')
AND invited_by != '00000000-0000-0000-0000-000000000000';

-- Log any remaining invalid records for manual review
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM team_members 
        WHERE member_id IS NULL OR member_id = '00000000-0000-0000-0000-000000000000'
    ) THEN
        RAISE NOTICE 'Warning: Some team_member records still have invalid member_id values';
    END IF;
END $$;