-- name: GetActiveDeploymentHistoryContainerByDeploymentId :one
SELECT
  c.*
FROM container c
JOIN deployment_history dh
  ON dh.id = c.deployment_history_id
WHERE
  dh.deployment_id = $1;

-- name: CreateContainer :exec
INSERT INTO container (name,deployment_history_id)
VALUES ($1,$2);
