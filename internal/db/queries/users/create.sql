INSERT INTO users (email, username, password_hash) 
VALUES ($1, $2, $3) 
RETURNING id, username;