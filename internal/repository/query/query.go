package query

import _ "embed"

var (
	//go:embed create-department.sql
	CreateDepartment string
	//go:embed create-employee.sql
	CreateEmployee string
	//go:embed delete-departments.sql
	DeleteDepartments string
	//go:embed delete-employees-by-department-id.sql
	DeleteEmployeesByDepartmentId string
	//go:embed department-tree.sql
	DepartmentTree string
	//go:embed list-employees-by-department-id.sql
	ListEmployeesByDepartmentId string
	//go:embed update-department.sql
	UpdateDepartment string
	//go:embed update-department-for-employees.sql
	UpdateDepartmentForEmployees string
)
