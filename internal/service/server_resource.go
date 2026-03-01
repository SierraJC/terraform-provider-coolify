package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/flatten"
	"terraform-provider-coolify/internal/provider/generated/resource_server"
	"terraform-provider-coolify/internal/provider/util"
	"terraform-provider-coolify/internal/wait"
)

var (
	_ resource.Resource                = &serverResource{}
	_ resource.ResourceWithConfigure   = &serverResource{}
	_ resource.ResourceWithImportState = &serverResource{}
)

func NewServerResource() resource.Resource {
	return &serverResource{}
}

type serverResource struct {
	client *api.ClientWithResponses
}

// serverResourceModel wraps the generated model with additional fields
type serverResourceModel struct {
	resource_server.ServerModel
	WaitForValidation types.Bool     `tfsdk:"wait_for_validation"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

func (r *serverResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *serverResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_server.ServerResourceSchema(ctx)
	resp.Schema.Description = "Create, read, update, and delete a Coolify server resource." +
		"\n**NOTE:** This resource is not fully implemented and may not work as expected because the Coolify API is incomplete."

	requiredAttrs := []string{"name", "private_key_uuid", "ip", "instant_validate"}
	for _, attr := range requiredAttrs {
		makeResourceAttributeRequired(resp.Schema.Attributes, attr)
	}

	validateNonEmptyStrings := []string{"description", "user", ""}
	for _, attr := range validateNonEmptyStrings {
		makeResourceAttributeNonEmpty(resp.Schema.Attributes, attr)
	}

	// Add wait_for_validation attribute
	resp.Schema.Attributes["wait_for_validation"] = schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Default:     booldefault.StaticBool(false),
		Description: "Wait for server validation (is_reachable && is_usable) to complete during creation. Defaults to false.",
	}

	// Add timeouts block
	resp.Schema.Blocks = map[string]schema.Block{
		"timeouts": timeouts.Block(ctx, timeouts.Opts{
			Create: true,
		}),
	}
}

func (r *serverResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.client, resp)
}

func (r *serverResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serverResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating server", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})
	createResp, err := r.client.CreateServerWithResponse(ctx, api.CreateServerJSONRequestBody{
		Description:     plan.Description.ValueStringPointer(),
		Name:            plan.Name.ValueStringPointer(),
		InstantValidate: plan.InstantValidate.ValueBoolPointer(),
		Ip:              plan.Ip.ValueStringPointer(),
		IsBuildServer:   plan.IsBuildServer.ValueBoolPointer(),
		Port: func() *int {
			if plan.Port.IsUnknown() || plan.Port.IsNull() {
				return nil
			}
			value := int(*plan.Port.ValueInt64Pointer())
			return &value
		}(),
		PrivateKeyUuid: plan.PrivateKeyUuid.ValueStringPointer(),
		User:           plan.User.ValueStringPointer(),
	})

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

	data, _ := r.ReadFromAPI(ctx, &resp.Diagnostics, *createResp.JSON201.Uuid)
	data.WaitForValidation = plan.WaitForValidation
	data.Timeouts = plan.Timeouts

	// Check if we should wait for validation
	if !plan.WaitForValidation.IsNull() && plan.WaitForValidation.ValueBool() {
		// Get timeout from timeouts block (default 3 minutes)
		var timeoutsValue timeouts.Value
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("timeouts"), &timeoutsValue)...)
		if resp.Diagnostics.HasError() {
			return
		}

		createTimeout, diags := timeoutsValue.Create(ctx, 3*time.Minute)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, createTimeout)
		defer cancel()

		serverUuid := *createResp.JSON201.Uuid
		tflog.Debug(ctx, "Waiting for server validation", map[string]interface{}{
			"uuid":    serverUuid,
			"timeout": createTimeout.String(),
		})

		err = wait.WaitForCondition(timeoutCtx, 5*time.Second, func() (bool, error) {
			// Re-read server state
			currentData, found := r.ReadFromAPI(timeoutCtx, &resp.Diagnostics, serverUuid)
			if !found {
				return false, fmt.Errorf("server not found during validation wait")
			}

			// Check for validation errors first
			if !currentData.ValidationLogs.IsNull() && currentData.ValidationLogs.ValueString() != "" {
				logContent := currentData.ValidationLogs.ValueString()
				if strings.Contains(strings.ToLower(logContent), "error") || strings.Contains(strings.ToLower(logContent), "failed") {
					return false, fmt.Errorf("server validation failed: %s", logContent)
				}
			}

			// Check settings flags
			isReachable := !currentData.Settings.IsReachable.IsNull() && currentData.Settings.IsReachable.ValueBool()
			isUsable := !currentData.Settings.IsUsable.IsNull() && currentData.Settings.IsUsable.ValueBool()

			if isReachable && isUsable {
				data = currentData // Update data with validated state
				return true, nil
			}

			return false, nil // Keep waiting
		})

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				resp.Diagnostics.AddError(
					"Server validation timeout",
					fmt.Sprintf("Server %s did not become validated within %s", serverUuid, createTimeout),
				)
			} else if errors.Is(err, context.Canceled) {
				resp.Diagnostics.AddError(
					"Server validation cancelled",
					fmt.Sprintf("Server %s validation was cancelled: %s", serverUuid, err),
				)
			} else {
				resp.Diagnostics.AddError(
					"Server validation failed",
					fmt.Sprintf("Server %s validation failed: %s", serverUuid, err),
				)
			}
			return
		}

		tflog.Debug(ctx, "Server validation completed", map[string]interface{}{
			"uuid": serverUuid,
		})
	}

	r.copyMissingAttributes(&plan, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
func (r *serverResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serverResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading server", map[string]interface{}{
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

func (r *serverResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan serverResourceModel
	var state serverResourceModel

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
	tflog.Debug(ctx, "Updating server", map[string]interface{}{
		"uuid": uuid,
	})
	updateResp, err := r.client.UpdateServerByUuidWithResponse(ctx, uuid, api.UpdateServerByUuidJSONRequestBody{
		Description:     plan.Description.ValueStringPointer(),
		Name:            plan.Name.ValueStringPointer(),
		InstantValidate: plan.InstantValidate.ValueBoolPointer(),
		Ip:              plan.Ip.ValueStringPointer(),
		IsBuildServer:   plan.IsBuildServer.ValueBoolPointer(),
		Port: func() *int { // todo: make a reusable fn for these inline conversions
			if plan.Port.IsUnknown() || plan.Port.IsNull() {
				return nil
			}
			value := int(*plan.Port.ValueInt64Pointer())
			return &value
		}(),
		PrivateKeyUuid: plan.PrivateKeyUuid.ValueStringPointer(),
		User: func() *string {
			if plan.User.IsUnknown() {
				return nil
			}
			return plan.User.ValueStringPointer()
		}(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating server: uuid=%s", uuid),
			err.Error(),
		)
		return
	}

	if updateResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code updating server",
			fmt.Sprintf("Received %s updating server: uuid=%s. Details: %s", updateResp.Status(), uuid, updateResp.Body))
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

func (r *serverResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serverResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting server", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	deleteResp, err := r.client.DeleteServerByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete server, got error: %s", err))
		return
	}

	if deleteResp.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code deleting server",
			fmt.Sprintf("Received %s deleting server: %s. Details: %s", deleteResp.Status(), state, deleteResp.Body))
		return
	}
}

func (r *serverResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *serverResource) copyMissingAttributes(
	plan *serverResourceModel,
	data *serverResourceModel,
) {
	// Values that are not returned in API response
	data.InstantValidate = plan.InstantValidate
	data.PrivateKeyUuid = plan.PrivateKeyUuid
	data.WaitForValidation = plan.WaitForValidation
	data.Timeouts = plan.Timeouts

	if plan.PrivateKeyUuid.IsNull() {
		data.PrivateKeyUuid = types.StringValue("")
	}

	// Values that are incorrectly mapped in API
	data.Id = data.Settings.ServerId
}

func (r *serverResource) ReadFromAPI(
	ctx context.Context,
	diags *diag.Diagnostics,
	uuid string,
) (serverResourceModel, bool) {
	readResp, err := r.client.GetServerByUuidWithResponse(ctx, uuid)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading server: uuid=%s", uuid),
			err.Error(),
		)
		return serverResourceModel{}, false
	}

	if readResp.StatusCode() == http.StatusNotFound {
		return serverResourceModel{}, false
	}

	if readResp.StatusCode() != http.StatusOK {
		diags.AddError(
			"Unexpected HTTP status code reading server",
			fmt.Sprintf("Received %s for server: uuid=%s. Details: %s", readResp.Status(), uuid, readResp.Body))
		return serverResourceModel{}, false
	}

	return serverResourceModel{
		ServerModel: r.ApiToModel(ctx, diags, readResp.JSON200),
	}, true
}

func (r *serverResource) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *api.Server,
) resource_server.ServerModel {
	settings := resource_server.NewSettingsValueMust(
		resource_server.SettingsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"concurrent_builds":                     flatten.Int64(response.Settings.ConcurrentBuilds),
			"created_at":                            flatten.String(response.Settings.CreatedAt),
			"delete_unused_networks":                flatten.Bool(response.Settings.DeleteUnusedNetworks),
			"delete_unused_volumes":                 flatten.Bool(response.Settings.DeleteUnusedVolumes),
			"docker_cleanup_frequency":              flatten.String(response.Settings.DockerCleanupFrequency),
			"docker_cleanup_threshold":              flatten.Int64(response.Settings.DockerCleanupThreshold),
			"dynamic_timeout":                       flatten.Int64(response.Settings.DynamicTimeout),
			"force_disabled":                        flatten.Bool(response.Settings.ForceDisabled),
			"force_server_cleanup":                  flatten.Bool(response.Settings.ForceServerCleanup),
			"id":                                    flatten.Int64(response.Settings.Id),
			"is_build_server":                       flatten.Bool(response.Settings.IsBuildServer),
			"is_cloudflare_tunnel":                  flatten.Bool(response.Settings.IsCloudflareTunnel),
			"is_jump_server":                        flatten.Bool(response.Settings.IsJumpServer),
			"is_logdrain_axiom_enabled":             flatten.Bool(response.Settings.IsLogdrainAxiomEnabled),
			"is_logdrain_custom_enabled":            flatten.Bool(response.Settings.IsLogdrainCustomEnabled),
			"is_logdrain_highlight_enabled":         flatten.Bool(response.Settings.IsLogdrainHighlightEnabled),
			"is_logdrain_newrelic_enabled":          flatten.Bool(response.Settings.IsLogdrainNewrelicEnabled),
			"is_metrics_enabled":                    flatten.Bool(response.Settings.IsMetricsEnabled),
			"is_reachable":                          flatten.Bool(response.Settings.IsReachable),
			"is_sentinel_enabled":                   flatten.Bool(response.Settings.IsSentinelEnabled),
			"is_swarm_manager":                      flatten.Bool(response.Settings.IsSwarmManager),
			"is_swarm_worker":                       flatten.Bool(response.Settings.IsSwarmWorker),
			"is_usable":                             flatten.Bool(response.Settings.IsUsable),
			"logdrain_axiom_api_key":                flatten.String(response.Settings.LogdrainAxiomApiKey),
			"logdrain_axiom_dataset_name":           flatten.String(response.Settings.LogdrainAxiomDatasetName),
			"logdrain_custom_config":                flatten.String(response.Settings.LogdrainCustomConfig),
			"logdrain_custom_config_parser":         flatten.String(response.Settings.LogdrainCustomConfigParser),
			"logdrain_highlight_project_id":         flatten.String(response.Settings.LogdrainHighlightProjectId),
			"logdrain_newrelic_base_uri":            flatten.String(response.Settings.LogdrainNewrelicBaseUri),
			"logdrain_newrelic_license_key":         flatten.String(response.Settings.LogdrainNewrelicLicenseKey),
			"sentinel_metrics_history_days":         flatten.Int64(response.Settings.SentinelMetricsHistoryDays),
			"sentinel_metrics_refresh_rate_seconds": flatten.Int64(response.Settings.SentinelMetricsRefreshRateSeconds),
			"sentinel_token":                        flatten.String(response.Settings.SentinelToken),
			"server_id":                             flatten.Int64(response.Settings.ServerId),
			"updated_at":                            flatten.String(response.Settings.UpdatedAt),
			"wildcard_domain":                       flatten.String(response.Settings.WildcardDomain),
		},
	)

	return resource_server.ServerModel{
		Description:                   flatten.String(response.Description),
		HighDiskUsageNotificationSent: flatten.Bool(response.HighDiskUsageNotificationSent), // missing
		Id:                            flatten.Int64(response.Id),
		Ip:                            flatten.String(response.Ip),
		IsBuildServer:                 flatten.Bool(response.Settings.IsBuildServer),
		LogDrainNotificationSent:      flatten.Bool(response.LogDrainNotificationSent),
		Name:                          flatten.String(response.Name),
		Port:                          flatten.Int64(response.Port),
		SwarmCluster:                  flatten.String(response.SwarmCluster),
		UnreachableCount:              flatten.Int64(response.UnreachableCount),
		UnreachableNotificationSent:   flatten.Bool(response.UnreachableNotificationSent),
		User:                          flatten.String(response.User),
		Uuid:                          flatten.String(response.Uuid),
		ValidationLogs:                flatten.String(response.ValidationLogs),

		// Proxy:                         resource_server.NewProxyValueUnknown(),
		ProxyType:       flatten.String((*string)(response.ProxyType)), // enum value
		PrivateKeyUuid:  types.StringUnknown(),
		InstantValidate: types.BoolUnknown(),
		Settings:        settings,
	}
}
