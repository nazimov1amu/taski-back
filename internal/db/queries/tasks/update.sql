WITH updated AS (
    UPDATE tasks
    SET project_id = $2,
        title = $3,
        description = $4,
        completed = $5,
        start_time = $6,
        end_time = $7,
        updated_at = now()
    WHERE id = $1
    RETURNING
        id,
        project_id,
        title,
        description,
        completed,
        start_time,
        end_time,
        created_at,
        updated_at
)
SELECT
    u.id,
    u.project_id,
    u.title,
    u.description,
    u.completed,
    u.start_time,
    u.end_time,
    u.created_at,
    u.updated_at,
    p.name AS project_name
FROM updated u
LEFT JOIN projects p ON p.id = u.project_id;
