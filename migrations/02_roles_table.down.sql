-- Remove the role_id column from the users table
ALTER TABLE IF EXISTS users
DROP COLUMN IF EXISTS role_id;
DROP TABLE IF EXISTS roles;
