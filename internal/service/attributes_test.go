package service

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stretchr/testify/assert"
)

func TestMakeResourceAttributeRequired(t *testing.T) {
	tests := []struct {
		name        string
		attributes  map[string]resource_schema.Attribute
		attrName    string
		expectedErr string
	}{
		{
			name: "attribute not found",
			attributes: map[string]resource_schema.Attribute{
				"existing_attr": resource_schema.StringAttribute{},
			},
			attrName:    "missing_attr",
			expectedErr: "attribute missing_attr not found",
		},
		{
			name: "unsupported attribute type",
			attributes: map[string]resource_schema.Attribute{
				"unsupported_attr": resource_schema.DynamicAttribute{},
			},
			attrName:    "unsupported_attr",
			expectedErr: "unsupported attribute type for unsupported_attr",
		},
		{
			name: "string attribute",
			attributes: map[string]resource_schema.Attribute{
				"string_attr": resource_schema.StringAttribute{},
			},
			attrName:    "string_attr",
			expectedErr: "",
		},
		{
			name: "bool attribute",
			attributes: map[string]resource_schema.Attribute{
				"bool_attr": resource_schema.BoolAttribute{},
			},
			attrName:    "bool_attr",
			expectedErr: "",
		},
		{
			name: "int64 attribute",
			attributes: map[string]resource_schema.Attribute{
				"int64_attr": resource_schema.Int64Attribute{},
			},
			attrName:    "int64_attr",
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := makeResourceAttributeRequired(tt.attributes, tt.attrName)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				attr := tt.attributes[tt.attrName]
				switch typedAttr := attr.(type) {
				case resource_schema.StringAttribute:
					assert.True(t, typedAttr.Required)
					assert.False(t, typedAttr.Optional)
					assert.False(t, typedAttr.Computed)
				case resource_schema.BoolAttribute:
					assert.True(t, typedAttr.Required)
					assert.False(t, typedAttr.Optional)
					assert.False(t, typedAttr.Computed)
				case resource_schema.Int64Attribute:
					assert.True(t, typedAttr.Required)
					assert.False(t, typedAttr.Optional)
					assert.False(t, typedAttr.Computed)
				}
			}
		})
	}
}

func TestMakeResourceAttributeSensitive(t *testing.T) {
	tests := []struct {
		name        string
		attributes  map[string]resource_schema.Attribute
		attrName    string
		expectedErr string
	}{
		{
			name: "attribute not found",
			attributes: map[string]resource_schema.Attribute{
				"existing_attr": resource_schema.StringAttribute{},
			},
			attrName:    "missing_attr",
			expectedErr: "attribute missing_attr not found",
		},
		{
			name: "unsupported attribute type",
			attributes: map[string]resource_schema.Attribute{
				"unsupported_attr": resource_schema.DynamicAttribute{},
			},
			attrName:    "unsupported_attr",
			expectedErr: "unsupported attribute type for unsupported_attr",
		},
		{
			name: "string attribute",
			attributes: map[string]resource_schema.Attribute{
				"string_attr": resource_schema.StringAttribute{},
			},
			attrName:    "string_attr",
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := makeResourceAttributeSensitive(tt.attributes, tt.attrName)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				attr := tt.attributes[tt.attrName]
				switch typedAttr := attr.(type) {
				case datasource_schema.StringAttribute:
					assert.True(t, typedAttr.Sensitive)
				}
			}
		})
	}
}

func TestMakeDataSourceAttributeSensitive(t *testing.T) {
	tests := []struct {
		name        string
		attributes  map[string]datasource_schema.Attribute
		attrName    string
		expectedErr string
	}{
		{
			name: "attribute not found",
			attributes: map[string]datasource_schema.Attribute{
				"existing_attr": datasource_schema.StringAttribute{},
			},
			attrName:    "missing_attr",
			expectedErr: "attribute missing_attr not found",
		},
		{
			name: "unsupported attribute type",
			attributes: map[string]datasource_schema.Attribute{
				"unsupported_attr": datasource_schema.DynamicAttribute{},
			},
			attrName:    "unsupported_attr",
			expectedErr: "unsupported attribute type for unsupported_attr",
		},
		{
			name: "string attribute",
			attributes: map[string]datasource_schema.Attribute{
				"string_attr": datasource_schema.StringAttribute{},
			},
			attrName:    "string_attr",
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := makeDataSourceAttributeSensitive(tt.attributes, tt.attrName)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				attr := tt.attributes[tt.attrName]
				switch typedAttr := attr.(type) {
				case datasource_schema.StringAttribute:
					assert.True(t, typedAttr.Sensitive)
				}
			}
		})
	}
}

func TestSetResourceDefaultValue(t *testing.T) {
	tests := []struct {
		name         string
		attributes   map[string]resource_schema.Attribute
		attrName     string
		defaultValue interface{}
		expectedErr  string
	}{
		{
			name: "attribute not found",
			attributes: map[string]resource_schema.Attribute{
				"existing_attr": resource_schema.StringAttribute{},
			},
			attrName:     "missing_attr",
			defaultValue: "default",
			expectedErr:  "attribute missing_attr not found",
		},
		{
			name: "unsupported attribute type",
			attributes: map[string]resource_schema.Attribute{
				"unsupported_attr": resource_schema.DynamicAttribute{},
			},
			attrName:     "unsupported_attr",
			defaultValue: nil,
			expectedErr:  "unsupported attribute type for unsupported_attr",
		},
		{
			name: "string attribute with default",
			attributes: map[string]resource_schema.Attribute{
				"string_attr": resource_schema.StringAttribute{},
			},
			attrName:     "string_attr",
			defaultValue: "default",
			expectedErr:  "",
		},
		{
			name: "bool attribute with default",
			attributes: map[string]resource_schema.Attribute{
				"bool_attr": resource_schema.BoolAttribute{},
			},
			attrName:     "bool_attr",
			defaultValue: true,
			expectedErr:  "",
		},
		{
			name: "int64 attribute with default",
			attributes: map[string]resource_schema.Attribute{
				"int64_attr": resource_schema.Int64Attribute{},
			},
			attrName:     "int64_attr",
			defaultValue: int64(42),
			expectedErr:  "",
		},
		{
			name: "string attribute with wrong default type",
			attributes: map[string]resource_schema.Attribute{
				"string_attr": resource_schema.StringAttribute{},
			},
			attrName:     "string_attr",
			defaultValue: 123, // wrong type
			expectedErr:  "",
		},
		{
			name: "string attribute with nil default",
			attributes: map[string]resource_schema.Attribute{
				"string_attr": resource_schema.StringAttribute{},
			},
			attrName:     "string_attr",
			defaultValue: nil,
			expectedErr:  "",
		},
		{
			name: "bool attribute with nil default",
			attributes: map[string]resource_schema.Attribute{
				"bool_attr": resource_schema.BoolAttribute{},
			},
			attrName:     "bool_attr",
			defaultValue: nil,
			expectedErr:  "",
		},
		{
			name: "int64 attribute with nil default",
			attributes: map[string]resource_schema.Attribute{
				"int64_attr": resource_schema.Int64Attribute{},
			},
			attrName:     "int64_attr",
			defaultValue: nil,
			expectedErr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setResourceDefaultValue(tt.attributes, tt.attrName, tt.defaultValue)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				attr := tt.attributes[tt.attrName]
				switch typedAttr := attr.(type) {
				case resource_schema.StringAttribute:
					if strVal, ok := tt.defaultValue.(string); ok {
						assert.Equal(t, stringdefault.StaticString(strVal), typedAttr.Default)
					} else {
						assert.Nil(t, typedAttr.Default)
					}
				case resource_schema.BoolAttribute:
					if boolVal, ok := tt.defaultValue.(bool); ok {
						assert.Equal(t, booldefault.StaticBool(boolVal), typedAttr.Default)
					} else {
						assert.Nil(t, typedAttr.Default)
					}
				case resource_schema.Int64Attribute:
					if intVal, ok := tt.defaultValue.(int64); ok {
						assert.Equal(t, int64default.StaticInt64(intVal), typedAttr.Default)
					} else {
						assert.Nil(t, typedAttr.Default)
					}
				default:
					t.Errorf("Unexpected attribute type")
				}
			}
		})
	}
}

func TestMakeResourceAttributeNonEmpty(t *testing.T) {
	tests := []struct {
		name        string
		attributes  map[string]resource_schema.Attribute
		attrName    string
		expectedErr string
	}{
		{
			name: "attribute not found",
			attributes: map[string]resource_schema.Attribute{
				"existing_attr": resource_schema.StringAttribute{},
			},
			attrName:    "missing_attr",
			expectedErr: "attribute missing_attr not found",
		},
		{
			name: "unsupported attribute type",
			attributes: map[string]resource_schema.Attribute{
				"unsupported_attr": resource_schema.BoolAttribute{},
			},
			attrName:    "unsupported_attr",
			expectedErr: "non-empty validation only supported for string attributes: unsupported_attr",
		},
		{
			name: "string attribute",
			attributes: map[string]resource_schema.Attribute{
				"string_attr": resource_schema.StringAttribute{},
			},
			attrName:    "string_attr",
			expectedErr: "",
		},
		{
			name: "string attribute with existing validators",
			attributes: map[string]resource_schema.Attribute{
				"string_attr": resource_schema.StringAttribute{
					Validators: []validator.String{
						stringvalidator.LengthBetween(1, 255),
					},
				},
			},
			attrName:    "string_attr",
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := makeResourceAttributeNonEmpty(tt.attributes, tt.attrName)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				attr := tt.attributes[tt.attrName]
				switch typedAttr := attr.(type) {
				case resource_schema.StringAttribute:
					// Check that the validator was added
					assert.NotEmpty(t, typedAttr.Validators)
					// The last validator should be our LengthAtLeast validator
					lastValidator := typedAttr.Validators[len(typedAttr.Validators)-1]
					assert.IsType(t, stringvalidator.LengthAtLeast(1), lastValidator)
				}
			}
		})
	}
}
