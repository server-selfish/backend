CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  provider TEXT NOT NULL DEFAULT 'github',
  provider_user_id BIGINT NOT NULL,
  username TEXT NOT NULL,
  email TEXT,
  avatar_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now(),
  CONSTRAINT users_provider_provider_user_id_key UNIQUE (provider, provider_user_id)
);

CREATE TABLE IF NOT EXISTS github_installations (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  installation_id BIGINT NOT NULL UNIQUE,
  account_login TEXT,
  account_id BIGINT,
  target_type TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_github_installations_user_installation
  ON github_installations (user_id, installation_id);

CREATE INDEX IF NOT EXISTS idx_github_installations_user_id
  ON github_installations (user_id);

CREATE INDEX IF NOT EXISTS idx_github_installations_installation_id
  ON github_installations (installation_id);
