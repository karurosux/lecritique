-- Revert the trigger fix (not recommended, but provided for completeness)

DROP TRIGGER IF EXISTS create_account_owner_trigger ON accounts;
DROP FUNCTION IF EXISTS create_account_owner();