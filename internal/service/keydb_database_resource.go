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
	_ resource.Resource                = &keydbDatabaseResource{}
	_ resource.ResourceWithConfigure   = &keydbDatabaseResource{}
	_ resource.ResourceWithImportState = &keydbDatabaseResource{}
	_ resource.ResourceWithModifyPlan  = &keydbDatabaseResource{}
)

type keydbDatabaseResourceModel = keydbDatabaseModel

func NewKeyDBDatabaseResource() resource.Resource {
	return &keydbDatabaseResource{}
}

type keydbDatabaseResource struct {
	client *api.ClientWithResponses
}

func (r *keydbDatabaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keydb_database"
}

func (r *keydbDatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	commonSchema := commonDatabaseModel{}.CommonSchema(ctx)
	keydbSchema := schema.Schema{
		Description: "Create, read, update, and delete a Coolify database (KeyDB) resource.",
		Attributes: map[string]schema.Attribute{
			"keydb_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "KeyDB password",
			},
			"keydb_conf": schema.StringAttribute{
				Optional:    true,
				Description: "KeyDB conf",
			},
		},
	}

	resp.Schema = sutil.MergeResourceSchemas(commonSchema, keydbSchema)
}

func (r *keydbDatabaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.client, resp)
}

func (r *keydbDatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan keydbDatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating KeyDB database", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	createResp, err := r.client.CreateDatabaseKeydbWithResponse(ctx, api.CreateDatabaseKeydbJSONRequestBody{
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
		KeydbPassword:           plan.KeydbPassword.ValueStringPointer(),
		KeydbConf:               sutil.Base64EncodeAttr(plan.KeydbConf),
		ProjectUuid:             plan.ProjectUuid.ValueString(),
		PublicPort:              expand.Int64(plan.PublicPort),
		ServerUuid:              plan.ServerUuid.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating KeyDB database",
			err.Error(),
		)
		return
	}

	if createResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code creating KeyDB database",
			fmt.Sprintf("Received %s creating KeyDB database. Details: %s", createResp.Status(), createResp.Body),
		)
		return
	}

	data, _ := r.ReadFromAPI(ctx, &resp.Diagnostics, createResp.JSON201.Uuid, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keydbDatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state keydbDatabaseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading KeyDB database", map[string]interface{}{
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

func (r *keydbDatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan keydbDatabaseResourceModel
	var state keydbDatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := plan.Uuid.ValueString()

	tflog.Debug(ctx, "Updating KeyDB database", map[string]interface{}{
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
		KeydbPassword:           plan.KeydbPassword.ValueStringPointer(),
		KeydbConf:               sutil.Base64EncodeAttr(plan.KeydbConf),
		PublicPort:              expand.Int64(plan.PublicPort),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating KeyDB database: uuid=%s", uuid),
			err.Error(),
		)
		return
	}

	if updateResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code updating KeyDB database",
			fmt.Sprintf("Received %s updating KeyDB database: uuid=%s. Details: %s", updateResp.Status(), uuid, updateResp.Body))
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

func (r *keydbDatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state keydbDatabaseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting KeyDB database", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	deleteResp, err := r.client.DeleteDatabaseByUuidWithResponse(ctx, state.Uuid.ValueString(), &api.DeleteDatabaseByUuidParams{
		DeleteConfigurations:    types.BoolValue(true).ValueBoolPointer(),
		DeleteVolumes:           types.BoolValue(true).ValueBoolPointer(),
		DockerCleanup:           types.BoolValue(true).ValueBoolPointer(),
		DeleteConnectedNetworks: types.BoolValue(false).ValueBoolPointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete KeyDB database, got error: %s", err))
		return
	}

	if deleteResp.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code deleting KeyDB database",
			fmt.Sprintf("Received %s deleting KeyDB database: %s. Details: %s", deleteResp.Status(), state, deleteResp.Body))
		return
	}
}

func (r *keydbDatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *keydbDatabaseResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *keydbDatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() || plan == nil || state == nil {
		return
	}

	// If the password changes, the internal URL will change
	if !plan.KeydbPassword.Equal(state.KeydbPassword) {
		plan.InternalDbUrl = types.StringUnknown()
		resp.Plan.Set(ctx, &plan)
	}
}

// MARK: Helper functions

func (r *keydbDatabaseResource) ReadFromAPI(
	ctx context.Context,
	diags *diag.Diagnostics,
	uuid string,
	state keydbDatabaseResourceModel,
) (keydbDatabaseResourceModel, bool) {
	readResp, err := r.client.GetDatabaseByUuidWithResponse(ctx, uuid)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading KeyDB database: uuid=%s", uuid),
			err.Error(),
		)
		return keydbDatabaseResourceModel{}, false
	}

	if readResp.StatusCode() == http.StatusNotFound {
		return keydbDatabaseResourceModel{}, false
	}

	if readResp.StatusCode() != http.StatusOK {
		diags.AddError(
			"Unexpected HTTP status code reading KeyDB database",
			fmt.Sprintf("Received %s for KeyDB database: uuid=%s. Details: %s", readResp.Status(), uuid, readResp.Body))
		return keydbDatabaseResourceModel{}, false
	}

	result, err := keydbDatabaseResourceModel{}.FromAPI(readResp.JSON200, state)
	if err != nil {
		diags.AddError("Error converting API response to model", err.Error())
		return keydbDatabaseResourceModel{}, false
	}

	return result, true
}
