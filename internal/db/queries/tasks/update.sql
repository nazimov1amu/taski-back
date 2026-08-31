WITH updated AS (
    UPDATE tasks
    SET project_id = $2,
        title = $3,
        description = $4,
        completed = $5,
        planned_at = $6,
        start_time = $7,
        end_time = $8,
        updated_at = now()
    WHERE id = $1
    RETURNING
        id,
        project_id,
        title,
        description,
        completed,
        planned_at,
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
    u.planned_at,
    u.start_time,
    u.end_time,
    u.created_at,
    u.updated_at,
    p.name AS project_name
FROM updated u
LEFT JOIN projects p ON p.id = u.project_id;
