WITH updated AS (
    UPDATE tasks
    SET project_id = COALESCE($2, project_id),
        title = COALESCE($3, title),
        description = COALESCE($4, description),
        completed = COALESCE($5, completed),
        planned_at = COALESCE($6, planned_at),
        start_time = COALESCE($7, start_time),
        end_time = COALESCE($8, end_time),
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
