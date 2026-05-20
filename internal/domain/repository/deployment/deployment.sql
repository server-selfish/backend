-- name: GetDeploymentsByProjectId :many
SELECT
  CAST (d.id AS VARCHAR) AS deployment_id,
  d.name AS deployment_name,
  d.git_remote_url as git_remote_url,
  CAST (dh.branch AS VARCHAR) AS branch,
  CAST (dh.commit_id AS VARCHAR) AS commit_id,
  CAST (dh.commit_msg AS VARCHAR) AS commit_message,
  CAST (dh.version AS VARCHAR) AS deployment_version,
  dh.external_port AS port,
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
WHERE
  p.user_id = $1
  AND d.project_id = $2
  AND dh.is_active = true;

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

-- name: GetDeploymentHistoryByDeploymentId :many
SELECT
  dh.id,
  dh.branch,
  dh.commit_id,
  dh.commit_msg AS commit_message,
  dh.version AS deployment_version,
  dh.external_port AS port,
  dh.created_at,
  dh.updated_at
FROM deployment_history dh
JOIN deployment d
  ON dh.deployment_id = d.id
JOIN project p
  ON d.project_id = p.id
WHERE
  p.user_id = $1 AND
  dh.deployment_id = $2 AND
  dh.is_active = false
ORDER BY
  COALESCE(dh.updated_at, dh.created_at) DESC;

-- name: GetActiveDeploymentHistoryByDeploymentId :one
SELECT
  dh.id AS deployment_history_id,
  dh.branch AS branch,
  dh.commit_id AS commit_id,
  dh.commit_msg AS commit_message,
  dh.version AS deployment_version,
  dh.external_port AS port,
  dh.build_command AS build_command,
  dt.id AS techstack_id,
  dt.name AS techstack_name,
  dt.version AS techstack_version
FROM deployment_history dh
JOIN deployment d
  ON dh.deployment_id = d.id
JOIN project p
  ON d.project_id = p.id
JOIN deployment_techstack dt
  ON dh.deployment_techstack_id = dt.id
WHERE
  p.user_id = $1
  AND dh.deployment_id = $2
  AND dh.is_active = true;

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
  LOWER(dt."name") = $1;

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

-- name: CreateDeployment :exec
INSERT INTO public.deployment (name,git_remote_url,project_id,installation_id)
VALUES ($1,$2,$3,$4);

-- name: CreateDeploymentHistory :one
INSERT INTO public.deployment_history (deployment_id, branch, commit_id, commit_msg, version, external_port, deployment_techstack_id, build_command, build_folder, run_command)
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
