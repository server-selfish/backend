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
  git_remote_url VARCHAR NOT NULL,
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
