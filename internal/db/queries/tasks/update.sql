WITH updated AS (
    UPDATE tasks
    SET project_id = $2,
        parent_id = $3,
        title = $4,
        description = $5,
        completed = $6,
        planned_at = $7,
        start_time = $8,
        end_time = $9,
        updated_at = now()
    WHERE id = $1
    RETURNING
        id,
        project_id,
        parent_id,
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
    u.parent_id,
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
