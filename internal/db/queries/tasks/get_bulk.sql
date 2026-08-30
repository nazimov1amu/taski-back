SELECT *
FROM tasks
WHERE planned_at = $1::date;
