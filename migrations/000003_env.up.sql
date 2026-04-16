CREATE TABLE IF NOT EXISTS container_env(
  id SERIAL PRIMARY KEY,
  key VARCHAR NOT NULL,
  value VARCHAR NOT NULL,
  created_at TIMESTAMPTZ default now(),
  updated_at TIMESTAMPTZ
);
