-- 002_add_system_to_users.sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS system JSONB DEFAULT '[]'::jsonb;
    