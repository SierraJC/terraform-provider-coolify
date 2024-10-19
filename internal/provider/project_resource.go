package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/generated/resource_project"
	"terraform-provider-coolify/internal/provider/util"
)

var (
	_ resource.Resource                = (*projectResource)(nil)
	_ resource.ResourceWithConfigure   = (*projectResource)(nil)
	_ resource.ResourceWithImportState = (*projectResource)(nil)
)

func NewProjectResource() resource.Resource {
	return &projectResource{}
}

type projectResource struct {
	providerData CoolifyProviderData
}

// type projectResourceModel struct {
// 	Id types.String `tfsdk:"id"`
// }

func (r *projectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *projectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_project.ProjectResourceSchema(ctx)
}

func (r *projectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.providerData, resp)
}

func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state resource_project.ProjectModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := api.CreateProjectJSONRequestBody{
		Name:        state.Name.ValueStringPointer(),
		Description: state.Description.ValueStringPointer(),
	}

	tflog.Debug(ctx, "Creating project")
	createResp, err := r.providerData.client.CreateProjectWithResponse(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			err.Error(),
		)
		return
	}

	if createResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code creating project",
			fmt.Sprintf("Received %s creating project. Details: %s", createResp.Status(), createResp.Body),
		)
		return
	}

	assignStr(createResp.JSON201.Uuid, &state.Uuid)

	// GET /projects/{uuid}
	readResp, err := r.providerData.client.GetProjectByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading project: uuid=%s", state.Uuid.ValueString()),
			err.Error(),
		)
		return
	}

	if readResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading project",
			fmt.Sprintf("Received %s for project: uuid=%s. Details: %s", readResp.Status(), state.Uuid.ValueString(), readResp.Body))
		return
	}

	assignInt(readResp.JSON200.Id, &state.Id)
	assignStr(readResp.JSON200.Name, &state.Name)
	assignStr(readResp.JSON200.Uuid, &state.Uuid)
	assignStr(readResp.JSON200.Description, &state.Description)

	// Handle environments
	{
		var diags diag.Diagnostics
		var values []resource_project.EnvironmentsValue
		value := resource_project.EnvironmentsValue{}

		for _, env := range *readResp.JSON200.Environments {
			fr := resource_project.NewEnvironmentsValueMust(value.AttributeTypes(ctx), map[string]attr.Value{
				"created_at":  attrStr(env.CreatedAt),
				"description": attrStr(env.Description),
				"id":          attrInt(env.Id),
				"name":        attrStr(env.Name),
				"project_id":  attrInt(env.ProjectId),
				"updated_at":  attrStr(env.UpdatedAt),
			})
			values = append(values, fr)
		}

		envs, diag := types.ListValueFrom(ctx, value.Type(ctx), values)
		diags.Append(diag...)
		if diags.HasError() {
			// return basetypes.NewListUnknown(value.Type(ctx)), diags
			return
		}

		state.Environments = envs
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resource_project.ProjectModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading project: uuid=%s", state.Uuid.ValueString()))
	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	// GET /projects/{uuid}
	readResp, err := r.providerData.client.GetProjectByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading project: uuid=%s", state.Uuid.ValueString()),
			err.Error(),
		)
		return
	}

	if readResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading project",
			fmt.Sprintf("Received %s for project: uuid=%s. Details: %s", readResp.Status(), state.Uuid.ValueString(), readResp.Body))
		return
	}

	assignInt(readResp.JSON200.Id, &state.Id)
	assignStr(readResp.JSON200.Name, &state.Name)
	assignStr(readResp.JSON200.Uuid, &state.Uuid)
	assignStr(readResp.JSON200.Description, &state.Description)

	// Handle environments
	{
		var diags diag.Diagnostics
		var values []resource_project.EnvironmentsValue
		value := resource_project.EnvironmentsValue{}

		for _, env := range *readResp.JSON200.Environments {
			fr := resource_project.NewEnvironmentsValueMust(value.AttributeTypes(ctx), map[string]attr.Value{
				"created_at":  attrStr(env.CreatedAt),
				"description": attrStr(env.Description),
				"id":          attrInt(env.Id),
				"name":        attrStr(env.Name),
				"project_id":  attrInt(env.ProjectId),
				"updated_at":  attrStr(env.UpdatedAt),
			})
			values = append(values, fr)
		}

		envs, diag := types.ListValueFrom(ctx, value.Type(ctx), values)
		diags.Append(diag...)
		if diags.HasError() {
			// return basetypes.NewListUnknown(value.Type(ctx)), diags
			return
		}

		state.Environments = envs
	}

	// elements := []attr.Value{}
	// for _, env := range *readResp.JSON200.Environments {

	// 	// envValue := resource_project.EnvironmentsValue{
	// 	// 	// CreatedAt:   types.StringValue(env.CreatedAt),
	// 	// 	// Description: types.StringValue(env.Description),
	// 	// 	// Id:          types.Int64Value(int64(env.Id)),
	// 	// 	// Name:        types.StringValue(env.Name),
	// 	// 	// ProjectId:   types.Int64Value(int64(env.ProjectId)),
	// 	// 	// UpdatedAt:   types.StringValue(env.UpdatedAt),
	// 	// }
	// 	// assignStr(env.CreatedAt, &envValue.CreatedAt)
	// 	// assignStr(env.Description, &envValue.Description)
	// 	// assignInt(env.Id, &envValue.Id)
	// 	// assignStr(env.Name, &envValue.Name)
	// 	// assignInt(env.ProjectId, &envValue.ProjectId)
	// 	// assignStr(env.UpdatedAt, &envValue.UpdatedAt)

	// 	state.Environments, _ = state.Environments.
	// }

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state resource_project.ProjectModel
	var plan resource_project.ProjectModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	body := api.UpdateProjectByUuidJSONRequestBody{
		Name:        plan.Name.ValueStringPointer(),
		Description: plan.Description.ValueStringPointer(),
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating project: uuid=%s", state.Uuid.ValueString()))
	updateResp, err := r.providerData.client.UpdateProjectByUuidWithResponse(ctx, state.Uuid.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating project: uuid=%s", state.Uuid.ValueString()),
			err.Error(),
		)
		return
	}

	if updateResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code updating project",
			fmt.Sprintf("Received %s updating project: uuid=%s. Details: %s", updateResp.Status(), state.Uuid.ValueString(), updateResp.Body))
		return
	}

	assignStr(updateResp.JSON201.Uuid, &state.Uuid)

	// GET /projects/{uuid}
	readResp, err := r.providerData.client.GetProjectByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading project: uuid=%s", state.Uuid.ValueString()),
			err.Error(),
		)
		return
	}

	if readResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading project",
			fmt.Sprintf("Received %s for project: uuid=%s. Details: %s", readResp.Status(), state.Uuid.ValueString(), readResp.Body))
		return
	}

	assignInt(readResp.JSON200.Id, &state.Id)
	assignStr(readResp.JSON200.Name, &state.Name)
	assignStr(readResp.JSON200.Uuid, &state.Uuid)
	assignStr(readResp.JSON200.Description, &state.Description)

	// Handle environments
	{
		var diags diag.Diagnostics
		var values []resource_project.EnvironmentsValue
		value := resource_project.EnvironmentsValue{}

		for _, env := range *readResp.JSON200.Environments {
			fr := resource_project.NewEnvironmentsValueMust(value.AttributeTypes(ctx), map[string]attr.Value{
				"created_at":  attrStr(env.CreatedAt),
				"description": attrStr(env.Description),
				"id":          attrInt(env.Id),
				"name":        attrStr(env.Name),
				"project_id":  attrInt(env.ProjectId),
				"updated_at":  attrStr(env.UpdatedAt),
			})
			values = append(values, fr)
		}

		envs, diag := types.ListValueFrom(ctx, value.Type(ctx), values)
		diags.Append(diag...)
		if diags.HasError() {
			// return basetypes.NewListUnknown(value.Type(ctx)), diags
			return
		}

		state.Environments = envs
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resource_project.ProjectModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteResp, err := r.providerData.client.DeleteProjectByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete project, got error: %s", err))
		return
	}

	if deleteResp.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code deleting project",
			fmt.Sprintf("Received %s deleting project: %s. Details: %s", deleteResp.Status(), state.Uuid.ValueString(), deleteResp.Body))
		return
	}
}

func (r *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
