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
