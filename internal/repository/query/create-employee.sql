INSERT INTO employee (
	department_id, full_name, position, hired_at
) VALUES (?, ?, ?, ?)
RETURNING *
