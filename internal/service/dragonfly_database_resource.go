package service

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
	"terraform-provider-coolify/internal/expand"
	"terraform-provider-coolify/internal/provider/util"
	sutil "terraform-provider-coolify/internal/service/util"
)

var (
	_ resource.Resource                = &dragonflyDatabaseResource{}
	_ resource.ResourceWithConfigure   = &dragonflyDatabaseResource{}
	_ resource.ResourceWithImportState = &dragonflyDatabaseResource{}
	_ resource.ResourceWithModifyPlan  = &dragonflyDatabaseResource{}
)

type dragonflyDatabaseResourceModel = dragonflyDatabaseModel

func NewDragonflyDatabaseResource() resource.Resource {
	return &dragonflyDatabaseResource{}
}

type dragonflyDatabaseResource struct {
	client *api.ClientWithResponses
}

func (r *dragonflyDatabaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dragonfly_database"
}

func (r *dragonflyDatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	commonSchema := commonDatabaseModel{}.CommonSchema(ctx)
	dragonflySchema := schema.Schema{
		Description: "Create, read, update, and delete a Coolify database (DragonFly) resource.",
		Attributes: map[string]schema.Attribute{
			"dragonfly_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "DragonFly password",
			},
		},
	}

	resp.Schema = sutil.MergeResourceSchemas(commonSchema, dragonflySchema)
}

func (r *dragonflyDatabaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.client, resp)
}

func (r *dragonflyDatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dragonflyDatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating DragonFly database", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	createResp, err := r.client.CreateDatabaseDragonflyWithResponse(ctx, api.CreateDatabaseDragonflyJSONRequestBody{
		Description:             plan.Description.ValueStringPointer(),
		Name:                    plan.Name.ValueStringPointer(),
		DestinationUuid:         plan.DestinationUuid.ValueStringPointer(),
		EnvironmentName:         plan.EnvironmentName.ValueString(),
		EnvironmentUuid:         plan.EnvironmentUuid.ValueString(),
		Image:                   plan.Image.ValueStringPointer(),
		InstantDeploy:           plan.InstantDeploy.ValueBoolPointer(),
		IsPublic:                plan.IsPublic.ValueBoolPointer(),
		LimitsCpuShares:         expand.Int64(plan.LimitsCpuShares),
		LimitsCpus:              plan.LimitsCpus.ValueStringPointer(),
		LimitsCpuset:            plan.LimitsCpuset.ValueStringPointer(),
		LimitsMemory:            plan.LimitsMemory.ValueStringPointer(),
		LimitsMemoryReservation: plan.LimitsMemoryReservation.ValueStringPointer(),
		LimitsMemorySwap:        plan.LimitsMemorySwap.ValueStringPointer(),
		LimitsMemorySwappiness:  expand.Int64(plan.LimitsMemorySwappiness),
		DragonflyPassword:       plan.DragonflyPassword.ValueStringPointer(),
		ProjectUuid:             plan.ProjectUuid.ValueString(),
		PublicPort:              expand.Int64(plan.PublicPort),
		ServerUuid:              plan.ServerUuid.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating DragonFly database",
			err.Error(),
		)
		return
	}

	if createResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code creating DragonFly database",
			fmt.Sprintf("Received %s creating DragonFly database. Details: %s", createResp.Status(), createResp.Body),
		)
		return
	}

	data, _ := r.ReadFromAPI(ctx, &resp.Diagnostics, createResp.JSON201.Uuid, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dragonflyDatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dragonflyDatabaseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading DragonFly database", map[string]interface{}{
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

func (r *dragonflyDatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dragonflyDatabaseResourceModel
	var state dragonflyDatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := plan.Uuid.ValueString()

	tflog.Debug(ctx, "Updating DragonFly database", map[string]interface{}{
		"uuid": uuid,
	})

	updateResp, err := r.client.UpdateDatabaseByUuidWithResponse(ctx, uuid, api.UpdateDatabaseByUuidJSONRequestBody{
		Description:             plan.Description.ValueStringPointer(),
		Image:                   plan.Image.ValueStringPointer(),
		IsPublic:                plan.IsPublic.ValueBoolPointer(),
		LimitsCpuShares:         expand.Int64(plan.LimitsCpuShares),
		LimitsCpus:              plan.LimitsCpus.ValueStringPointer(),
		LimitsCpuset:            plan.LimitsCpuset.ValueStringPointer(),
		LimitsMemory:            plan.LimitsMemory.ValueStringPointer(),
		LimitsMemoryReservation: plan.LimitsMemoryReservation.ValueStringPointer(),
		LimitsMemorySwap:        plan.LimitsMemorySwap.ValueStringPointer(),
		LimitsMemorySwappiness:  expand.Int64(plan.LimitsMemorySwappiness),
		Name:                    plan.Name.ValueStringPointer(),
		DragonflyPassword:       plan.DragonflyPassword.ValueStringPointer(),
		PublicPort:              expand.Int64(plan.PublicPort),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating DragonFly database: uuid=%s", uuid),
			err.Error(),
		)
		return
	}

	if updateResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code updating DragonFly database",
			fmt.Sprintf("Received %s updating DragonFly database: uuid=%s. Details: %s", updateResp.Status(), uuid, updateResp.Body))
		return
	}

	if plan.InstantDeploy.ValueBool() {
		r.client.RestartDatabaseByUuid(ctx, uuid)
	}

	data, ok := r.ReadFromAPI(ctx, &resp.Diagnostics, uuid, plan)
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dragonflyDatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dragonflyDatabaseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting DragonFly database", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	deleteResp, err := r.client.DeleteDatabaseByUuidWithResponse(ctx, state.Uuid.ValueString(), &api.DeleteDatabaseByUuidParams{
		DeleteConfigurations:    types.BoolValue(true).ValueBoolPointer(),
		DeleteVolumes:           types.BoolValue(true).ValueBoolPointer(),
		DockerCleanup:           types.BoolValue(true).ValueBoolPointer(),
		DeleteConnectedNetworks: types.BoolValue(false).ValueBoolPointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete DragonFly database, got error: %s", err))
		return
	}

	if deleteResp.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code deleting DragonFly database",
			fmt.Sprintf("Received %s deleting DragonFly database: %s. Details: %s", deleteResp.Status(), state, deleteResp.Body))
		return
	}
}

func (r *dragonflyDatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ids := strings.Split(req.ID, "/")
	if len(ids) != 4 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID should be in the format: <server_uuid>/<project_uuid>/<environment_name>/<database_uuid>",
		)
		return
	}

	serverUuid, projectUuid, environmentName, uuid := ids[0], ids[1], ids[2], ids[3]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("server_uuid"), serverUuid)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_uuid"), projectUuid)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_name"), environmentName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), uuid)...)
}

func (r *dragonflyDatabaseResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *dragonflyDatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() || plan == nil || state == nil {
		return
	}

	// If the password changes, the internal URL will change
	if !plan.DragonflyPassword.Equal(state.DragonflyPassword) {
		plan.InternalDbUrl = types.StringUnknown()
		resp.Plan.Set(ctx, &plan)
	}
}

// MARK: Helper functions

func (r *dragonflyDatabaseResource) ReadFromAPI(
	ctx context.Context,
	diags *diag.Diagnostics,
	uuid string,
	state dragonflyDatabaseResourceModel,
) (dragonflyDatabaseResourceModel, bool) {
	readResp, err := r.client.GetDatabaseByUuidWithResponse(ctx, uuid)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading DragonFly database: uuid=%s", uuid),
			err.Error(),
		)
		return dragonflyDatabaseResourceModel{}, false
	}

	if readResp.StatusCode() == http.StatusNotFound {
		return dragonflyDatabaseResourceModel{}, false
	}

	if readResp.StatusCode() != http.StatusOK {
		diags.AddError(
			"Unexpected HTTP status code reading DragonFly database",
			fmt.Sprintf("Received %s for DragonFly database: uuid=%s. Details: %s", readResp.Status(), uuid, readResp.Body))
		return dragonflyDatabaseResourceModel{}, false
	}

	result, err := dragonflyDatabaseResourceModel{}.FromAPI(readResp.JSON200, state)
	if err != nil {
		diags.AddError("Error converting API response to model", err.Error())
		return dragonflyDatabaseResourceModel{}, false
	}

	return result, true
}
