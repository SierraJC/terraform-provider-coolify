// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package datasource_private_keys

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func PrivateKeysDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"private_keys": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"fingerprint": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"is_git_related": schema.BoolAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"private_key": schema.StringAttribute{
							Computed: true,
						},
						"team_id": schema.Int64Attribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
							Computed: true,
						},
						"uuid": schema.StringAttribute{
							Computed: true,
						},
					},
					CustomType: PrivateKeysType{
						ObjectType: types.ObjectType{
							AttrTypes: PrivateKeysValue{}.AttributeTypes(ctx),
						},
					},
				},
				Computed: true,
			},
		},
	}
}

type PrivateKeysModel struct {
	PrivateKeys types.Set `tfsdk:"private_keys"`
}

var _ basetypes.ObjectTypable = PrivateKeysType{}

type PrivateKeysType struct {
	basetypes.ObjectType
}

func (t PrivateKeysType) Equal(o attr.Type) bool {
	other, ok := o.(PrivateKeysType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t PrivateKeysType) String() string {
	return "PrivateKeysType"
}

func (t PrivateKeysType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	createdAtAttribute, ok := attributes["created_at"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`created_at is missing from object`)

		return nil, diags
	}

	createdAtVal, ok := createdAtAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`created_at expected to be basetypes.StringValue, was: %T`, createdAtAttribute))
	}

	descriptionAttribute, ok := attributes["description"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`description is missing from object`)

		return nil, diags
	}

	descriptionVal, ok := descriptionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`description expected to be basetypes.StringValue, was: %T`, descriptionAttribute))
	}

	fingerprintAttribute, ok := attributes["fingerprint"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`fingerprint is missing from object`)

		return nil, diags
	}

	fingerprintVal, ok := fingerprintAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`fingerprint expected to be basetypes.StringValue, was: %T`, fingerprintAttribute))
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return nil, diags
	}

	idVal, ok := idAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be basetypes.Int64Value, was: %T`, idAttribute))
	}

	isGitRelatedAttribute, ok := attributes["is_git_related"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`is_git_related is missing from object`)

		return nil, diags
	}

	isGitRelatedVal, ok := isGitRelatedAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`is_git_related expected to be basetypes.BoolValue, was: %T`, isGitRelatedAttribute))
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return nil, diags
	}

	nameVal, ok := nameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be basetypes.StringValue, was: %T`, nameAttribute))
	}

	privateKeyAttribute, ok := attributes["private_key"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`private_key is missing from object`)

		return nil, diags
	}

	privateKeyVal, ok := privateKeyAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`private_key expected to be basetypes.StringValue, was: %T`, privateKeyAttribute))
	}

	teamIdAttribute, ok := attributes["team_id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`team_id is missing from object`)

		return nil, diags
	}

	teamIdVal, ok := teamIdAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`team_id expected to be basetypes.Int64Value, was: %T`, teamIdAttribute))
	}

	updatedAtAttribute, ok := attributes["updated_at"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`updated_at is missing from object`)

		return nil, diags
	}

	updatedAtVal, ok := updatedAtAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`updated_at expected to be basetypes.StringValue, was: %T`, updatedAtAttribute))
	}

	uuidAttribute, ok := attributes["uuid"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`uuid is missing from object`)

		return nil, diags
	}

	uuidVal, ok := uuidAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`uuid expected to be basetypes.StringValue, was: %T`, uuidAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return PrivateKeysValue{
		CreatedAt:    createdAtVal,
		Description:  descriptionVal,
		Fingerprint:  fingerprintVal,
		Id:           idVal,
		IsGitRelated: isGitRelatedVal,
		Name:         nameVal,
		PrivateKey:   privateKeyVal,
		TeamId:       teamIdVal,
		UpdatedAt:    updatedAtVal,
		Uuid:         uuidVal,
		state:        attr.ValueStateKnown,
	}, diags
}

func NewPrivateKeysValueNull() PrivateKeysValue {
	return PrivateKeysValue{
		state: attr.ValueStateNull,
	}
}

func NewPrivateKeysValueUnknown() PrivateKeysValue {
	return PrivateKeysValue{
		state: attr.ValueStateUnknown,
	}
}

func NewPrivateKeysValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (PrivateKeysValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing PrivateKeysValue Attribute Value",
				"While creating a PrivateKeysValue value, a missing attribute value was detected. "+
					"A PrivateKeysValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("PrivateKeysValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid PrivateKeysValue Attribute Type",
				"While creating a PrivateKeysValue value, an invalid attribute value was detected. "+
					"A PrivateKeysValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("PrivateKeysValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("PrivateKeysValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra PrivateKeysValue Attribute Value",
				"While creating a PrivateKeysValue value, an extra attribute value was detected. "+
					"A PrivateKeysValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra PrivateKeysValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewPrivateKeysValueUnknown(), diags
	}

	createdAtAttribute, ok := attributes["created_at"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`created_at is missing from object`)

		return NewPrivateKeysValueUnknown(), diags
	}

	createdAtVal, ok := createdAtAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`created_at expected to be basetypes.StringValue, was: %T`, createdAtAttribute))
	}

	descriptionAttribute, ok := attributes["description"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`description is missing from object`)

		return NewPrivateKeysValueUnknown(), diags
	}

	descriptionVal, ok := descriptionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`description expected to be basetypes.StringValue, was: %T`, descriptionAttribute))
	}

	fingerprintAttribute, ok := attributes["fingerprint"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`fingerprint is missing from object`)

		return NewPrivateKeysValueUnknown(), diags
	}

	fingerprintVal, ok := fingerprintAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`fingerprint expected to be basetypes.StringValue, was: %T`, fingerprintAttribute))
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return NewPrivateKeysValueUnknown(), diags
	}

	idVal, ok := idAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be basetypes.Int64Value, was: %T`, idAttribute))
	}

	isGitRelatedAttribute, ok := attributes["is_git_related"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`is_git_related is missing from object`)

		return NewPrivateKeysValueUnknown(), diags
	}

	isGitRelatedVal, ok := isGitRelatedAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`is_git_related expected to be basetypes.BoolValue, was: %T`, isGitRelatedAttribute))
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return NewPrivateKeysValueUnknown(), diags
	}

	nameVal, ok := nameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be basetypes.StringValue, was: %T`, nameAttribute))
	}

	privateKeyAttribute, ok := attributes["private_key"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`private_key is missing from object`)

		return NewPrivateKeysValueUnknown(), diags
	}

	privateKeyVal, ok := privateKeyAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`private_key expected to be basetypes.StringValue, was: %T`, privateKeyAttribute))
	}

	teamIdAttribute, ok := attributes["team_id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`team_id is missing from object`)

		return NewPrivateKeysValueUnknown(), diags
	}

	teamIdVal, ok := teamIdAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`team_id expected to be basetypes.Int64Value, was: %T`, teamIdAttribute))
	}

	updatedAtAttribute, ok := attributes["updated_at"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`updated_at is missing from object`)

		return NewPrivateKeysValueUnknown(), diags
	}

	updatedAtVal, ok := updatedAtAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`updated_at expected to be basetypes.StringValue, was: %T`, updatedAtAttribute))
	}

	uuidAttribute, ok := attributes["uuid"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`uuid is missing from object`)

		return NewPrivateKeysValueUnknown(), diags
	}

	uuidVal, ok := uuidAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`uuid expected to be basetypes.StringValue, was: %T`, uuidAttribute))
	}

	if diags.HasError() {
		return NewPrivateKeysValueUnknown(), diags
	}

	return PrivateKeysValue{
		CreatedAt:    createdAtVal,
		Description:  descriptionVal,
		Fingerprint:  fingerprintVal,
		Id:           idVal,
		IsGitRelated: isGitRelatedVal,
		Name:         nameVal,
		PrivateKey:   privateKeyVal,
		TeamId:       teamIdVal,
		UpdatedAt:    updatedAtVal,
		Uuid:         uuidVal,
		state:        attr.ValueStateKnown,
	}, diags
}

func NewPrivateKeysValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) PrivateKeysValue {
	object, diags := NewPrivateKeysValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewPrivateKeysValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t PrivateKeysType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewPrivateKeysValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewPrivateKeysValueUnknown(), nil
	}

	if in.IsNull() {
		return NewPrivateKeysValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewPrivateKeysValueMust(PrivateKeysValue{}.AttributeTypes(ctx), attributes), nil
}

func (t PrivateKeysType) ValueType(ctx context.Context) attr.Value {
	return PrivateKeysValue{}
}

var _ basetypes.ObjectValuable = PrivateKeysValue{}

type PrivateKeysValue struct {
	CreatedAt    basetypes.StringValue `tfsdk:"created_at"`
	Description  basetypes.StringValue `tfsdk:"description"`
	Fingerprint  basetypes.StringValue `tfsdk:"fingerprint"`
	Id           basetypes.Int64Value  `tfsdk:"id"`
	IsGitRelated basetypes.BoolValue   `tfsdk:"is_git_related"`
	Name         basetypes.StringValue `tfsdk:"name"`
	PrivateKey   basetypes.StringValue `tfsdk:"private_key"`
	TeamId       basetypes.Int64Value  `tfsdk:"team_id"`
	UpdatedAt    basetypes.StringValue `tfsdk:"updated_at"`
	Uuid         basetypes.StringValue `tfsdk:"uuid"`
	state        attr.ValueState
}

func (v PrivateKeysValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 10)

	var val tftypes.Value
	var err error

	attrTypes["created_at"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["description"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["fingerprint"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["id"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["is_git_related"] = basetypes.BoolType{}.TerraformType(ctx)
	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["private_key"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["team_id"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["updated_at"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["uuid"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 10)

		val, err = v.CreatedAt.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["created_at"] = val

		val, err = v.Description.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["description"] = val

		val, err = v.Fingerprint.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["fingerprint"] = val

		val, err = v.Id.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["id"] = val

		val, err = v.IsGitRelated.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["is_git_related"] = val

		val, err = v.Name.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["name"] = val

		val, err = v.PrivateKey.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["private_key"] = val

		val, err = v.TeamId.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["team_id"] = val

		val, err = v.UpdatedAt.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["updated_at"] = val

		val, err = v.Uuid.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["uuid"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v PrivateKeysValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v PrivateKeysValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v PrivateKeysValue) String() string {
	return "PrivateKeysValue"
}

func (v PrivateKeysValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"created_at":     basetypes.StringType{},
		"description":    basetypes.StringType{},
		"fingerprint":    basetypes.StringType{},
		"id":             basetypes.Int64Type{},
		"is_git_related": basetypes.BoolType{},
		"name":           basetypes.StringType{},
		"private_key":    basetypes.StringType{},
		"team_id":        basetypes.Int64Type{},
		"updated_at":     basetypes.StringType{},
		"uuid":           basetypes.StringType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"created_at":     v.CreatedAt,
			"description":    v.Description,
			"fingerprint":    v.Fingerprint,
			"id":             v.Id,
			"is_git_related": v.IsGitRelated,
			"name":           v.Name,
			"private_key":    v.PrivateKey,
			"team_id":        v.TeamId,
			"updated_at":     v.UpdatedAt,
			"uuid":           v.Uuid,
		})

	return objVal, diags
}

func (v PrivateKeysValue) Equal(o attr.Value) bool {
	other, ok := o.(PrivateKeysValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.CreatedAt.Equal(other.CreatedAt) {
		return false
	}

	if !v.Description.Equal(other.Description) {
		return false
	}

	if !v.Fingerprint.Equal(other.Fingerprint) {
		return false
	}

	if !v.Id.Equal(other.Id) {
		return false
	}

	if !v.IsGitRelated.Equal(other.IsGitRelated) {
		return false
	}

	if !v.Name.Equal(other.Name) {
		return false
	}

	if !v.PrivateKey.Equal(other.PrivateKey) {
		return false
	}

	if !v.TeamId.Equal(other.TeamId) {
		return false
	}

	if !v.UpdatedAt.Equal(other.UpdatedAt) {
		return false
	}

	if !v.Uuid.Equal(other.Uuid) {
		return false
	}

	return true
}

func (v PrivateKeysValue) Type(ctx context.Context) attr.Type {
	return PrivateKeysType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v PrivateKeysValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"created_at":     basetypes.StringType{},
		"description":    basetypes.StringType{},
		"fingerprint":    basetypes.StringType{},
		"id":             basetypes.Int64Type{},
		"is_git_related": basetypes.BoolType{},
		"name":           basetypes.StringType{},
		"private_key":    basetypes.StringType{},
		"team_id":        basetypes.Int64Type{},
		"updated_at":     basetypes.StringType{},
		"uuid":           basetypes.StringType{},
	}
}
