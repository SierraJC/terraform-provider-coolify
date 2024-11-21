package provider

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type filterBlockModel struct {
	Name   types.String `tfsdk:"name"`
	Values types.List   `tfsdk:"values"`
}

// createDatasourceFilter creates a filter block for a datasource schema.
func createDatasourceFilter(allowedFields []string) schema.Block {
	return schema.ListNestedBlock{
		MarkdownDescription: "Filter results by values",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					MarkdownDescription: fmt.Sprintf("Name of the field to filter on. Valid names are `%s`", strings.Join(allowedFields, "`, `")),
					Required:            true,
					Validators: []validator.String{
						stringvalidator.OneOf(allowedFields...),
					},
				},
				"values": schema.ListAttribute{
					Required:            true,
					MarkdownDescription: "List of values to match against - if any value matches, the filter is satisfied (**OR** operation). Non-string values will be converted to strings if possible, ie `true` -> `\"true\"`",
					ElementType:         types.StringType,
				},
			},
		},
	}
}

// filterOnAttributes filters a map of attributes based on a list of filter blocks.
func filterOnAttributes(attributes map[string]attr.Value, filters []filterBlockModel) bool {
	if len(filters) == 0 {
		return true
	}

	for _, filter := range filters {
		if attr, ok := attributes[filter.Name.ValueString()]; ok {
			attrValueString, err := attributeValueToString(attr)
			if err != nil {
				return false
			}

			filterValues := []string{}
			filter.Values.ElementsAs(context.Background(), &filterValues, false)
			if !slices.Contains(filterValues, attrValueString) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

// attributeValueToString converts any supported attribute value to its string representation.
func attributeValueToString(value attr.Value) (string, error) {
	switch v := value.(type) {
	case types.String:
		return v.ValueString(), nil
	case types.Bool:
		return fmt.Sprintf("%t", v.ValueBool()), nil
	case types.Int64:
		return fmt.Sprintf("%d", v.ValueInt64()), nil
	case types.Int32:
		return fmt.Sprintf("%d", v.ValueInt32()), nil
	case types.Float64:
		return fmt.Sprintf("%f", v.ValueFloat64()), nil
	case types.Float32:
		return fmt.Sprintf("%f", v.ValueFloat32()), nil
	default:
		return "", fmt.Errorf("unsupported attribute type: %T", value)
	}
}
