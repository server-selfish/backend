CREATE TABLE IF NOT EXISTS project(
	id UUID PRIMARY KEY DEFAULT uuidv7(),
  name VARCHAR UNIQUE NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS container(
	id UUID PRIMARY KEY DEFAULT uuidv7(),
  name VARCHAR NOT NULL,
  port_exposed integer[],
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ
);
  
CREATE TABLE IF NOT EXISTS deployment(
	id UUID PRIMARY KEY DEFAULT uuidv7(),
  name VARCHAR NOT NULL,
  s3_remote_url VARCHAR,
  version varchar,
  project_id UUID NOT NULL REFERENCES project(id),
  container_id UUID NOT NULL REFERENCES container(id),
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ,
  CONSTRAINT deployment_project_name_unique UNIQUE (project_id, name)
);
