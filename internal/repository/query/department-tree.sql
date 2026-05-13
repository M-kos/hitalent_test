WITH RECURSIVE department_tree AS (
    SELECT
        id,
        name,
        parent_id,
        0 AS depth
    FROM department
    WHERE id = ?

    UNION ALL

    SELECT
        d.id,
        d.name,
        d.parent_id,
        dt.depth + 1
    FROM department d
    INNER JOIN department_tree dt
        ON d.parent_id = dt.id
    WHERE dt.depth < ?
)
SELECT *
FROM department_tree;
