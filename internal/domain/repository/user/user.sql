-- name: UpsertGithubUser :one
INSERT INTO users (
  id,
  provider,
  provider_user_id,
  username,
  email,
  avatar_url,
  created_at,
  updated_at
) VALUES (
  $1,
  'github',
  $2,
  $3,
  $4,
  $5,
  now(),
  now()
)
ON CONFLICT (provider, provider_user_id)
DO UPDATE SET
  username = EXCLUDED.username,
  email = EXCLUDED.email,
  avatar_url = EXCLUDED.avatar_url,
  updated_at = now()
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: CreateAuthSession :one
INSERT INTO auth_sessions (
  id,
  user_id,
  refresh_token_hash,
  expires_at,
  user_agent,
  ip_address,
  created_at,
  updated_at
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  now(),
  now()
)
RETURNING *;

-- name: GetAuthSessionByRefreshTokenHash :one
SELECT *
FROM auth_sessions
WHERE refresh_token_hash = $1
  AND revoked_at IS NULL
  AND expires_at > now();

-- name: RevokeAuthSessionByID :exec
UPDATE auth_sessions
SET revoked_at = now(),
    updated_at = now()
WHERE id = $1
  AND revoked_at IS NULL;

-- name: RotateAuthSessionToken :exec
UPDATE auth_sessions
SET refresh_token_hash = $2,
    expires_at = $3,
    updated_at = now()
WHERE id = $1
  AND revoked_at IS NULL;
