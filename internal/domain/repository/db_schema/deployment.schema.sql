CREATE TABLE IF NOT EXISTS deployment(
	id UUID PRIMARY KEY DEFAULT uuidv7(),
  name VARCHAR NOT NULL,
  git_remote_url VARCHAR NOT NULL,
  project_id UUID NOT NULL REFERENCES project(id) ON DELETE CASCADE,
  installation_id BIGINT NOT NULL REFERENCES github_installations(installation_id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ,
  CONSTRAINT deployment_project_name_unique UNIQUE (project_id, name)
);

CREATE TABLE IF NOT EXISTS deployment_techstack(
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL,
  version VARCHAR NOT NULL,
  docker_base_image VARCHAR NOT NULL,
  docker_runtime_image VARCHAR NOT NULL,
  created_at TIMESTAMPTZ default now(),
  updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS deployment_history(
  id SERIAL PRIMARY KEY,
  deployment_id UUID NOT NULL REFERENCES deployment(id) ON DELETE CASCADE,
  branch VARCHAR NOT NULL,
  commit_id VARCHAR NOT NULL,
  commit_msg TEXT NOT NULL,
  version VARCHAR NOT NULL,
  external_port integer[],
  deployment_techstack_id INTEGER NOT NULL REFERENCES deployment_techstack(id) ON DELETE CASCADE,
  build_command VARCHAR,
  build_folder VARCHAR,
  run_command VARCHAR,
  is_active BOOL NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ default now(),
  updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS container(
	id UUID PRIMARY KEY DEFAULT uuidv7(),
  name VARCHAR NOT NULL,
  deployment_history_id INTEGER NOT NULL UNIQUE REFERENCES deployment_history(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS container_env(
  id SERIAL PRIMARY KEY,
  container_id UUID NOT NULL REFERENCES container(id) ON DELETE CASCADE,
  key VARCHAR NOT NULL,
  value VARCHAR NOT NULL,
  created_at TIMESTAMPTZ default now(),
  updated_at TIMESTAMPTZ
);


CREATE TABLE IF NOT EXISTS project(
	id UUID PRIMARY KEY DEFAULT uuidv7(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR UNIQUE NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ
);

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
