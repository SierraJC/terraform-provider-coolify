package provider

import (
	"fmt"

	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func makeResourceAttributeRequired(
	attributes map[string]resource_schema.Attribute,
	attrName string,
) error {
	attr, ok := attributes[attrName]
	if !ok {
		return fmt.Errorf("attribute %s not found", attrName)
	}

	switch typedAttr := attr.(type) {
	case resource_schema.StringAttribute:
		typedAttr.Required = true
		typedAttr.Optional = false
		typedAttr.Computed = false
		attributes[attrName] = typedAttr
	case resource_schema.BoolAttribute:
		typedAttr.Required = true
		typedAttr.Optional = false
		typedAttr.Computed = false
		attributes[attrName] = typedAttr
	case resource_schema.Int64Attribute:
		typedAttr.Required = true
		typedAttr.Optional = false
		typedAttr.Computed = false
		attributes[attrName] = typedAttr
	default:
		return fmt.Errorf("unsupported attribute type for %s", attrName)
	}

	return nil
}

func makeResourceAttributeSensitive(
	attributes map[string]resource_schema.Attribute,
	attrName string,
) error {
	attr, ok := attributes[attrName]
	if !ok {
		return fmt.Errorf("attribute %s not found", attrName)
	}

	switch typedAttr := attr.(type) {
	case resource_schema.StringAttribute:
		typedAttr.Sensitive = true
		attributes[attrName] = typedAttr
	default:
		return fmt.Errorf("unsupported attribute type for %s", attrName)
	}

	return nil
}

func makeDataSourceAttributeSensitive(
	attributes map[string]datasource_schema.Attribute,
	attrName string,
) error {
	attr, ok := attributes[attrName]
	if !ok {
		return fmt.Errorf("attribute %s not found", attrName)
	}

	switch typedAttr := attr.(type) {
	case datasource_schema.StringAttribute:
		typedAttr.Sensitive = true
		attributes[attrName] = typedAttr
	default:
		return fmt.Errorf("unsupported attribute type for %s", attrName)
	}

	return nil
}
