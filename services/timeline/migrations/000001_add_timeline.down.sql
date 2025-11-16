-- Drop the users table if it exists
DROP TABLE IF EXISTS timeline.posts CASCADE;
DROP TABLE IF EXISTS timeline.users CASCADE;
DROP TABLE IF EXISTS timeline.follow CASCADE;

-- Drop the users schema if it's now empty
DROP SCHEMA IF EXISTS timeline CASCADE;

