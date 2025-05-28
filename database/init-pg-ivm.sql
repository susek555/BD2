CREATE EXTENSION IF NOT EXISTS pg_ivm;

-- Necessary privileges
GRANT USAGE ON SCHEMA public TO bd2_user;
GRANT CREATE ON SCHEMA public TO bd2_user;
