INSERT INTO tasks (project_id, title, description, planned_at, start_time, end_time)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
