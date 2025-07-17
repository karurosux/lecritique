-- Rollback: This migration fixed invalid member_ids
-- Since we can't reliably restore the previous invalid state,
-- this down migration is a no-op to prevent data corruption

-- Log a warning
DO $$
BEGIN
    RAISE NOTICE 'Warning: Down migration for fix_owner_member_ids cannot restore previous invalid states';
END $$;