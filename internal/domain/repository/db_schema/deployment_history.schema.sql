CREATE TABLE IF NOT EXISTS deployment_history(
  id SERIAL PRIMARY KEY,
  deployment_id UUID NOT NULL REFERENCES deployment(id) ON DELETE CASCADE,
  branch VARCHAR NOT NULL,
  commit_msg TEXT NOT NULL,
  created_at TIMESTAMPTZ default now(),
);
