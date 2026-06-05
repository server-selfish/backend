-- name: GetAllProjects :many
SELECT
  *
FROM public.project p
WHERE
  p.user_id = $1
ORDER BY
  COALESCE(p.updated_at, p.created_at) DESC;

-- name: GetProjectById :one
SELECT
  *
FROM
  public.project p
WHERE
  p.id = $1
  AND p.user_id = $2;

-- name: GetProjectByName :one
SELECT
  p.name AS project_name
  ,p.description AS project_description
  ,p.created_at AS project_created_at
  ,p.updated_at AS project_updated_at
  ,COALESCE(
    json_agg(
      DISTINCT jsonb_build_object(
        'deployment_name', d.name,
        'techstack_name', dt.name,
        'container_name', c.name
      )
    ) FILTER (WHERE d.name IS NOT NULL),
    '[]'
  )::jsonb AS deployments
FROM
  public.project p
LEFT JOIN deployment d
	ON d.project_id = p.id
LEFT JOIN deployment_history dh
	ON dh.deployment_id = d.id
		AND dh.is_active IS true
LEFT JOIN deployment_techstack dt
	ON dt.id = dh.deployment_techstack_id
LEFT JOIN container c
	ON c.deployment_history_id = dh.id
WHERE
  p.user_id = $1
  AND p.name ILIKE $2
GROUP BY
  p.name, p.description, p.created_at, p.updated_at;

-- name: GetProjectByNameDetail :one
WITH container_ports AS(
	SELECT
	  cp.container_id,
	  COALESCE(
	    json_agg(
	      DISTINCT jsonb_build_object(
	        'external', cp.external,
	        'internal', cp.internal,
	        'protocol', cp.protocol
	      )
	    ), '[]'
	  )::jsonb as port
	FROM container_port cp
	GROUP BY cp.container_id
)
SELECT
  p.name AS project_name
  ,p.description AS project_description
  ,p.created_at AS project_created_at
  ,p.updated_at AS project_updated_at
  ,COALESCE(
    json_agg(
      DISTINCT jsonb_build_object(
        'deployment_name', d.name,
        'deployment_url', d.git_remote_url,
        'deployment_created_at',d.created_at,
        'deployment_updated_at',d.updated_at,
        'deployment_branch',dh.branch,
        'deployment_version',dh."version",
        'deployment_commit_msg',dh.commit_msg,
        'deployment_port',cp.port,
        'deployment_history_created_at',dh.created_at,
        'techstack_name', dt.name,
        'techstack_version', dt.version,
        'container_name', c.name
      )
    ) FILTER (WHERE d.name IS NOT NULL),
    '[]'
  )::jsonb AS deployments
FROM
  public.project p
LEFT JOIN deployment d
	ON d.project_id = p.id
LEFT JOIN deployment_history dh
	ON dh.deployment_id = d.id
		AND dh.is_active IS true
LEFT JOIN deployment_techstack dt
	ON dt.id = dh.deployment_techstack_id
LEFT JOIN container c
	ON c.deployment_history_id = dh.id
LEFT JOIN container_ports cp
  ON cp.container_id = c.id
WHERE
  p.user_id = $1
  AND p.name ILIKE $2
GROUP BY
  p.name, p.description, p.created_at, p.updated_at;

-- name: CreateProject :exec
INSERT INTO public.project (user_id,name,description)
VALUES ($1,$2,$3);

-- name: UpdateProjectById :exec
UPDATE public.project
SET name = $1,
    description = $2,
    updated_at = now()
WHERE
  id = $3
  AND user_id = $4;

-- name: DeleteProjectById :exec
DELETE FROM public.project
WHERE
  id = $1
  AND user_id = $2;
