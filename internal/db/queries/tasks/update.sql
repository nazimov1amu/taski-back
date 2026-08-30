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
RETURNING *;
