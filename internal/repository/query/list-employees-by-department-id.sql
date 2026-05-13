SELECT id, department_id, full_name, position,
       hired_at, created_at
FROM employee
WHERE department_id = ANY(?)
ORDER BY full_name
