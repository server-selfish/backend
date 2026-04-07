CREATE TABLE IF NOT EXISTS container(
	id UUID PRIMARY KEY DEFAULT uuidv7(),
  name VARCHAR NOT NULL,
  port_exposed integer[],
  deployment_id UUID NOT NULL UNIQUE REFERENCES deployment(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ
);
