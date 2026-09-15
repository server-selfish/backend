-- name: GetDeploymentsByProjectId :many
SELECT
  CAST (d.id AS VARCHAR) AS deployment_id,
  d.name AS deployment_name,
  d.git_remote_url as git_remote_url,
  CAST (dh.branch AS VARCHAR) AS branch,
  CAST (dh.commit_id AS VARCHAR) AS commit_id,
  CAST (dh.commit_msg AS VARCHAR) AS commit_message,
  CAST (dh.version AS VARCHAR) AS deployment_version,
  COALESCE(
    json_agg(
      DISTINCT jsonb_build_object(
        'external', cp.external,
        'internal', cp.internal,
        'protocol', cp.protocol
      )
    ) FILTER (WHERE cp.external IS NOT NULL), '[]'
  )::jsonb as port,
  dt.name AS techstack_name,
  dt.version AS techstack_version,
  CAST (c.id AS VARCHAR) AS container_id,
  dh.created_at AS created_at,
  dh.updated_at AS updated_at
FROM
  deployment d
JOIN project p
  ON p.id = d.project_id
LEFT JOIN deployment_history dh
  ON d.id = dh.deployment_id
JOIN deployment_techstack dt
  ON dh.deployment_techstack_id = dt.id
LEFT JOIN container c
  ON c.deployment_history_id = dh.id
LEFT JOIN container_port cp
  ON cp.container_id = c.id
WHERE
  p.user_id = $1
  AND d.project_id = $2
  AND dh.is_active = true
GROUP BY
  deployment_id;

-- name: GetDeploymentByDeploymentId :one
SELECT
  d.*
FROM deployment d
JOIN project p
  ON p.id = d.project_id
WHERE
  p.user_id = $1
  AND d.id = $2;

-- name: GetProjectByDeploymentId :one
SELECT
  p.*
FROM deployment d
JOIN project p
  ON d.project_id = p.id
WHERE
  p.user_id = $1
  AND d.id = $2;

-- name: GetDeploymentHistoryByDeploymentName :many
SELECT
  dh.id,
  dh.branch,
  dh.commit_id,
  dh.commit_msg AS commit_message,
  dh.version AS deployment_version,
  dh.build_command AS build_command,
  dt.id AS techstack_id,
  dt.name AS techstack_name,
  dt.version AS techstack_version,
  dh.created_at,
  dh.updated_at
FROM deployment_history dh
JOIN deployment d
  ON dh.deployment_id = d.id
JOIN project p
  ON d.project_id = p.id
JOIN deployment_techstack dt
  ON dh.deployment_techstack_id = dt.id
WHERE
  p.user_id = $1
  AND p.name ILIKE $2
  AND d.name ILIKE $3
  AND dh.is_active = false
GROUP BY
  dh.id,
  dh.branch,
  dh.commit_id,
  dh.commit_msg,
  dh.version,
  dh.build_command,
  dt.id,
  dt.name,
  dt.version,
  dh.created_at,
  dh.updated_at
ORDER BY
  COALESCE(dh.updated_at, dh.created_at) DESC;

-- name: GetActiveDeploymentDetailByDeploymentName :one
SELECT
  p.name as project_name,
  d.id,
  d.name as deployment_name,
  d.description as deployment_description,
  d.repository_id,
  d.git_remote_url,
  d.installation_id,
  dh.branch,
  dh.version,
  dh.commit_id,
  dh.build_command,
  dh.build_folder,
  dh.run_command,
  dh.main_file_path,
  dt.id as techstack_id,
  dt.name as techstack_name,
  dt.docker_base_image,
  dt.docker_runtime_image,
  c.name as container_name,
  COALESCE(
    json_agg(
      DISTINCT jsonb_build_object(
        'key', ce.key,
        'value', ce.value
      )
    ) FILTER (WHERE ce.key IS NOT NULL), '[]'
  )::jsonb as env,
  COALESCE(
    json_agg(
      DISTINCT jsonb_build_object(
        'external', cp.external,
        'internal', cp.internal,
        'protocol', cp.protocol
      )
    ) FILTER (WHERE cp.external IS NOT NULL), '[]'
  )::jsonb as port
FROM deployment_history dh
JOIN deployment d
  ON dh.deployment_id = d.id
JOIN project p
  ON d.project_id = p.id
JOIN deployment_techstack dt
  ON dh.deployment_techstack_id = dt.id
JOIN container c
  ON c.deployment_history_id = dh.id
LEFT JOIN container_env ce
  ON ce.container_id = c.id
LEFT JOIN container_port cp
  ON cp.container_id = c.id
WHERE
  p.user_id = $1
  AND p.name ILIKE $2
  AND d.name ILIKE $3
  AND dh.is_active = true
GROUP BY
  p.id,
  d.id,
  dh.id,
  dt.id,
  c.id
ORDER BY
  COALESCE(dh.updated_at, dh.created_at) DESC;

-- name: GetActiveDeploymentHistoryByDeploymentName :one
SELECT
  gi.account_login as github_account,
  d.name as deployment_name,
  d.description as deployment_description,
  d.installation_id as installation_id,
  d.git_remote_url AS remote_url,
  dh.id AS deployment_history_id,
  dh.branch AS branch,
  dh.commit_id AS commit_id,
  dh.commit_msg AS commit_message,
  dh.version AS deployment_version,
  COALESCE(
    json_agg(
      DISTINCT jsonb_build_object(
        'key', ce.key,
        'value', ce.value
      )
    ) FILTER (WHERE ce.key IS NOT NULL), '[]'
  )::jsonb as env,
  COALESCE(
    json_agg(
      DISTINCT jsonb_build_object(
        'external', cp.external,
        'internal', cp.internal,
        'protocol', cp.protocol
      )
    ) FILTER (WHERE cp.external IS NOT NULL), '[]'
  )::jsonb as port,
  dh.build_command AS build_command,
  dh.build_folder AS build_folder,
  dh.main_file_path as main_file_path,
  dt.id AS techstack_id,
  dt.name AS techstack_name,
  dt.version AS techstack_version,
  c.name as container_name,
  dh.created_at,
  d.updated_at as deployment_updated_at,
  dh.updated_at
FROM deployment_history dh
JOIN container c
  ON c.deployment_history_id = dh.id
LEFT JOIN container_port cp
  ON cp.container_id = c.id
LEFT JOIN container_env ce
  ON ce.container_id = c.id
JOIN deployment d
  ON dh.deployment_id = d.id
JOIN project p
  ON d.project_id = p.id
JOIN deployment_techstack dt
  ON dh.deployment_techstack_id = dt.id
JOIN github_installations gi
  ON d.installation_id = gi.installation_id
WHERE
  p.user_id = $1
  AND p.name ILIKE $2
  AND d.name ILIKE $3
  AND dh.is_active = true
GROUP BY
  d.id,
  dh.id,
  dt.id,
  gi.id,
  c.name;

-- name: GetDeploymentHistoryCountByCommitID :one
SELECT
  count(*)
FROM deployment_history dh
JOIN deployment d
  ON d.id = dh.deployment_id
WHERE
  d.id = $1
  AND dh.commit_id = $2;

-- name: GetTechstackByTechstackId :one
SELECT
  *
FROM deployment_techstack dt
WHERE
  id = $1;

-- name: GetTechstackName :many
SELECT DISTINCT ON (LOWER(name))
  name
FROM deployment_techstack;

-- name: GetTechstackVersionByName :many
SELECT
  dt.id AS id,
  dt.VERSION AS version
FROM deployment_techstack dt
WHERE
  LOWER(dt."name") = $1
ORDER BY
  split_part(version, '.', 1)::int DESC,
  split_part(version, '.', 2)::int DESC,
  split_part(version, '.', 3)::int DESC;

-- name: SetActiveDeploymentHistoryNonActiveByDeploymentId :exec
UPDATE deployment_history dh
SET is_active = false
FROM deployment d
JOIN project p
  ON p.id = d.project_id
WHERE
  d.id = dh.deployment_id
  AND p.user_id = $1
  AND dh.deployment_id = $2
  AND dh.is_active = true;

-- name: SetNonActiveDeploymentHistoryActiveByDeploymentHistoryId :exec
UPDATE deployment_history dh
SET is_active = true
FROM deployment d
JOIN project p
  ON p.id = d.project_id
WHERE
  d.id = dh.deployment_id
  AND p.user_id = $1
  AND dh.id = $2
  AND dh.is_active = false;

-- name: UpsertDeployment :one
INSERT INTO deployment (
    name,
    git_remote_url,
    project_id,
    installation_id,
    repository_id
)
SELECT
    $1,
    $2,
    p.id,
    $3,
    $4
FROM project p
WHERE p.name = $5
  AND p.user_id = $6
ON CONFLICT (name, project_id)
DO UPDATE
SET name = deployment.name
RETURNING id;

-- name: CreateDeploymentHistory :one
INSERT INTO public.deployment_history (deployment_id, branch, commit_id, commit_msg, version, deployment_techstack_id, build_command, build_folder, run_command, main_file_path)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
RETURNING id;

-- name: DeleteDeploymentByDeploymentId :exec
DELETE FROM
  deployment d
USING project p
WHERE
  p.id = d.project_id
  AND p.user_id = $1
  AND d.id = $2;

-- name: DeleteDeploymentByDeploymentName :exec
DELETE FROM
  deployment d
USING project p
WHERE
  p.user_id = $1
  AND p.name = $2
  AND d.name = $3;

-- name: UpdateDeploymentDescriptionByDeploymentId :exec
UPDATE deployment d
SET description = $3
FROM project p
WHERE
  p.id = d.project_id
  AND p.user_id = $1
  AND d.id = $2;
