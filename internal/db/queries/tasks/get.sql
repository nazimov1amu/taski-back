SELECT
    t.id,
    t.project_id,
    t.title,
    t.description,
    t.completed,
    t.planned_at,
    t.start_time,
    t.end_time,
    t.created_at,
    t.updated_at,
    p.name AS project_name
FROM tasks t
LEFT JOIN projects p ON p.id = t.project_id
WHERE t.id = $1;
