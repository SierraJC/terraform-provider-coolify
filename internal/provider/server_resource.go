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
	"terraform-provider-coolify/internal/provider/generated/resource_server"
	"terraform-provider-coolify/internal/provider/util"
)

var (
	_ resource.Resource                = (*serverResource)(nil)
	_ resource.ResourceWithConfigure   = (*serverResource)(nil)
	_ resource.ResourceWithImportState = (*serverResource)(nil)
)

func NewServerResource() resource.Resource {
	return &serverResource{}
}

type serverResource struct {
	providerData CoolifyProviderData
}

// type serverResourceModel struct {
// 	Id types.String `tfsdk:"id"`
// }

func (r *serverResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *serverResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_server.ServerResourceSchema(ctx)
	resp.Schema.Description = "Provides a Coolify server resource, which can be used to create, read, update, and delete servers."

}

func (r *serverResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.providerData, resp)
}

func (r *serverResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state resource_server.ServerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := api.CreateServerJSONRequestBody{
		Description:     state.Description.ValueStringPointer(),
		InstantValidate: state.InstantValidate.ValueBoolPointer(),
		Ip:              state.Ip.ValueString(),
		IsBuildServer:   state.IsBuildServer.ValueBoolPointer(),
		Name:            state.Name.ValueStringPointer(),
		Port:            int64ToIntPointer(state.Port),
		PrivateKeyUuid:  state.PrivateKeyUuid.ValueString(),
		User:            state.User.ValueStringPointer(),
	}

	tflog.Debug(ctx, "Creating server")
	createResp, err := r.providerData.client.CreateServerWithResponse(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating server",
			err.Error(),
		)
		return
	}

	if createResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code creating server",
			fmt.Sprintf("Received %s creating server. Details: %s", createResp.Status(), createResp.Body),
		)
		return
	}

	assignStr(createResp.JSON201.Uuid, &state.Uuid)

	// GET /servers/{uuid}
	readResp, err := r.providerData.client.GetServerByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading server: uuid=%s", state.Uuid.ValueString()),
			err.Error(),
		)
		return
	}

	if readResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading server",
			fmt.Sprintf("Received %s for server: uuid=%s. Details: %s", readResp.Status(), state.Uuid.ValueString(), readResp.Body))
		return
	}

	diags := setServerState(ctx, readResp.JSON200, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *serverResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resource_server.ServerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading server: uuid=%s", state.Uuid.ValueString()))
	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	// GET /servers/{uuid}
	readResp, err := r.providerData.client.GetServerByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading server: uuid=%s", state.Uuid.ValueString()),
			err.Error(),
		)
		return
	}

	if readResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading server",
			fmt.Sprintf("Received %s for server: uuid=%s. Details: %s", readResp.Status(), state.Uuid.ValueString(), readResp.Body))
		return
	}

	diags := setServerState(ctx, readResp.JSON200, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *serverResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state resource_server.ServerModel
	var plan resource_server.ServerModel

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

	body := api.UpdateServerByUuidJSONRequestBody{}

	tflog.Debug(ctx, fmt.Sprintf("Updating server: uuid=%s", state.Uuid.ValueString()))
	updateResp, err := r.providerData.client.UpdateServerByUuidWithResponse(ctx, state.Uuid.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating server: uuid=%s", state.Uuid.ValueString()),
			err.Error(),
		)
		return
	}

	if updateResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code updating server",
			fmt.Sprintf("Received %s updating server: uuid=%s. Details: %s", updateResp.Status(), state.Uuid.ValueString(), updateResp.Body))
		return
	}

	// assignStr(updateResp.JSON201.Uuid, &state.Uuid)

	// GET /servers/{uuid}
	readResp, err := r.providerData.client.GetServerByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading server: uuid=%s", state.Uuid.ValueString()),
			err.Error(),
		)
		return
	}

	if readResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading server",
			fmt.Sprintf("Received %s for server: uuid=%s. Details: %s", readResp.Status(), state.Uuid.ValueString(), readResp.Body))
		return
	}

	diags := setServerState(ctx, readResp.JSON200, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *serverResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resource_server.ServerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteResp, err := r.providerData.client.DeleteServerByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete server, got error: %s", err))
		return
	}

	if deleteResp.StatusCode() != http.StatusOK && deleteResp.StatusCode() != http.StatusNotFound {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code deleting server",
			fmt.Sprintf("Received %s deleting server: %s. Details: %s", deleteResp.Status(), state.Uuid.ValueString(), deleteResp.Body))
		return
	}

}

func (r *serverResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func setServerState(ctx context.Context, apiResp *api.Server, state *resource_server.ServerModel) (diags diag.Diagnostics) {
	assignStr(apiResp.Description, &state.Description)
	assignBool(apiResp.HighDiskUsageNotificationSent, &state.HighDiskUsageNotificationSent)
	assignInt(apiResp.Settings.Id, &state.Id)
	// state.InstantValidate = types.BoolNull()
	assignStr(apiResp.Ip, &state.Ip)
	assignBool(apiResp.Settings.IsBuildServer, &state.IsBuildServer)
	assignBool(apiResp.LogDrainNotificationSent, &state.LogDrainNotificationSent)
	assignStr(apiResp.Name, &state.Name)
	assignInt(apiResp.Port, &state.Port)
	// state.PrivateKeyUuid = types.StringNull()
	assignStr(apiResp.SwarmCluster, &state.SwarmCluster)
	assignInt(apiResp.UnreachableCount, &state.UnreachableCount)
	assignBool(apiResp.UnreachableNotificationSent, &state.UnreachableNotificationSent)
	assignStr(apiResp.User, &state.User)
	assignStr(apiResp.Uuid, &state.Uuid)
	assignStr(apiResp.ValidationLogs, &state.ValidationLogs)
	state.Settings = resource_server.NewSettingsValueNull()

	// todo: get proxy
	proxy, diags := getServerProxyState(ctx, apiResp.Proxy)
	if diags.HasError() {
		return diags
	}
	state.Proxy = proxy

	settings, diags := getServerSettingsState(ctx, apiResp.Settings)
	if diags.HasError() {
		return diags
	}
	state.Settings = settings

	return diags
}

func getServerProxyState(ctx context.Context, proxies *map[string]interface{}) (resource_server.ProxyValue, diag.Diagnostics) {
	return resource_server.NewProxyValueNull(), nil
}

func getServerSettingsState(ctx context.Context, settings *api.ServerSetting) (resource_server.SettingsValue, diag.Diagnostics) {
	// todo: make a function that will compare the result of this function with SettingsValue, to ensure no fields are missing
	return resource_server.NewSettingsValue(resource_server.SettingsValue{}.AttributeTypes(ctx), map[string]attr.Value{
		// "cleanup_after_percentage":      types.Int64PointerValue(intPointerToInt64Pointer(settings.CleanupAfterPercentage)),
		// "concurrent_builds":             types.Int64PointerValue(intPointerToInt64Pointer(settings.ConcurrentBuilds)),
		"created_at":                    types.StringPointerValue(settings.CreatedAt),
		"dynamic_timeout":               types.Int64PointerValue(intPointerToInt64Pointer(settings.DynamicTimeout)),
		"force_disabled":                types.BoolPointerValue(settings.ForceDisabled),
		"id":                            types.Int64PointerValue(intPointerToInt64Pointer(settings.Id)),
		"is_build_server":               types.BoolPointerValue(settings.IsBuildServer),
		"is_cloudflare_tunnel":          types.BoolPointerValue(settings.IsCloudflareTunnel),
		"is_jump_server":                types.BoolPointerValue(settings.IsJumpServer),
		"is_logdrain_axiom_enabled":     types.BoolPointerValue(settings.IsLogdrainAxiomEnabled),
		"is_logdrain_custom_enabled":    types.BoolPointerValue(settings.IsLogdrainCustomEnabled),
		"is_logdrain_highlight_enabled": types.BoolPointerValue(settings.IsLogdrainHighlightEnabled),
		"is_logdrain_newrelic_enabled":  types.BoolPointerValue(settings.IsLogdrainNewrelicEnabled),
		"is_metrics_enabled":            types.BoolPointerValue(settings.IsMetricsEnabled),
		"is_reachable":                  types.BoolPointerValue(settings.IsReachable),
		// "is_server_api_enabled":         types.BoolPointerValue(settings.IsServerApiEnabled),
		"is_swarm_manager":              types.BoolPointerValue(settings.IsSwarmManager),
		"is_swarm_worker":               types.BoolPointerValue(settings.IsSwarmWorker),
		"is_usable":                     types.BoolPointerValue(settings.IsUsable),
		"logdrain_axiom_api_key":        types.StringPointerValue(settings.LogdrainAxiomApiKey),
		"logdrain_axiom_dataset_name":   types.StringPointerValue(settings.LogdrainAxiomDatasetName),
		"logdrain_custom_config":        types.StringPointerValue(settings.LogdrainCustomConfig),
		"logdrain_custom_config_parser": types.StringPointerValue(settings.LogdrainCustomConfigParser),
		"logdrain_highlight_project_id": types.StringPointerValue(settings.LogdrainHighlightProjectId),
		"logdrain_newrelic_base_uri":    types.StringPointerValue(settings.LogdrainNewrelicBaseUri),
		"logdrain_newrelic_license_key": types.StringPointerValue(settings.LogdrainNewrelicLicenseKey),
		// "metrics_history_days":          types.Int64PointerValue(intPointerToInt64Pointer(settings.MetricsHistoryDays)),
		// "metrics_refresh_rate_seconds":  types.Int64PointerValue(intPointerToInt64Pointer(settings.MetricsRefreshRateSeconds)),
		// "metrics_token":                 types.StringPointerValue(settings.MetricsToken),
		"server_id":                     types.Int64PointerValue(intPointerToInt64Pointer(settings.ServerId)),
		"updated_at":                    types.StringPointerValue(settings.UpdatedAt),
		"wildcard_domain":               types.StringPointerValue(settings.WildcardDomain),
	})
}
