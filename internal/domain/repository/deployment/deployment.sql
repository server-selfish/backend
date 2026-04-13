-- name: GetDeploymentsByProjectId :many
SELECT
  CAST (d.id AS VARCHAR) AS deployment_id,
  d.name AS deployment_name,
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
LEFT JOIN deployment_history dh
  ON d.id = dh.deployment_id
JOIN deployment_techstack dt
  ON dh.deployment_techstack_id = dt.id
LEFT JOIN container c
  ON c.deployment_history_id = dh.id
WHERE
  d.project_id = $1 AND
  dh.is_active = true;

-- name: GetDeploymentByDeploymentId :one
SELECT
  *
FROM deployment d
WHERE
  d.id = $1;

-- name: GetDeploymentHistoryByDeploymentId :many
SELECT
  id,
  branch,
  commit_id,
  commit_msg AS commit_message,
  version AS deployment_version,
  external_port AS port,
  created_at,
  updated_at
FROM deployment_history dh
WHERE
  dh.deployment_id = $1 AND
  dh.is_active = false
ORDER BY
  COALESCE(dh.updated_at, dh.created_at) DESC;

-- name: GetActiveDeploymentHistoryByDeploymentId :one
SELECT
  dh.id AS deployment_history_id,
  dh.git_remote_url AS git_remote_url,
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
JOIN deployment_techstack dt
  ON dh.deployment_techstack_id = dt.id
WHERE
  dh.deployment_id = $1 AND
  dh.is_active = true;

-- name: GetActiveDeploymentHistoryContainerByDeploymentId :one
SELECT
  c.*
FROM deployment_history dh
JOIN container c
  ON dh.id = c.deployment_history_id
WHERE
  dh.deployment_id = $1;

-- name: GetTechstackByTechstackId :one
SELECT
  *
FROM deployment_techstack dt
WHERE
  id = $1;

-- name: SetActiveDeploymentHistoryNonActiveByDeploymentId :exec
UPDATE deployment_history dh
SET is_active = false
WHERE
  dh.deployment_id = $1 AND
  dh.is_active = true;

-- name: SetNonActiveDeploymentHistoryActiveByDeploymentHistoryId :exec
UPDATE deployment_history dh
SET is_active = true
WHERE
  dh.id = $1 AND
  dh.is_active = false;

-- name: CreateDeployment :exec
INSERT INTO public.deployment (name,project_id)
VALUES ($1,$2);

-- name: CreateDeploymentHistory :exec
INSERT INTO public.deployment_history (deployment_id, git_remote_url, branch, commit_id, commit_msg, version, external_port, deployment_techstack_id, build_command, build_folder, run_command)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);

-- name: DeleteDeploymentByDeploymentId :exec
DELETE FROM
  deployment d
WHERE
  d.id = $1;
