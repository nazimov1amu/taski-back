SELECT *
FROM projects
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;
