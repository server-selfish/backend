-- name: UpsertGithubInstallation :one
INSERT INTO github_installations (
  id,
  user_id,
  installation_id,
  account_login,
  account_id,
  target_type,
  updated_at
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  now()
)
ON CONFLICT (user_id, installation_id)
DO UPDATE SET
  account_login = EXCLUDED.account_login,
  account_id = EXCLUDED.account_id,
  target_type = EXCLUDED.target_type,
  updated_at = now()
RETURNING
  id,
  user_id,
  installation_id,
  account_login,
  account_id,
  target_type,
  created_at,
  updated_at;

-- name: ListGithubInstallationsByUserID :many
SELECT
  id,
  user_id,
  installation_id,
  account_login,
  account_id,
  target_type,
  created_at,
  updated_at
FROM github_installations
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetGithubInstallationByUserIDAndInstallationID :one
SELECT
  id,
  user_id,
  installation_id,
  account_login,
  account_id,
  target_type,
  created_at,
  updated_at
FROM github_installations
WHERE user_id = $1
  AND installation_id = $2
LIMIT 1;

-- name: DeleteGithubInstallationByUserIDAndInstallationID :exec
DELETE FROM github_installations
WHERE user_id = $1
  AND installation_id = $2;

-- name: DeleteGithubInstallationsByUserID :exec
DELETE FROM github_installations
WHERE user_id = $1;
