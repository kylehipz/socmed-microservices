-- Create pg extension
CREATE EXTENSION IF NOT EXISTS pgcrypto;
-- Create users schema
CREATE SCHEMA IF NOT EXISTS users;

-- Create users table
CREATE TABLE IF NOT EXISTS users.users (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);
