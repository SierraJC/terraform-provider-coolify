package provider

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestOptionalString(t *testing.T) {
	input := "test"
	assert.Equal(t, types.StringNull(), optionalString(nil))
	assert.Equal(t, types.StringValue(input), optionalString(&input))
}

func TestOptionalInt64(t *testing.T) {
	input := 42
	assert.Equal(t, types.Int64Null(), optionalInt64(nil))
	assert.Equal(t, types.Int64Value(int64(input)), optionalInt64(&input))
}

func TestOptionalBool(t *testing.T) {
	input := true
	assert.Equal(t, types.BoolNull(), optionalBool(nil))
	assert.Equal(t, types.BoolValue(input), optionalBool(&input))
}

func TestOptionalStringListValue(t *testing.T) {
	tests := []struct {
		name     string
		input    *[]string
		expected types.List
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: types.ListNull(types.StringType),
		},
		{
			name:     "empty array input",
			input:    &[]string{},
			expected: types.ListValueMust(types.StringType, []attr.Value{}),
		},
		{
			name:  "populated array input",
			input: &[]string{"one", "two", "three"},
			expected: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("one"),
				types.StringValue("two"),
				types.StringValue("three"),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, optionalStringListValue(tt.input))
		})
	}
}

func TestOptionalTime(t *testing.T) {
	now := time.Now()
	formattedTime := now.Format(time.RFC3339Nano)

	tests := []struct {
		name     string
		input    *time.Time
		expected types.String
	}{
		{"nil input", nil, types.StringNull()},
		{"valid time input", &now, types.StringValue(formattedTime)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, optionalTime(tt.input))
		})
	}
}
