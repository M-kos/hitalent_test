package dto

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDepartmentRequest_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		req   CreateDepartmentRequest
		valid bool
	}{
		{
			name:  "valid request with parent",
			req:   CreateDepartmentRequest{Name: "Backend", ParentID: intPtr(1)},
			valid: true,
		},
		{
			name:  "valid request without parent",
			req:   CreateDepartmentRequest{Name: "Frontend"},
			valid: true,
		},
		{
			name:  "name too short",
			req:   CreateDepartmentRequest{Name: "", ParentID: intPtr(1)},
			valid: false,
		},
		{
			name:  "name too long",
			req:   CreateDepartmentRequest{Name: string(make([]byte, 201)), ParentID: intPtr(1)},
			valid: false,
		},
		{
			name:  "name with spaces should be trimmed",
			req:   CreateDepartmentRequest{Name: "  Backend2  "},
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
				assert.Equal(t, strings.TrimSpace(tt.req.Name), tt.req.Name, "name should be trimmed")
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestCreateDepartmentRequest_ToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		req         CreateDepartmentRequest
		expectedID  int
		expectedNil bool
	}{
		{
			name:       "converts to domain with parent",
			req:        CreateDepartmentRequest{Name: "Backend", ParentID: intPtr(1)},
			expectedID: 1,
		},
		{
			name:        "converts to domain without parent",
			req:         CreateDepartmentRequest{Name: "Frontend"},
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dep := tt.req.ToDomain()
			assert.Equal(t, tt.req.Name, dep.Name)
			if tt.expectedNil {
				assert.Nil(t, dep.Parent)
			} else {
				assert.NotNil(t, dep.Parent)
				assert.Equal(t, tt.expectedID, dep.Parent.ID)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
