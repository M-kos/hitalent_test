INSERT INTO department (
	name, parent_id
) VALUES (?, ?)
ON CONFLICT (name) DO NOTHING
RETURNING *
