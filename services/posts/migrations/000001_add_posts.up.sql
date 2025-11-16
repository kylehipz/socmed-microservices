-- Create pg extension
CREATE EXTENSION IF NOT EXISTS pgcrypto;
-- Create posts schema
CREATE SCHEMA IF NOT EXISTS posts;

-- Create posts table
CREATE TABLE IF NOT EXISTS posts.posts (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  author_id uuid NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create index on author id
CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts.posts (author_id);
