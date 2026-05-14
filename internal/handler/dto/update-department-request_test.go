package dto

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateDepartmentRequest_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		req   UpdateDepartmentRequest
		valid bool
	}{
		{
			name:  "valid request with name",
			req:   UpdateDepartmentRequest{Name: "Backend"},
			valid: true,
		},
		{
			name:  "valid request with parent",
			req:   UpdateDepartmentRequest{ParentID: intPtr(1)},
			valid: true,
		},
		{
			name:  "valid empty request",
			req:   UpdateDepartmentRequest{},
			valid: true,
		},
		{
			name:  "name too long",
			req:   UpdateDepartmentRequest{Name: string(make([]byte, 201))},
			valid: false,
		},
		{
			name:  "name with spaces should be trimmed",
			req:   UpdateDepartmentRequest{Name: "  Backend2  "},
			valid: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.req.Validate()
			if tt.valid {
				assert.NoError(t, err)
				if tt.req.Name != "" {
					assert.Equal(t, strings.TrimSpace(tt.req.Name), tt.req.Name, "name should be trimmed")
				}
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestUpdateDepartmentRequest_ToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		req          UpdateDepartmentRequest
		id           int
		expectedName string
		expectedNil  bool
		expectedID   int
	}{
		{
			name:         "converts to domain with name and parent",
			req:          UpdateDepartmentRequest{Name: "Backend", ParentID: intPtr(1)},
			id:           5,
			expectedName: "Backend",
			expectedID:   1,
		},
		{
			name:         "converts to domain with only name",
			req:          UpdateDepartmentRequest{Name: "Frontend"},
			id:           5,
			expectedName: "Frontend",
			expectedNil:  true,
		},
		{
			name:         "converts to domain with only parent",
			req:          UpdateDepartmentRequest{ParentID: intPtr(1)},
			id:           5,
			expectedName: "",
			expectedID:   1,
		},
		{
			name:         "converts to domain with no changes",
			req:          UpdateDepartmentRequest{},
			id:           5,
			expectedName: "",
			expectedNil:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dep := tt.req.ToDomain(tt.id)
			assert.Equal(t, tt.id, dep.ID)
			assert.Equal(t, tt.expectedName, dep.Name)
			if tt.expectedNil {
				assert.Nil(t, dep.Parent)
			} else {
				assert.NotNil(t, dep.Parent)
				assert.Equal(t, tt.expectedID, dep.Parent.ID)
			}
		})
	}
}
