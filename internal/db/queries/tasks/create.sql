WITH created AS (
    INSERT INTO tasks (project_id, title, description, start_time, end_time)
    VALUES ($1, $2, $3, $4, $5)
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
    c.id,
    c.project_id,
    c.title,
    c.description,
    c.completed,
    c.start_time,
    c.end_time,
    c.created_at,
    c.updated_at,
    p.name AS project_name
FROM created c
LEFT JOIN projects p ON p.id = c.project_id;
