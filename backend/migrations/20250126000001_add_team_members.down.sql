-- Drop trigger and function
DROP TRIGGER IF EXISTS create_account_owner_trigger ON accounts;
DROP FUNCTION IF EXISTS create_account_owner();

-- Drop tables
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS users;