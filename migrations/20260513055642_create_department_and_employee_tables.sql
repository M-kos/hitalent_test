-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS department
(
	id BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL UNIQUE CHECK (name <> ''),
	parent_id BIGINT NULL REFERENCES department (id) ON DELETE SET NULL CHECK ( parent_id IS NULL OR parent_id != id),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

	CONSTRAINT department_parent_name_unique
		UNIQUE NULLS NOT DISTINCT(parent_id, name)
);

CREATE TABLE IF NOT EXISTS employee
(
	id BIGSERIAL PRIMARY KEY,
	department_id BIGINT NULL REFERENCES department (id) ON DELETE CASCADE,
	full_name TEXT NOT NULL,
	position TEXT NOT NULL,
	hired_at TIMESTAMPTZ NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_department_parent_id
	ON department(parent_id);

CREATE INDEX idx_employee_department_id
	ON employee(department_id);

-- =========================================================
-- DEPARTMENTS
-- =========================================================

INSERT INTO department (id, name, parent_id)
VALUES
	-- root without children
	(1, 'Board', NULL),

	-- root with depth = 6
	(2, 'Engineering', NULL),

	-- root with depth = 2
	(3, 'Operations', NULL);

-- ---------------------------------------------------------
-- Engineering tree (depth 6)
-- 2 -> 4 -> 5 -> 6 -> 7 -> 8 -> 9
-- ---------------------------------------------------------

INSERT INTO department (id, name, parent_id)
VALUES
	(4, 'Platform', 2),
	(5, 'Backend', 4),
	(6, 'Data', 5),
	(7, 'ML', 6),
	(8, 'Infrastructure', 7),
	(9, 'Observability', 8);

-- ---------------------------------------------------------
-- Operations tree (depth 2)
-- 3 -> 10 -> 11
-- ---------------------------------------------------------

INSERT INTO department (id, name, parent_id)
VALUES
	(10, 'Support', 3),
	(11, 'Regional Support', 10);

-- ---------------------------------------------------------
-- Additional branches
-- ---------------------------------------------------------

INSERT INTO department (id, name, parent_id)
VALUES
	(12, 'QA', 4),
	(13, 'Mobile', 4),
	(14, 'Security', 2),
	(15, 'Finance', 3);

-- =========================================================
-- EMPLOYEES
-- 1-2 employees per department
-- =========================================================

INSERT INTO employee (department_id, full_name, position, hired_at)
VALUES
	(1, 'Alice Johnson', 'Chairwoman', '2018-01-10'),
	(1, 'Robert Miles', 'Executive Assistant', '2020-03-15'),

	(2, 'John Carter', 'VP Engineering', '2017-06-01'),
	(2, 'Emma Brown', 'Engineering Coordinator', '2021-09-11'),

	(3, 'Olivia Smith', 'COO', '2019-04-20'),

	(4, 'Daniel White', 'Platform Lead', '2020-07-07'),
	(4, 'Sophia Green', 'Senior Platform Engineer', '2022-02-01'),

	(5, 'Michael Black', 'Backend Lead', '2018-11-12'),
	(5, 'Lucas Gray', 'Backend Engineer', '2023-01-16'),

	(6, 'Ethan Walker', 'Data Engineer', '2021-05-30'),

	(7, 'Mia Hall', 'ML Engineer', '2022-08-14'),
	(7, 'James Young', 'ML Researcher', '2023-04-09'),

	(8, 'Benjamin King', 'Infrastructure Engineer', '2020-12-12'),

	(9, 'Charlotte Scott', 'SRE', '2021-10-01'),
	(9, 'Henry Adams', 'Monitoring Engineer', '2024-01-01'),

	(10, 'Amelia Baker', 'Support Manager', '2019-09-09'),

	(11, 'Jack Turner', 'Regional Specialist', '2022-06-18'),
	(11, 'Lily Parker', 'Support Engineer', '2023-11-01'),

	(12, 'David Evans', 'QA Lead', '2020-02-22'),

	(13, 'Grace Hill', 'Mobile Engineer', '2021-03-03'),
	(13, 'Samuel Lee', 'iOS Engineer', '2024-02-14'),

	(14, 'Victoria Wright', 'Security Engineer', '2018-08-08'),

	(15, 'Andrew Harris', 'Financial Analyst', '2020-01-20'),
	(15, 'Ella Nelson', 'Accountant', '2022-05-12');

-- =========================================================
-- reset sequences after manual ids
-- =========================================================

SELECT setval(
		       'department_id_seq',
		       (SELECT MAX(id) FROM department)
       );

SELECT setval(
		       'employee_id_seq',
		       (SELECT MAX(id) FROM employee)
       );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS employee;
DROP TABLE IF EXISTS department;
-- +goose StatementEnd
