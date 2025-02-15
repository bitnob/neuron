package utils

import (
	"testing"
)

type TestStruct struct {
	Name  string
	Email string
	Age   int
}

func TestValidator_Validate(t *testing.T) {
	validator := NewValidator()
	tests := []struct {
		name    string
		data    TestStruct
		rules   []ValidationRule
		wantErr bool
	}{
		{
			name: "valid data",
			data: TestStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   25,
			},
			rules: []ValidationRule{
				{
					Field:   "Name",
					Rule:    "required",
					Message: "Name is required",
				},
				{
					Field:   "Email",
					Rule:    "email",
					Message: "Invalid email format",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			data: TestStruct{
				Name:  "John Doe",
				Email: "invalid-email",
				Age:   25,
			},
			rules: []ValidationRule{
				{
					Field:   "Email",
					Rule:    "email",
					Message: "Invalid email format",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := validator.Validate(tt.data, tt.rules)
			if isValid == tt.wantErr {
				t.Errorf("Validate() = %v, want %v", isValid, !tt.wantErr)
			}
		})
	}
}

func BenchmarkValidator_Validate(b *testing.B) {
	validator := NewValidator()
	data := TestStruct{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}
	rules := []ValidationRule{
		{
			Field:   "Name",
			Rule:    "required",
			Message: "Name is required",
		},
		{
			Field:   "Email",
			Rule:    "email",
			Message: "Invalid email format",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.Validate(data, rules)
	}
}
