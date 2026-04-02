-- name: GetAllProjects :many
SELECT
  *
FROM public.project p;

-- name: GetProjectById :one
SELECT
  *
FROM
  public.project p
WHERE
  p.id = $1;

-- name: CreateProject :exec
INSERT INTO public.project (name,description)
VALUES ($1,$2);

-- name: UpdateProjectById :exec
UPDATE public.project
SET name = $1,
    description = $2,
    updated_at = now()
WHERE id = $3;

-- name: DeleteProjectById :exec
DELETE FROM public.project
where id = $1;
