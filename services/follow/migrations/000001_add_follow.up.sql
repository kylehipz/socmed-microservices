CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE SCHEMA IF NOT EXISTS follow;

CREATE TABLE IF NOT EXISTS follow.users (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  user_name TEXT NOT NULL UNIQUE,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

CREATE TABLE IF NOT EXISTS follow.follow (
  follower_id uuid NOT NULL,
  followee_id uuid NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  PRIMARY KEY (follower_id, followee_id),
  CONSTRAINT fk_follower FOREIGN KEY (follower_id)
    REFERENCES follow.users(id) ON DELETE CASCADE,
  CONSTRAINT fk_followee FOREIGN KEY (followee_id)
    REFERENCES follow.users(id) ON DELETE CASCADE
);
