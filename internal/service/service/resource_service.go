package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/util"
)

var (
	_ resource.Resource                = &ServiceResource{}
	_ resource.ResourceWithConfigure   = &ServiceResource{}
	_ resource.ResourceWithImportState = &ServiceResource{}
)

type ServiceResourceModel = ServiceModel

func NewServiceResource() resource.Resource {
	return &ServiceResource{}
}

type ServiceResource struct {
	client *api.ClientWithResponses
}

func (r *ServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (r *ServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ServiceModel{}.Schema(ctx)
}

func (r *ServiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.client, resp)
}

func (r *ServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServiceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating service", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	res, err := r.client.CreateServiceWithResponse(ctx, plan.ToAPICreate())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating service",
			err.Error(),
		)
		return
	}

	if res.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code creating service",
			fmt.Sprintf("Received %s creating service. Details: %s", res.Status(), res.Body),
		)
		return
	}

	data, _ := r.ReadFromAPI(ctx, &resp.Diagnostics, *res.JSON201.Uuid, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServiceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading service", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	data, ok := r.ReadFromAPI(ctx, &resp.Diagnostics, state.Uuid.ValueString(), state)
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServiceResourceModel
	var state ServiceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := plan.Uuid.ValueString()

	tflog.Debug(ctx, "Updating service", map[string]interface{}{
		"uuid": uuid,
	})

	updateResp, err := r.client.UpdateServiceByUuidWithResponse(ctx, uuid, plan.ToAPIUpdate())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating service: uuid=%s", uuid),
			err.Error(),
		)
		return
	}

	if updateResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code updating service",
			fmt.Sprintf("Received %s updating service: uuid=%s. Details: %s", updateResp.Status(), uuid, updateResp.Body))
		return
	}

	if plan.InstantDeploy.ValueBool() {
		r.client.RestartServiceByUuid(ctx, uuid)
	}

	data, ok := r.ReadFromAPI(ctx, &resp.Diagnostics, uuid, plan)
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServiceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting service", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	deleteResp, err := r.client.DeleteServiceByUuidWithResponse(ctx, state.Uuid.ValueString(), &api.DeleteServiceByUuidParams{
		DeleteConfigurations:    types.BoolValue(true).ValueBoolPointer(),
		DeleteVolumes:           types.BoolValue(true).ValueBoolPointer(),
		DockerCleanup:           types.BoolValue(true).ValueBoolPointer(),
		DeleteConnectedNetworks: types.BoolValue(false).ValueBoolPointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service, got error: %s", err))
		return
	}

	if deleteResp.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code deleting service",
			fmt.Sprintf("Received %s deleting service: %s. Details: %s", deleteResp.Status(), state, deleteResp.Body))
		return
	}
}

func (r *ServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ids := strings.Split(req.ID, "/")
	if len(ids) != 4 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID should be in the format: <server_uuid>/<project_uuid>/<environment_name>/<service_uuid>",
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

func (r *ServiceResource) ReadFromAPI(
	ctx context.Context,
	diags *diag.Diagnostics,
	uuid string,
	state ServiceResourceModel,
) (ServiceResourceModel, bool) {
	res, err := r.client.GetServiceByUuidWithResponse(ctx, uuid)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading service: uuid=%s", uuid),
			err.Error(),
		)
		return ServiceResourceModel{}, false
	}

	if res.StatusCode() == http.StatusNotFound {
		return ServiceResourceModel{}, false
	}

	if res.StatusCode() != http.StatusOK {
		diags.AddError(
			"Unexpected HTTP status code reading service",
			fmt.Sprintf("Received %s for service: uuid=%s. Details: %s", res.Status(), uuid, res.Body))
		return ServiceResourceModel{}, false
	}

	result := ServiceResourceModel{}.FromAPI(res.JSON200, state)

	return result, true
}
