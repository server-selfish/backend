-- name: GetActiveDeploymentHistoryContainerByDeploymentId :one
SELECT
  c.*
FROM container c
JOIN deployment_history dh
  ON dh.id = c.deployment_history_id
JOIN deployment d
  ON dh.deployment_id = d.id
JOIN project p
  ON d.project_id = p.id
WHERE
  p.user_id = $1
  AND dh.deployment_id = $2;

-- name: CreateContainer :exec
INSERT INTO container (name,deployment_history_id)
VALUES ($1,$2);
