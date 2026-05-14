WITH RECURSIVE department_tree AS (
    SELECT
        id,
        name,
        parent_id
    FROM department
    WHERE id = ?

    UNION ALL

    SELECT
        d.id,
        d.name,
        d.parent_id
    FROM department d
    INNER JOIN department_tree dt
        ON d.parent_id = dt.id
)
SELECT id
FROM department_tree;
