CREATE TABLE IF NOT EXISTS project(
	id UUID PRIMARY KEY DEFAULT uuidv7(),
  name VARCHAR UNIQUE NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ
);
