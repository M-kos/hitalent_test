INSERT INTO department (
	name, parent_id
) VALUES (?, ?)
ON CONFLICT(parent_id, name) DO NOTHING 
RETURNING *
