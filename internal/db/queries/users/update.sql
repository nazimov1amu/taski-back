UPDATE users 
SET email = COALESCE($2, email), username = COALESCE($3, username), password = COALESCE($4, password)
WHERE id = $1
RETURNING id, username;
