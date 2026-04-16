CREATE TABLE IF NOT EXISTS user(
  id UUID PRIMARY KEY,
  github_id BIGINT UNIQUE NOT NULL,
  username TEXT NOT NULL,
  email TEXT,
  avatar_url TEXT,
  created_at TIMESTAMPTZ default now(),
  updated_at TIMESTAMPTZ
);
