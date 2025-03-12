package application

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/util"
	sutil "terraform-provider-coolify/internal/service/util"
)

var (
	_ resource.Resource                = &dockerComposeApplicationResource{}
	_ resource.ResourceWithConfigure   = &dockerComposeApplicationResource{}
	_ resource.ResourceWithImportState = &dockerComposeApplicationResource{}
)

type dockerComposeApplicationResourceModel = dockerComposeApplicationModel

func NewDockerComposeApplicationResource() resource.Resource {
	return &dockerComposeApplicationResource{}
}

type dockerComposeApplicationResource struct {
	client *api.ClientWithResponses
}

func (r *dockerComposeApplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dockercompose_application"
}

func (r *dockerComposeApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	commonSchema := commonApplicationModel{}.CommonSchema(ctx)
	dockerComposeSchema := schema.Schema{
		Description: "Create, read, update, and delete a Coolify Dockercompose application resource.",
		Attributes: map[string]schema.Attribute{
			"dockercompose_raw": schema.StringAttribute{
				Optional:    true,
				Description: "Raw Dockercompose",
			},
		},
	}

	resp.Schema = sutil.MergeResourceSchemas(commonSchema, dockerComposeSchema)

}

func (r *dockerComposeApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.client, resp)
}

func (r *dockerComposeApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dockerComposeApplicationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Dockercompose application", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	createResp, err := r.client.CreateDockercomposeApplicationWithResponse(ctx, api.CreateDockercomposeApplicationJSONRequestBody{
		Description:      plan.Description.ValueStringPointer(),
		Name:             plan.Name.ValueStringPointer(),
		DestinationUuid:  plan.DestinationUuid.ValueStringPointer(),
		EnvironmentName:  plan.EnvironmentName.ValueString(),
		EnvironmentUuid:  plan.EnvironmentUuid.ValueString(),
		InstantDeploy:    plan.InstantDeploy.ValueBoolPointer(),
		ProjectUuid:      plan.ProjectUuid.ValueString(),
		ServerUuid:       plan.ServerUuid.ValueString(),
		DockerComposeRaw: *plan.DockerComposeRaw.ValueStringPointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Dockercompose application",
			err.Error(),
		)
		return
	}

	if createResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code creating Dockercompose application",
			fmt.Sprintf("Received %s creating Dockercompose application. Details: %s", createResp.Status(), createResp.Body),
		)
		return
	}

	data := r.ReadFromAPI(ctx, &resp.Diagnostics, *createResp.JSON201.Uuid, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dockerComposeApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dockerComposeApplicationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Dockercompose application", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	data := r.ReadFromAPI(ctx, &resp.Diagnostics, state.Uuid.ValueString(), state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dockerComposeApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dockerComposeApplicationResourceModel
	var state dockerComposeApplicationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := plan.Uuid.ValueString()

	tflog.Debug(ctx, "Updating Dockercompose application", map[string]interface{}{
		"uuid": uuid,
	})

	updateResp, err := r.client.UpdateApplicationByUuidWithResponse(ctx, uuid, api.UpdateApplicationByUuidJSONRequestBody{
		Description:      plan.Description.ValueStringPointer(),
		Name:             plan.Name.ValueStringPointer(),
		DestinationUuid:  plan.DestinationUuid.ValueStringPointer(),
		EnvironmentName:  plan.EnvironmentName.ValueStringPointer(),
		InstantDeploy:    plan.InstantDeploy.ValueBoolPointer(),
		ProjectUuid:      plan.ProjectUuid.ValueStringPointer(),
		ServerUuid:       plan.ServerUuid.ValueStringPointer(),
		DockerComposeRaw: plan.DockerComposeRaw.ValueStringPointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating Dockercompose application: uuid=%s", uuid),
			err.Error(),
		)
		return
	}

	if updateResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code updating Dockercompose application",
			fmt.Sprintf("Received %s updating Dockercompose application: uuid=%s. Details: %s", updateResp.Status(), uuid, updateResp.Body))
		return
	}

	if plan.InstantDeploy.ValueBool() {
		r.client.RestartApplicationByUuid(ctx, uuid)
	}

	data := r.ReadFromAPI(ctx, &resp.Diagnostics, uuid, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dockerComposeApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dockerComposeApplicationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Dockercompose application", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	deleteResp, err := r.client.DeleteApplicationByUuidWithResponse(ctx, state.Uuid.ValueString(), &api.DeleteApplicationByUuidParams{
		DeleteConfigurations:    types.BoolValue(true).ValueBoolPointer(),
		DeleteVolumes:           types.BoolValue(true).ValueBoolPointer(),
		DockerCleanup:           types.BoolValue(true).ValueBoolPointer(),
		DeleteConnectedNetworks: types.BoolValue(false).ValueBoolPointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Dockercompose application, got error: %s", err))
		return
	}

	if deleteResp.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code deleting Dockercompose application",
			fmt.Sprintf("Received %s deleting Dockercompose application: %s. Details: %s", deleteResp.Status(), state, deleteResp.Body))
		return
	}
}

func (r *dockerComposeApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ids := strings.Split(req.ID, "/")
	if len(ids) != 4 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID should be in the format: <server_uuid>/<project_uuid>/<environment_name>/<application_uuid>",
		)
		return
	}

	serverUuid, projectUuid, environmentName, uuid := ids[0], ids[1], ids[2], ids[3]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("server_uuid"), serverUuid)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_uuid"), projectUuid)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_name"), environmentName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), uuid)...)
}

// MARK: Helper functions

func (r *dockerComposeApplicationResource) ReadFromAPI(
	ctx context.Context,
	diags *diag.Diagnostics,
	uuid string,
	state dockerComposeApplicationResourceModel,
) dockerComposeApplicationResourceModel {
	readResp, err := r.client.GetApplicationByUuidWithResponse(ctx, uuid)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading Dockercompose application: uuid=%s", uuid),
			err.Error(),
		)
		return dockerComposeApplicationResourceModel{}
	}

	if readResp.StatusCode() != http.StatusOK {
		diags.AddError(
			"Unexpected HTTP status code reading Dockercompose application",
			fmt.Sprintf("Received %s for Dockercompose application: uuid=%s. Details: %s", readResp.Status(), uuid, readResp.Body))
		return dockerComposeApplicationResourceModel{}
	}

	result := dockerComposeApplicationResourceModel{}.FromAPI(readResp.JSON200, state)

	return result
}
