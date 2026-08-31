SELECT
    t.id,
    t.project_id,
    t.title,
    t.description,
    t.completed,
    t.start_time,
    t.end_time,
    t.created_at,
    t.updated_at,
    p.name AS project_name
FROM tasks t
LEFT JOIN projects p ON p.id = t.project_id
WHERE (
    $1::date IS NULL
    OR (
        t.start_time <= (($1::date + 1)::timestamptz)
        AND t.end_time   > ($1::date::timestamptz)
    )
)
AND ($2::uuid IS NULL OR t.project_id = $2::uuid)
AND (
    $3::text IS NULL
    OR t.title ILIKE '%' || $3 || '%'
    OR t.description ILIKE '%' || $3 || '%'
)
ORDER BY t.created_at DESC;
