-- Migration: Remove unused columns from players table
-- These columns are not used by any parser: registry_id, school, rank, citizenship, role

-- Remove unused columns
ALTER TABLE players DROP COLUMN IF EXISTS registry_id;
ALTER TABLE players DROP COLUMN IF EXISTS school;
ALTER TABLE players DROP COLUMN IF EXISTS rank;
ALTER TABLE players DROP COLUMN IF EXISTS citizenship;
ALTER TABLE players DROP COLUMN IF EXISTS role;
