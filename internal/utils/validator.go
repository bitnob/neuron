package utils

import (
	"reflect"
	"regexp"
)

// Validator provides validation utilities
type Validator struct {
	errors map[string]string
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		errors: make(map[string]string),
	}
}

// ValidationRule defines a validation rule
type ValidationRule struct {
	Field    string
	Rule     string
	Message  string
	Function func(interface{}) bool
}

// Validate performs validation on a struct
func (v *Validator) Validate(data interface{}, rules []ValidationRule) bool {
	v.errors = make(map[string]string)
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for _, rule := range rules {
		field := val.FieldByName(rule.Field)
		if !field.IsValid() {
			v.errors[rule.Field] = "Field not found"
			continue
		}

		value := field.Interface()

		// Built-in validation rules
		switch rule.Rule {
		case "required":
			if isEmpty(value) {
				v.errors[rule.Field] = rule.Message
			}
		case "email":
			if !isValidEmail(value.(string)) {
				v.errors[rule.Field] = rule.Message
			}
		case "min":
			if !isMinLength(value.(string), 6) {
				v.errors[rule.Field] = rule.Message
			}
		}

		// Custom validation function
		if rule.Function != nil && !rule.Function(value) {
			v.errors[rule.Field] = rule.Message
		}
	}

	return len(v.errors) == 0
}

// Common validation helpers
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return v == ""
	case int:
		return v == 0
	case bool:
		return !v
	case []interface{}:
		return len(v) == 0
	default:
		return false
	}
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}
