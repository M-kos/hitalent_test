-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS department
(
	id BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	parent_id BIGINT NULL REFERENCES department (id) ON DELETE SET NULL CHECK ( parent_id IS NULL OR parent_id != id),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS employee
(
	id BIGSERIAL PRIMARY KEY,
	department_id BIGINT NULL REFERENCES department (id) ON DELETE SET NULL,
	full_name TEXT NOT NULL,
	position TEXT NOT NULL,
	hired_at TIMESTAMPTZ NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_department_parent_id
	ON department(parent_id);

CREATE INDEX idx_employee_department_id
	ON employee(department_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS employee;
DROP TABLE IF EXISTS department;
-- +goose StatementEnd
