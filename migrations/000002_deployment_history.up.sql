CREATE TABLE IF NOT EXISTS deployment_techstack(
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL, --js, go, python
  version VARCHAR NOT NULL
)

CREATE TABLE IF NOT EXISTS deployment_history(
  id SERIAL PRIMARY KEY,
  deployment_id UUID NOT NULL REFERENCES deployment(id) ON DELETE CASCADE,
  git_remote_url VARCHAR NOT NULL,
  branch VARCHAR NOT NULL,
  commit_msg TEXT NOT NULL,
  version VARCHAR NOT NULL,
  external_port integer[],
  deployment_techstack INTEGER NOT NULL REFERENCES deployment_techstack(id) ON DELETE CASCADE,
  build_command VARCHAR,
  created_at TIMESTAMPTZ default now(),
);
