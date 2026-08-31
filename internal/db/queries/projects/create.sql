INSERT INTO projects (name, description)
VALUES ($1, $2)
RETURNING id, name, description;
