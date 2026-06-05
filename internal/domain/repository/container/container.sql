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

-- name: CreateContainer :one
INSERT INTO container (name,deployment_history_id)
VALUES ($1,$2)
RETURNING id;

-- name: GetContainerByName :one
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
  AND c.name = $2;

-- name: CreateContainerEnv :exec
INSERT INTO public.container_env (
  container_id,
  key,
  value
)
SELECT
  sqlc.arg(container_ids)::uuid,
  k.key,
  v.value
FROM unnest(sqlc.arg(keys)::text[])
  WITH ORDINALITY AS k(key, ord)
JOIN unnest(sqlc.arg(values)::text[])
  WITH ORDINALITY AS v(value, ord)
  USING (ord);

-- name: CreateContainerPort :exec
INSERT INTO public.container_port (
  container_id,
  external,
  internal,
  protocol
)
SELECT
  sqlc.arg(container_ids)::uuid,
  e.external,
  i.internal,
  p.protocol
FROM unnest(sqlc.arg(external)::integer[])
  WITH ORDINALITY AS e(external, ord)
JOIN unnest(sqlc.arg(internal)::integer[])
  WITH ORDINALITY AS i(internal, ord)
  USING (ord)
JOIN unnest(sqlc.arg(protocol)::varchar[])
  WITH ORDINALITY AS p(protocol, ord)
  USING (ord);
