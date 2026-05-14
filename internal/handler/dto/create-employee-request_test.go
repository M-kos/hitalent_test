package dto

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateEmployeeRequest_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		req   CreateEmployeeRequest
		valid bool
	}{
		{
			name:  "valid request",
			req:   CreateEmployeeRequest{FullName: "Aaa Bbb", Position: "Developer"},
			valid: true,
		},
		{
			name:  "full_name too short",
			req:   CreateEmployeeRequest{FullName: "", Position: "Developer"},
			valid: false,
		},
		{
			name:  "full_name too long",
			req:   CreateEmployeeRequest{FullName: string(make([]byte, 201)), Position: "Developer"},
			valid: false,
		},
		{
			name:  "position too short",
			req:   CreateEmployeeRequest{FullName: "Aaa Bbb", Position: ""},
			valid: false,
		},
		{
			name:  "position too long",
			req:   CreateEmployeeRequest{FullName: "Aaa Bbb", Position: string(make([]byte, 201))},
			valid: false,
		},
		{
			name:  "full_name with spaces should be trimmed",
			req:   CreateEmployeeRequest{FullName: "  Aaa Bbb  ", Position: "Developer"},
			valid: true,
		},
		{
			name:  "valid hired_at format",
			req:   CreateEmployeeRequest{FullName: "Aaa Bbb", Position: "Developer", HiredAt: strPtr("2026-01-01T00:00:00Z")},
			valid: true,
		},
		{
			name:  "invalid hired_at format",
			req:   CreateEmployeeRequest{FullName: "Aaa Bbb", Position: "Developer", HiredAt: strPtr("invalid-date")},
			valid: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.req.Validate()
			if tt.valid {
				assert.NoError(t, err)
				assert.Equal(t, strings.TrimSpace(tt.req.FullName), tt.req.FullName, "full_name should be trimmed")
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestCreateEmployeeRequest_ToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		req             CreateEmployeeRequest
		departmentID    int
		expectError     bool
		expectedHired   bool
		expectedHiredAt time.Time
	}{
		{
			name:         "converts to domain with valid data",
			req:          CreateEmployeeRequest{FullName: "Aaa Bbb", Position: "Developer"},
			departmentID: 5,
			expectError:  false,
		},
		{
			name:            "converts to domain with hired_at",
			req:             CreateEmployeeRequest{FullName: "Aaa Bbb", Position: "Developer", HiredAt: strPtr("2026-01-01T00:00:00Z")},
			departmentID:    5,
			expectError:     false,
			expectedHired:   true,
			expectedHiredAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "returns error on invalid hired_at",
			req:          CreateEmployeeRequest{FullName: "Aaa Bbb", Position: "Developer", HiredAt: strPtr("invalid-date")},
			departmentID: 5,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dep, err := tt.req.ToDomain(tt.departmentID)

			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.req.FullName, dep.FullName)
			assert.Equal(t, tt.req.Position, dep.Position)
			assert.Equal(t, tt.departmentID, *dep.DepartmentID)
			if tt.expectedHired {
				assert.NotNil(t, dep.HiredAt)
				assert.WithinDuration(t, tt.expectedHiredAt, *dep.HiredAt, time.Second)
			} else {
				assert.Nil(t, dep.HiredAt)
			}
			assert.WithinDuration(t, time.Now(), dep.CreatedAt, time.Second)
		})
	}
}

func strPtr(s string) *string {
	return &s
}
