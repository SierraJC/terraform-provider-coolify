package expand

import "github.com/hashicorp/terraform-plugin-framework/types"

func String(value types.String) *string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	return value.ValueStringPointer()
}

// StringOrNil returns nil if the string is null, unknown, or empty.
// This is useful for fields that should not be sent to the API if empty.
func StringOrNil(value types.String) *string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	str := value.ValueString()
	if str == "" {
		return nil
	}

	return &str
}

func RequiredString(value types.String) string {
	if value.IsNull() || value.IsUnknown() {
		return ""
	}

	return value.ValueString()
}
