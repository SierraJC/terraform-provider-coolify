package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/filter"
	"terraform-provider-coolify/internal/flatten"
	"terraform-provider-coolify/internal/provider/generated/datasource_servers"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &serversDataSource{}
var _ datasource.DataSourceWithConfigure = &serversDataSource{}

func NewServersDataSource() datasource.DataSource {
	return &serversDataSource{}
}

type serversDataSource struct {
	client *api.ClientWithResponses
}

type serversDataSourceWithFilterModel struct {
	datasource_servers.ServersModel
	Filter []filter.BlockModel `tfsdk:"filter"`
}

var serversFilterNames = []string{"id", "uuid", "user", "ip", "name", "description"}

func (d *serversDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servers"
}

func (d *serversDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_servers.ServersDataSourceSchema(ctx)
	resp.Schema.Description = "Get a list of Coolify servers."

	resp.Schema.Blocks = map[string]schema.Block{
		"filter": filter.CreateDatasourceFilter(serversFilterNames),
	}
}
func (d *serversDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.client, resp)
}

func (d *serversDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan serversDataSourceWithFilterModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	listResponse, err := d.client.ListServersWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading servers", err.Error(),
		)
		return
	}

	if listResponse.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading servers",
			fmt.Sprintf("Received %s for servers. Details: %s", listResponse.Status(), listResponse.Body),
		)
		return
	}

	state, diag := d.apiToModel(ctx, listResponse.JSON200, plan.Filter)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *serversDataSource) apiToModel(
	ctx context.Context,
	response *[]api.Server,
	filters []filter.BlockModel,
) (serversDataSourceWithFilterModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var elements []attr.Value

	for _, sv := range *response {
		settings, diag := datasource_servers.NewSettingsValueMust(
			datasource_servers.SettingsValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"concurrent_builds":                     flatten.Int64(sv.Settings.ConcurrentBuilds),
				"created_at":                            flatten.String(sv.Settings.CreatedAt),
				"delete_unused_networks":                flatten.Bool(sv.Settings.DeleteUnusedNetworks),
				"delete_unused_volumes":                 flatten.Bool(sv.Settings.DeleteUnusedVolumes),
				"deployment_queue_limit":                flatten.Int64(sv.Settings.DeploymentQueueLimit),
				"docker_cleanup_frequency":              flatten.String(sv.Settings.DockerCleanupFrequency),
				"docker_cleanup_threshold":              flatten.Int64(sv.Settings.DockerCleanupThreshold),
				"dynamic_timeout":                       flatten.Int64(sv.Settings.DynamicTimeout),
				"force_disabled":                        flatten.Bool(sv.Settings.ForceDisabled),
				"force_server_cleanup":                  flatten.Bool(sv.Settings.ForceServerCleanup),
				"id":                                    flatten.Int64(sv.Settings.Id),
				"is_build_server":                       flatten.Bool(sv.Settings.IsBuildServer),
				"is_cloudflare_tunnel":                  flatten.Bool(sv.Settings.IsCloudflareTunnel),
				"is_jump_server":                        flatten.Bool(sv.Settings.IsJumpServer),
				"is_logdrain_axiom_enabled":             flatten.Bool(sv.Settings.IsLogdrainAxiomEnabled),
				"is_logdrain_custom_enabled":            flatten.Bool(sv.Settings.IsLogdrainCustomEnabled),
				"is_logdrain_highlight_enabled":         flatten.Bool(sv.Settings.IsLogdrainHighlightEnabled),
				"is_logdrain_newrelic_enabled":          flatten.Bool(sv.Settings.IsLogdrainNewrelicEnabled),
				"is_metrics_enabled":                    flatten.Bool(sv.Settings.IsMetricsEnabled),
				"is_reachable":                          flatten.Bool(sv.Settings.IsReachable),
				"is_sentinel_enabled":                   flatten.Bool(sv.Settings.IsSentinelEnabled),
				"is_swarm_manager":                      flatten.Bool(sv.Settings.IsSwarmManager),
				"is_swarm_worker":                       flatten.Bool(sv.Settings.IsSwarmWorker),
				"is_terminal_enabled":                   flatten.Bool(sv.Settings.IsTerminalEnabled),
				"is_usable":                             flatten.Bool(sv.Settings.IsUsable),
				"logdrain_axiom_api_key":                flatten.String(sv.Settings.LogdrainAxiomApiKey),
				"logdrain_axiom_dataset_name":           flatten.String(sv.Settings.LogdrainAxiomDatasetName),
				"logdrain_custom_config":                flatten.String(sv.Settings.LogdrainCustomConfig),
				"logdrain_custom_config_parser":         flatten.String(sv.Settings.LogdrainCustomConfigParser),
				"logdrain_highlight_project_id":         flatten.String(sv.Settings.LogdrainHighlightProjectId),
				"logdrain_newrelic_base_uri":            flatten.String(sv.Settings.LogdrainNewrelicBaseUri),
				"logdrain_newrelic_license_key":         flatten.String(sv.Settings.LogdrainNewrelicLicenseKey),
				"sentinel_metrics_history_days":         flatten.Int64(sv.Settings.SentinelMetricsHistoryDays),
				"sentinel_metrics_refresh_rate_seconds": flatten.Int64(sv.Settings.SentinelMetricsRefreshRateSeconds),
				"sentinel_token":                        flatten.String(sv.Settings.SentinelToken),
				"server_id":                             flatten.Int64(sv.Settings.ServerId),
				"updated_at":                            flatten.String(sv.Settings.UpdatedAt),
				"wildcard_domain":                       flatten.String(sv.Settings.WildcardDomain),
			},
		).ToObjectValue(ctx)
		diags.Append(diag...)

		attributes := map[string]attr.Value{
			"description":                       flatten.String(sv.Description),
			"high_disk_usage_notification_sent": flatten.Bool(sv.HighDiskUsageNotificationSent),
			"id":                                flatten.Int64(sv.Settings.ServerId), // TODO: this should be `id` on root object, upstream spec is wrong
			"ip":                                flatten.String(sv.Ip),
			"log_drain_notification_sent":       flatten.Bool(sv.LogDrainNotificationSent),
			"name":                              flatten.String(sv.Name),
			"port":                              flatten.Int64(sv.Port),
			"settings":                          settings,
			"swarm_cluster":                     flatten.String(sv.SwarmCluster),
			"unreachable_count":                 flatten.Int64(sv.UnreachableCount),
			"unreachable_notification_sent":     flatten.Bool(sv.UnreachableNotificationSent),
			"proxy_type":                        flatten.String((*string)(sv.ProxyType)), // enum value
			"user":                              flatten.String(sv.User),
			"uuid":                              flatten.String(sv.Uuid),
			"validation_logs":                   flatten.String(sv.ValidationLogs),
		}

		if !filter.OnAttributes(attributes, filters) {
			continue
		}

		data, diag := datasource_servers.NewServersValue(
			datasource_servers.ServersValue{}.AttributeTypes(ctx),
			attributes)
		diags.Append(diag...)
		elements = append(elements, data)
	}

	dataSet, diag := types.SetValue(datasource_servers.ServersValue{}.Type(ctx), elements)
	diags.Append(diag...)

	return serversDataSourceWithFilterModel{
		ServersModel: datasource_servers.ServersModel{
			Servers: dataSet,
		},
		Filter: filters,
	}, diags
}
