WITH created AS (
    INSERT INTO tasks (project_id, parent_id, title, description, planned_at, start_time, end_time)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
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
    c.id,
    c.project_id,
    c.parent_id,
    c.title,
    c.description,
    c.completed,
    c.planned_at,
    c.start_time,
    c.end_time,
    c.created_at,
    c.updated_at,
    p.name AS project_name
FROM created c
LEFT JOIN projects p ON p.id = c.project_id;
