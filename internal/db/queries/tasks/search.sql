SELECT *
FROM tasks
WHERE title ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT 10;