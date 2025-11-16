CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE SCHEMA IF NOT EXISTS timeline;

CREATE TABLE IF NOT EXISTS timeline.users (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  user_name TEXT NOT NULL UNIQUE,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

CREATE TABLE IF NOT EXISTS timeline.follow (
  follower_id uuid NOT NULL,
  followee_id uuid NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  PRIMARY KEY (follower_id, followee_id),
  CONSTRAINT fk_follower FOREIGN KEY (follower_id)
    REFERENCES timeline.users(id) ON DELETE CASCADE,
  CONSTRAINT fk_followee FOREIGN KEY (followee_id)
    REFERENCES timeline.users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS timeline.posts (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  author_id uuid NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now(),
  CONSTRAINT fk_author FOREIGN KEY (author_id)
    REFERENCES timeline.users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_posts_author_id ON timeline.posts (author_id);
