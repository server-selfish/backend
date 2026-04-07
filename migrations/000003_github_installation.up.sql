CREATE TABLE IF NOT EXISTS github_installation (
  id SERIAL PRIMARY KEY,
  installation_id BIGINT NOT NULL,
  created_at TIMESTAMPZ DEFAULT now()
);
