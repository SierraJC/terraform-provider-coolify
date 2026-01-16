package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/flatten"
	"terraform-provider-coolify/internal/provider/generated/datasource_server"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &serverDataSource{}
var _ datasource.DataSourceWithConfigure = &serverDataSource{}

func NewServerDataSource() datasource.DataSource {
	return &serverDataSource{}
}

type serverDataSource struct {
	client *api.ClientWithResponses
}

func (d *serverDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (d *serverDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_server.ServerDataSourceSchema(ctx)
	resp.Schema.Description = "Get a Coolify server by `uuid`."
}

func (d *serverDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.client, resp)
}

func (d *serverDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan datasource_server.ServerModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.client.GetServerByUuidWithResponse(ctx, plan.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading server", err.Error(),
		)
		return
	}

	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading server",
			fmt.Sprintf("Received %s for server. Details: %s", response.Status(), string(response.Body)),
		)
		return
	}

	state := d.ApiToModel(ctx, &resp.Diagnostics, response.JSON200)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *serverDataSource) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *api.Server,
) datasource_server.ServerModel {

	settings := datasource_server.NewSettingsValueMust(
		datasource_server.SettingsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"concurrent_builds":                     flatten.Int64(response.Settings.ConcurrentBuilds),
			"created_at":                            flatten.String(response.Settings.CreatedAt),
			"delete_unused_networks":                flatten.Bool(response.Settings.DeleteUnusedNetworks),
			"delete_unused_volumes":                 flatten.Bool(response.Settings.DeleteUnusedVolumes),
			"deployment_queue_limit":                flatten.Int64(response.Settings.DeploymentQueueLimit),
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
			"is_terminal_enabled":                   flatten.Bool(response.Settings.IsTerminalEnabled),
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

	return datasource_server.ServerModel{
		Description:                   flatten.String(response.Description),
		HighDiskUsageNotificationSent: flatten.Bool(response.HighDiskUsageNotificationSent),
		Id:                            flatten.Int64(response.Settings.ServerId), // TODO: this should be `id` on root object, upstream spec is wrong
		Ip:                            flatten.String(response.Ip),
		LogDrainNotificationSent:      flatten.Bool(response.LogDrainNotificationSent),
		Name:                          flatten.String(response.Name),
		Port:                          flatten.Int64(response.Port),
		ProxyType:                     flatten.String((*string)(response.ProxyType)), // enum value
		Settings:                      settings,
		SwarmCluster:                  flatten.String(response.SwarmCluster),
		UnreachableCount:              flatten.Int64(response.UnreachableCount),
		UnreachableNotificationSent:   flatten.Bool(response.UnreachableNotificationSent),
		User:                          flatten.String(response.User),
		Uuid:                          flatten.String(response.Uuid),
		ValidationLogs:                flatten.String(response.ValidationLogs),
	}
}
