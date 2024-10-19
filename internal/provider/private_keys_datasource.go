package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/generated/datasource_private_keys"
)

var _ datasource.DataSource = (*privateKeysDataSource)(nil)

func NewPrivateKeysDataSource() datasource.DataSource {
	return &privateKeysDataSource{}
}

type privateKeysDataSource struct {
	client *api.APIClient
}

func (d *privateKeysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_keys"
}

func (d *privateKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_private_keys.PrivateKeysDataSourceSchema(ctx)
}

func (d *privateKeysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*api.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected coolifyProviderConfig, got: %T. Please report this issue to the provider developers.",
		)
		return
	}

	d.client = client
}

func (d *privateKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_private_keys.PrivateKeysModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read Coolify API - GET /security/keys
	res, err := d.client.ListPrivateKeysWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Private Keys",
			"Unexpected API error: "+err.Error(),
		)
	}

	if res.StatusCode() != 200 {
		resp.Diagnostics.AddError("API call failed", "Expected HTTP 200 but received "+res.Status())
	}

	deviceApStat, diag := SdkToTerraform(ctx, res.JSON200)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := resp.State.SetAttribute(ctx, path.Root("private_keys"), deviceApStat); err != nil {
		resp.Diagnostics.Append(err...)
		return
	}

	// -------------
	// var elements []datasource_private_keys.PrivateKeysValue
	// for _, d := range *res.JSON200 {

	// 	val := datasource_private_keys.PrivateKeysValue{
	// 		CreatedAt:    types.StringValue(*d.CreatedAt),
	// 		Description:  types.StringValue(*d.Description),
	// 		Id:           types.Int64Value(int64(*d.Id)),
	// 		IsGitRelated: types.BoolValue(*d.IsGitRelated),
	// 		Name:         types.StringValue(*d.Name),
	// 		PrivateKey:   types.StringValue(*d.PrivateKey),
	// 		TeamId:       types.Int64Value(int64(*d.TeamId)),
	// 		UpdatedAt:    types.StringValue(*d.UpdatedAt),
	// 		Uuid:         types.StringValue(*d.Uuid),
	// 	}

	// 	objVal, diag := val.ToObjectValue(ctx)

	// newObjVal, _ := datasource_private_keys.NewPrivateKeysValue(objVal.AttributeTypes(ctx), objVal.Attributes())
	// dataSet, diag := types.SetValueFrom(ctx, datasource_private_keys.PrivateKeysValue{}.Type(ctx), []datasource_private_keys.PrivateKeysValue{newObjVal})
	// fmt.Printf("dataSet: %+v", dataSet)
	// role := autogen.RolesValue{
	// 	DatabaseName: types.StringValue(v.DatabaseName),
	// 				RoleName:     types.StringValue(v.RoleName),
	// }
	// objVal, _ := value.ToObjectValue(ctx)
	// newRoleValue, _ := autogen.NewRolesValue(objVal.AttributeTypes(ctx), objVal.Attributes())
	// rolesSet, diagnostic := types.SetValueFrom(ctx, autogen.RolesValue{}.Type(ctx), []autogen.RolesValue{newRoleValue})

	// 	resp.Diagnostics.Append(diag...)
	// 	if resp.Diagnostics.HasError() {
	// 		return
	// 	}
	// 	nodesValue, diags := datasource_private_keys.NewPrivateKeysValue(objVal.AttributeTypes(ctx), objVal.Attributes())
	// 	resp.Diagnostics.Append(diags...)
	// 	if resp.Diagnostics.HasError() {
	// 		return
	// 	}
	// 	elements = append(elements, nodesValue)

	// }

	// nodesListVal, diags := types.SetValueFrom(ctx, datasource_private_keys.PrivateKeysValue{}.Type(ctx), elements)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// data.PrivateKeys = nodesListVal
	// resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// data, _ := datasource_private_keys.PrivateKeysValue{
	// 	CreatedAt:    created_at,
	// 	Description:  description,
	// 	Id:           id,
	// 	IsGitRelated: is_git_related,
	// 	Name:         name,
	// 	PrivateKey:   private_key,
	// 	TeamId:       team_id,
	// 	UpdatedAt:    updated_at,
	// 	Uuid:         uuid,
	// }.ToObjectValue(ctx)

	// a, b := datasource_private_keys.NewPrivateKeysValue(data.AttributeTypes(ctx), data.Attributes())

	// model := datasource_private_keys.PrivateKeysModel{
	// 	PrivateKeys: deviceApStat,
	// }

	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// if err := resp.State.Set(ctx, model); err != nil {
	// 	resp.Diagnostics.Append(err...)
	// 	return
	// }
}

func SdkToTerraform(ctx context.Context, pk *[]api.PrivateKey) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// var elements []datasource_private_keys.PrivateKeysValue
	// for _, d := range *pk {
	// 	elem := siteSdkToTerraform(ctx, &diags, d)
	// 	elements = append(elements, elem)
	// }

	// dataSet, diag := types.SetValueFrom(context.Background(), datasource_private_keys.PrivateKeysType{}, elements)
	// if diag.HasError() {
	// 	diags.Append(diag...)
	// }

	// keys, d := tftypes.ListValueFrom(ctx, types.ObjectType{AttrTypes: datasource_private_keys.PrivateKeysType{}.AttributeTypes()}, elements)
	// diags.Append(d...)

	// fmt.Printf("keys: %+v", keys)

	var elements []attr.Value
	for _, d := range *pk {
		elem := siteSdkToTerraform(ctx, &diags, d)
		elements = append(elements, elem)
	}

	dataSet, err := types.SetValue(datasource_private_keys.PrivateKeysValue{}.Type(ctx), elements)
	if err != nil {
		diags.Append(err...)
	}

	return dataSet, diags
}

func siteSdkToTerraform(ctx context.Context, diags *diag.Diagnostics, d api.PrivateKey) datasource_private_keys.PrivateKeysValue {
	var created_at basetypes.StringValue
	var description basetypes.StringValue
	var id basetypes.Int64Value
	var is_git_related basetypes.BoolValue
	var name basetypes.StringValue
	var private_key basetypes.StringValue
	var team_id basetypes.Int64Value
	var updated_at basetypes.StringValue
	var uuid basetypes.StringValue

	if d.CreatedAt != nil {
		created_at = types.StringValue(*d.CreatedAt)
	}
	if d.Description != nil {
		description = types.StringValue(*d.Description)
	}
	if d.Id != nil {
		id = types.Int64Value(int64(*d.Id))
	}
	if d.IsGitRelated != nil {
		is_git_related = types.BoolValue(*d.IsGitRelated)
	}
	if d.Name != nil {
		name = types.StringValue(*d.Name)
	}
	if d.PrivateKey != nil {
		private_key = types.StringValue(*d.PrivateKey)
	}
	if d.TeamId != nil {
		team_id = types.Int64Value(int64(*d.TeamId))
	}
	if d.UpdatedAt != nil {
		updated_at = types.StringValue(*d.UpdatedAt)
	}
	if d.Uuid != nil {
		uuid = types.StringValue(*d.Uuid)
	}

	data_map_attr_type := datasource_private_keys.PrivateKeysValue{}.AttributeTypes(ctx)
	data_map_value := map[string]attr.Value{
		"created_at":     created_at,
		"description":    description,
		"id":             id,
		"is_git_related": is_git_related,
		"name":           name,
		"private_key":    private_key,
		"team_id":        team_id,
		"updated_at":     updated_at,
		"uuid":           uuid,
	}
	data, diag := datasource_private_keys.NewPrivateKeysValue(data_map_attr_type, data_map_value)
	diags.Append(diag...)

	return data
}
