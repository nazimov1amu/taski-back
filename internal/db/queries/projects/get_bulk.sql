SELECT id, name, description
FROM projects
WHERE name ILIKE '%' || $1 || '%' 
    OR description ILIKE '%' || $1 || '%' 
    OR $1 IS NULL
ORDER BY name DESC;
