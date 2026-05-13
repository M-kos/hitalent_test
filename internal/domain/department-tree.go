package domain

type DepartmentTree struct {
	Department         *Department
	ChildrenDepartment []*Department
	Employees          []*Employee
}
