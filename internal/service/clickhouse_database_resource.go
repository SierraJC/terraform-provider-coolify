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
	_ resource.Resource                = &clickhouseDatabaseResource{}
	_ resource.ResourceWithConfigure   = &clickhouseDatabaseResource{}
	_ resource.ResourceWithImportState = &clickhouseDatabaseResource{}
	_ resource.ResourceWithModifyPlan  = &clickhouseDatabaseResource{}
)

type clickhouseDatabaseResourceModel = clickhouseDatabaseModel

func NewClickhouseDatabaseResource() resource.Resource {
	return &clickhouseDatabaseResource{}
}

type clickhouseDatabaseResource struct {
	client *api.ClientWithResponses
}

func (r *clickhouseDatabaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clickhouse_database"
}

func (r *clickhouseDatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	commonSchema := commonDatabaseModel{}.CommonSchema(ctx)
	clickhouseSchema := schema.Schema{
		Description: "Create, read, update, and delete a Coolify database (Clickhouse) resource.",
		Attributes: map[string]schema.Attribute{
			"clickhouse_admin_user": schema.StringAttribute{
				Required:    true,
				Description: "Clickhouse admin user",
			},
			"clickhouse_admin_password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Clickhouse admin password",
			},
		},
	}

	resp.Schema = sutil.MergeResourceSchemas(commonSchema, clickhouseSchema)
}

func (r *clickhouseDatabaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.client, resp)
}

func (r *clickhouseDatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan clickhouseDatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Clickhouse database", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	createResp, err := r.client.CreateDatabaseClickhouseWithResponse(ctx, api.CreateDatabaseClickhouseJSONRequestBody{
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
		ClickhouseAdminUser:     plan.ClickhouseAdminUser.ValueStringPointer(),
		ClickhouseAdminPassword: plan.ClickhouseAdminPassword.ValueStringPointer(),
		ProjectUuid:             plan.ProjectUuid.ValueString(),
		PublicPort:              expand.Int64(plan.PublicPort),
		ServerUuid:              plan.ServerUuid.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Clickhouse database",
			err.Error(),
		)
		return
	}

	if createResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code creating Clickhouse database",
			fmt.Sprintf("Received %s creating Clickhouse database. Details: %s", createResp.Status(), createResp.Body),
		)
		return
	}

	data, _ := r.ReadFromAPI(ctx, &resp.Diagnostics, createResp.JSON201.Uuid, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clickhouseDatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state clickhouseDatabaseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Clickhouse database", map[string]interface{}{
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

func (r *clickhouseDatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan clickhouseDatabaseResourceModel
	var state clickhouseDatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := plan.Uuid.ValueString()

	tflog.Debug(ctx, "Updating Clickhouse database", map[string]interface{}{
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
		ClickhouseAdminUser:     plan.ClickhouseAdminUser.ValueStringPointer(),
		ClickhouseAdminPassword: plan.ClickhouseAdminPassword.ValueStringPointer(),
		PublicPort:              expand.Int64(plan.PublicPort),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating Clickhouse database: uuid=%s", uuid),
			err.Error(),
		)
		return
	}

	if updateResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code updating Clickhouse database",
			fmt.Sprintf("Received %s updating Clickhouse database: uuid=%s. Details: %s", updateResp.Status(), uuid, updateResp.Body))
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

func (r *clickhouseDatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state clickhouseDatabaseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Clickhouse database", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	deleteResp, err := r.client.DeleteDatabaseByUuidWithResponse(ctx, state.Uuid.ValueString(), &api.DeleteDatabaseByUuidParams{
		DeleteConfigurations:    types.BoolValue(true).ValueBoolPointer(),
		DeleteVolumes:           types.BoolValue(true).ValueBoolPointer(),
		DockerCleanup:           types.BoolValue(true).ValueBoolPointer(),
		DeleteConnectedNetworks: types.BoolValue(false).ValueBoolPointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Clickhouse database, got error: %s", err))
		return
	}

	if deleteResp.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code deleting Clickhouse database",
			fmt.Sprintf("Received %s deleting Clickhouse database: %s. Details: %s", deleteResp.Status(), state, deleteResp.Body))
		return
	}
}

func (r *clickhouseDatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *clickhouseDatabaseResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *clickhouseDatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() || plan == nil || state == nil {
		return
	}

	// If the username or password changes, the internal URL will change
	if !(plan.ClickhouseAdminUser.Equal(state.ClickhouseAdminUser) &&
		plan.ClickhouseAdminPassword.Equal(state.ClickhouseAdminPassword)) {
		plan.InternalDbUrl = types.StringUnknown()
		resp.Plan.Set(ctx, &plan)
	}
}

// MARK: Helper functions

func (r *clickhouseDatabaseResource) ReadFromAPI(
	ctx context.Context,
	diags *diag.Diagnostics,
	uuid string,
	state clickhouseDatabaseResourceModel,
) (clickhouseDatabaseResourceModel, bool) {
	readResp, err := r.client.GetDatabaseByUuidWithResponse(ctx, uuid)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading Clickhouse database: uuid=%s", uuid),
			err.Error(),
		)
		return clickhouseDatabaseResourceModel{}, false
	}

	if readResp.StatusCode() == http.StatusNotFound {
		return clickhouseDatabaseResourceModel{}, false
	}

	if readResp.StatusCode() != http.StatusOK {
		diags.AddError(
			"Unexpected HTTP status code reading Clickhouse database",
			fmt.Sprintf("Received %s for Clickhouse database: uuid=%s. Details: %s", readResp.Status(), uuid, readResp.Body))
		return clickhouseDatabaseResourceModel{}, false
	}

	result, err := clickhouseDatabaseResourceModel{}.FromAPI(readResp.JSON200, state)
	if err != nil {
		diags.AddError("Error converting API response to model", err.Error())
		return clickhouseDatabaseResourceModel{}, false
	}

	return result, true
}
