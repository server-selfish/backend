CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  provider TEXT NOT NULL DEFAULT 'github',
  provider_user_id BIGINT NOT NULL,
  username TEXT NOT NULL,
  email TEXT,
  avatar_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ,
  CONSTRAINT users_provider_provider_user_id_key UNIQUE (provider, provider_user_id)
);

CREATE TABLE IF NOT EXISTS project(
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR UNIQUE NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ
);
