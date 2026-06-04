CREATE TABLE IF NOT EXISTS deployment_techstack(
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL,
  version VARCHAR NOT NULL,
  docker_base_image VARCHAR NOT NULL,
  docker_runtime_image VARCHAR NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS deployment_history(
  id SERIAL PRIMARY KEY,
  deployment_id UUID NOT NULL REFERENCES deployment(id) ON DELETE CASCADE,
  branch VARCHAR NOT NULL,
  commit_id VARCHAR NOT NULL,
  commit_msg TEXT NOT NULL,
  version VARCHAR NOT NULL,
  deployment_techstack_id INTEGER NOT NULL REFERENCES deployment_techstack(id) ON DELETE CASCADE,
  build_command VARCHAR,
  build_folder VARCHAR,
  run_command VARCHAR,
  is_active BOOL NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS container(
	id UUID PRIMARY KEY DEFAULT uuidv7(),
  name VARCHAR NOT NULL,
  deployment_history_id INTEGER NOT NULL UNIQUE REFERENCES deployment_history(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS container_port(
  id SERIAL PRIMARY KEY,
  container_id UUID NOT NULL REFERENCES container(id) ON DELETE CASCADE,
  external INTEGER NOT NULL CHECK (external BETWEEN 1 AND 65535),
  internal INTEGER NOT NULL CHECK (internal BETWEEN 1 AND 65535),
  protocol VARCHAR(10) NOT NULL CHECK (protocol IN ('tcp', 'udp')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);
