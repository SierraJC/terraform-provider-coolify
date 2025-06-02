package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/flatten"
	"terraform-provider-coolify/internal/provider/generated/datasource_project"
	"terraform-provider-coolify/internal/provider/generated/resource_project"
	"terraform-provider-coolify/internal/provider/util"
)

var (
	_ resource.Resource                = &projectResource{}
	_ resource.ResourceWithConfigure   = &projectResource{}
	_ resource.ResourceWithImportState = &projectResource{}
)

func NewProjectResource() resource.Resource {
	return &projectResource{}
}

type projectResource struct {
	client *api.ClientWithResponses
}

func (r *projectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *projectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_project.ProjectResourceSchema(ctx)
	resp.Schema.Description = "Create, read, update, and delete a Coolify project resource."

	if nameAttr, ok := resp.Schema.Attributes["name"].(schema.StringAttribute); ok {
		nameAttr.Required = true
		nameAttr.Optional = false
		nameAttr.Computed = false
		resp.Schema.Attributes["name"] = nameAttr
	}

	validateNonEmptyStrings := []string{"name", "description"}
	for _, attr := range validateNonEmptyStrings {
		makeResourceAttributeNonEmpty(resp.Schema.Attributes, attr)
	}
}

func (r *projectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.client, resp)
}

func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resource_project.ProjectModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating project", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})
	createResp, err := r.client.CreateProjectWithResponse(ctx, api.CreateProjectJSONRequestBody{
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.Name.ValueStringPointer(),
	})

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

	data, _ := r.ReadFromAPI(ctx, &resp.Diagnostics, *createResp.JSON201.Uuid)
	r.copyMissingAttributes(&plan, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resource_project.ProjectModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading project", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	data, ok := r.ReadFromAPI(ctx, &resp.Diagnostics, state.Uuid.ValueString())
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}
	r.copyMissingAttributes(&state, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resource_project.ProjectModel
	var state resource_project.ProjectModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()

	if uuid == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	// Update API call logic
	tflog.Debug(ctx, "Updating project", map[string]interface{}{
		"uuid": uuid,
	})
	updateResp, err := r.client.UpdateProjectByUuidWithResponse(ctx, uuid, api.UpdateProjectByUuidJSONRequestBody{
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.Name.ValueStringPointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating project: uuid=%s", uuid),
			err.Error(),
		)
		return
	}

	if updateResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code updating project",
			fmt.Sprintf("Received %s updating project: uuid=%s. Details: %s", updateResp.Status(), uuid, updateResp.Body))
		return
	}

	data, ok := r.ReadFromAPI(ctx, &resp.Diagnostics, uuid)
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}
	r.copyMissingAttributes(&plan, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resource_project.ProjectModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting project", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	deleteResp, err := r.client.DeleteProjectByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete project, got error: %s", err))
		return
	}

	if deleteResp.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code deleting project",
			fmt.Sprintf("Received %s deleting project: %s. Details: %s", deleteResp.Status(), state, deleteResp.Body))
		return
	}
}

func (r *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *projectResource) copyMissingAttributes(
	plan *resource_project.ProjectModel,
	data *resource_project.ProjectModel,
) {

}

func (r *projectResource) ReadFromAPI(
	ctx context.Context,
	diags *diag.Diagnostics,
	uuid string,
) (resource_project.ProjectModel, bool) {
	readResp, err := r.client.GetProjectByUuidWithResponse(ctx, uuid)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading project: uuid=%s", uuid),
			err.Error(),
		)
		return resource_project.ProjectModel{}, false
	}

	if readResp.StatusCode() == http.StatusNotFound {
		return resource_project.ProjectModel{}, false
	}

	if readResp.StatusCode() != http.StatusOK {
		diags.AddError(
			"Unexpected HTTP status code reading project",
			fmt.Sprintf("Received %s for project: uuid=%s. Details: %s", readResp.Status(), uuid, readResp.Body))
		return resource_project.ProjectModel{}, false
	}

	return r.ApiToModel(ctx, diags, readResp.JSON200), true
}

func (r *projectResource) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *api.Project,
) resource_project.ProjectModel {
	var elements []attr.Value
	for _, env := range *response.Environments {
		attributes := map[string]attr.Value{
			"created_at":  flatten.String(env.CreatedAt),
			"description": flatten.String(env.Description),
			"id":          flatten.Int64(env.Id),
			"name":        flatten.String(env.Name),
			"project_id":  flatten.Int64(env.ProjectId),
			"updated_at":  flatten.String(env.UpdatedAt),
		}

		data, diag := datasource_project.NewEnvironmentsValue(
			datasource_project.EnvironmentsValue{}.AttributeTypes(ctx),
			attributes)
		diags.Append(diag...)
		elements = append(elements, data)
	}
	dataList, diag := types.ListValueFrom(ctx, datasource_project.EnvironmentsValue{}.Type(ctx), elements)
	diags.Append(diag...)

	return resource_project.ProjectModel{
		Description:  flatten.String(response.Description),
		Environments: dataList,
		Id:           flatten.Int64(response.Id),
		Name:         flatten.String(response.Name),
		Uuid:         flatten.String(response.Uuid),
	}
}
