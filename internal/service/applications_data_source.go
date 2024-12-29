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
	"terraform-provider-coolify/internal/provider/generated/datasource_applications"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &applicationsDataSource{}
var _ datasource.DataSourceWithConfigure = &applicationsDataSource{}

func NewApplicationsDataSource() datasource.DataSource {
	return &applicationsDataSource{}
}

type applicationsDataSource struct {
	client *api.ClientWithResponses
}

type applicationsDataSourceWithFilterModel struct {
	datasource_applications.ApplicationsModel
	Filter []filter.BlockModel `tfsdk:"filter"`
}

var applicationsFilterNames = []string{"id", "uuid", "name", "description", "fqdn"}

func (d *applicationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_applications"
}

func (d *applicationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_applications.ApplicationsDataSourceSchema(ctx)
	resp.Schema.Description = "Get a list of Coolify applications."
	resp.Schema.Blocks = map[string]schema.Block{
		"filter": filter.CreateDatasourceFilter(applicationsFilterNames),
	}

	// Mark sensitive attributes
	sensitiveAttrs := []string{"manual_webhook_secret_bitbucket", "manual_webhook_secret_gitea", "manual_webhook_secret_github", "manual_webhook_secret_gitlab"}
	for _, attr := range sensitiveAttrs {
		makeDataSourceAttributeSensitive(
			resp.Schema.Attributes["applications"].(schema.SetNestedAttribute).NestedObject.Attributes,
			attr,
		)
	}
}

func (d *applicationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.client, resp)
}

func (d *applicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan applicationsDataSourceWithFilterModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	listResponse, err := d.client.ListApplicationsWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading applications", err.Error(),
		)
		return
	}

	if listResponse.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading applications",
			fmt.Sprintf("Received %s for applications. Details: %s", listResponse.Status(), string(listResponse.Body)),
		)
		return
	}

	state, diag := d.ApiToModel(ctx, listResponse.JSON200, plan.Filter)
	resp.Diagnostics.Append(diag...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *applicationsDataSource) ApiToModel(
	ctx context.Context,
	response *[]api.Application,
	filters []filter.BlockModel,
) (applicationsDataSourceWithFilterModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var applications []attr.Value

	for _, application := range *response {
		attributes := map[string]attr.Value{
			"base_directory":                      flatten.String(application.BaseDirectory),
			"build_command":                       flatten.String(application.BuildCommand),
			"build_pack":                          flatten.String((*string)(application.BuildPack)), // enum value
			"compose_parsing_version":             flatten.String(application.ComposeParsingVersion),
			"config_hash":                         flatten.String(application.ConfigHash),
			"created_at":                          flatten.Time(application.CreatedAt),
			"custom_docker_run_options":           flatten.String(application.CustomDockerRunOptions),
			"custom_healthcheck_found":            flatten.Bool(application.CustomHealthcheckFound),
			"custom_labels":                       flatten.String(application.CustomLabels),
			"custom_nginx_configuration":          flatten.String(application.CustomNginxConfiguration),
			"deleted_at":                          flatten.Time(application.DeletedAt),
			"description":                         flatten.String(application.Description),
			"destination_id":                      flatten.Int64(application.DestinationId),
			"destination_type":                    flatten.String(application.DestinationType),
			"docker_compose":                      flatten.String(application.DockerCompose),
			"docker_compose_custom_build_command": flatten.String(application.DockerComposeCustomBuildCommand),
			"docker_compose_custom_start_command": flatten.String(application.DockerComposeCustomStartCommand),
			"docker_compose_domains":              flatten.String(application.DockerComposeDomains),
			"docker_compose_location":             flatten.String(application.DockerComposeLocation),
			"docker_compose_raw":                  flatten.String(application.DockerComposeRaw),
			"docker_registry_image_name":          flatten.String(application.DockerRegistryImageName),
			"docker_registry_image_tag":           flatten.String(application.DockerRegistryImageTag),
			"dockerfile":                          flatten.String(application.Dockerfile),
			"dockerfile_location":                 flatten.String(application.DockerfileLocation),
			"dockerfile_target_build":             flatten.String(application.DockerfileTargetBuild),
			"environment_id":                      flatten.Int64(application.EnvironmentId),
			"fqdn":                                flatten.String(application.Fqdn),
			"git_branch":                          flatten.String(application.GitBranch),
			"git_commit_sha":                      flatten.String(application.GitCommitSha),
			"git_full_url":                        flatten.String(application.GitFullUrl),
			"git_repository":                      flatten.String(application.GitRepository),
			"health_check_enabled":                flatten.Bool(application.HealthCheckEnabled),
			"health_check_host":                   flatten.String(application.HealthCheckHost),
			"health_check_interval":               flatten.Int64(application.HealthCheckInterval),
			"health_check_method":                 flatten.String(application.HealthCheckMethod),
			"health_check_path":                   flatten.String(application.HealthCheckPath),
			"health_check_port":                   flatten.String(application.HealthCheckPort),
			"health_check_response_text":          flatten.String(application.HealthCheckResponseText),
			"health_check_retries":                flatten.Int64(application.HealthCheckRetries),
			"health_check_return_code":            flatten.Int64(application.HealthCheckReturnCode),
			"health_check_scheme":                 flatten.String(application.HealthCheckScheme),
			"health_check_start_period":           flatten.Int64(application.HealthCheckStartPeriod),
			"health_check_timeout":                flatten.Int64(application.HealthCheckTimeout),
			"id":                                  flatten.Int64(application.Id),
			"install_command":                     flatten.String(application.InstallCommand),
			"limits_cpu_shares":                   flatten.Int64(application.LimitsCpuShares),
			"limits_cpus":                         flatten.String(application.LimitsCpus),
			"limits_cpuset":                       flatten.String(application.LimitsCpuset),
			"limits_memory":                       flatten.String(application.LimitsMemory),
			"limits_memory_reservation":           flatten.String(application.LimitsMemoryReservation),
			"limits_memory_swap":                  flatten.String(application.LimitsMemorySwap),
			"limits_memory_swappiness":            flatten.Int64(application.LimitsMemorySwappiness),
			"manual_webhook_secret_bitbucket":     flatten.String(application.ManualWebhookSecretBitbucket),
			"manual_webhook_secret_gitea":         flatten.String(application.ManualWebhookSecretGitea),
			"manual_webhook_secret_github":        flatten.String(application.ManualWebhookSecretGithub),
			"manual_webhook_secret_gitlab":        flatten.String(application.ManualWebhookSecretGitlab),
			"name":                                flatten.String(application.Name),
			"ports_exposes":                       flatten.String(application.PortsExposes),
			"ports_mappings":                      flatten.String(application.PortsMappings),
			"post_deployment_command":             flatten.String(application.PostDeploymentCommand),
			"post_deployment_command_container":   flatten.String(application.PostDeploymentCommandContainer),
			"pre_deployment_command":              flatten.String(application.PreDeploymentCommand),
			"pre_deployment_command_container":    flatten.String(application.PreDeploymentCommandContainer),
			"preview_url_template":                flatten.String(application.PreviewUrlTemplate),
			"private_key_id":                      flatten.Int64(application.PrivateKeyId),
			"publish_directory":                   flatten.String(application.PublishDirectory),
			"redirect":                            flatten.String((*string)(application.Redirect)), // enum value
			"repository_project_id":               flatten.Int64(application.RepositoryProjectId),
			"source_id":                           flatten.Int64(application.SourceId),
			"start_command":                       flatten.String(application.StartCommand),
			"static_image":                        flatten.String(application.StaticImage),
			"status":                              flatten.String(application.Status),
			"swarm_placement_constraints":         flatten.String(application.SwarmPlacementConstraints),
			"swarm_replicas":                      flatten.Int64(application.SwarmReplicas),
			"updated_at":                          flatten.Time(application.UpdatedAt),
			"uuid":                                flatten.String(application.Uuid),
			"watch_paths":                         flatten.String(application.WatchPaths),
		}

		if !filter.OnAttributes(attributes, filters) {
			continue
		}

		data, diag := datasource_applications.NewApplicationsValue(
			datasource_applications.ApplicationsValue{}.AttributeTypes(ctx),
			attributes)
		diags.Append(diag...)
		applications = append(applications, data)
	}

	dataSet, diag := types.SetValue(datasource_applications.ApplicationsValue{}.Type(ctx), applications)
	diags.Append(diag...)

	return applicationsDataSourceWithFilterModel{
		ApplicationsModel: datasource_applications.ApplicationsModel{
			Applications: dataSet,
		},
		Filter: filters,
	}, diags
}
