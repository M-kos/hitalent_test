SELECT id, department_id, full_name, position,
       hired_at, created_at
FROM employee
WHERE department_id IN ?
ORDER BY full_name
