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

	// todo: Mark sensitive attributes
	resp.Schema.Blocks = map[string]schema.Block{
		"filter": filter.CreateDatasourceFilter(applicationsFilterNames),
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
			"BaseDirectory":                   flatten.String(application.BaseDirectory),
			"BuildCommand":                    flatten.String(application.BuildCommand),
			"BuildPack":                       flatten.String((*string)(application.BuildPack)), // enum value
			"ComposeParsingVersion":           flatten.String(application.ComposeParsingVersion),
			"ConfigHash":                      flatten.String(application.ConfigHash),
			"CreatedAt":                       flatten.Time(application.CreatedAt),
			"CustomDockerRunOptions":          flatten.String(application.CustomDockerRunOptions),
			"CustomHealthcheckFound":          flatten.Bool(application.CustomHealthcheckFound),
			"CustomLabels":                    flatten.String(application.CustomLabels),
			"CustomNginxConfiguration":        flatten.String(application.CustomNginxConfiguration),
			"DeletedAt":                       flatten.Time(application.DeletedAt),
			"Description":                     flatten.String(application.Description),
			"DestinationId":                   flatten.Int64(application.DestinationId),
			"DestinationType":                 flatten.String(application.DestinationType),
			"DockerCompose":                   flatten.String(application.DockerCompose),
			"DockerComposeCustomBuildCommand": flatten.String(application.DockerComposeCustomBuildCommand),
			"DockerComposeCustomStartCommand": flatten.String(application.DockerComposeCustomStartCommand),
			"DockerComposeDomains":            flatten.String(application.DockerComposeDomains),
			"DockerComposeLocation":           flatten.String(application.DockerComposeLocation),
			"DockerComposeRaw":                flatten.String(application.DockerComposeRaw),
			"DockerRegistryImageName":         flatten.String(application.DockerRegistryImageName),
			"DockerRegistryImageTag":          flatten.String(application.DockerRegistryImageTag),
			"Dockerfile":                      flatten.String(application.Dockerfile),
			"DockerfileLocation":              flatten.String(application.DockerfileLocation),
			"DockerfileTargetBuild":           flatten.String(application.DockerfileTargetBuild),
			"EnvironmentId":                   flatten.Int64(application.EnvironmentId),
			"Fqdn":                            flatten.String(application.Fqdn),
			"GitBranch":                       flatten.String(application.GitBranch),
			"GitCommitSha":                    flatten.String(application.GitCommitSha),
			"GitFullUrl":                      flatten.String(application.GitFullUrl),
			"GitRepository":                   flatten.String(application.GitRepository),
			"HealthCheckEnabled":              flatten.Bool(application.HealthCheckEnabled),
			"HealthCheckHost":                 flatten.String(application.HealthCheckHost),
			"HealthCheckInterval":             flatten.Int64(application.HealthCheckInterval),
			"HealthCheckMethod":               flatten.String(application.HealthCheckMethod),
			"HealthCheckPath":                 flatten.String(application.HealthCheckPath),
			"HealthCheckPort":                 flatten.String(application.HealthCheckPort),
			"HealthCheckResponseText":         flatten.String(application.HealthCheckResponseText),
			"HealthCheckRetries":              flatten.Int64(application.HealthCheckRetries),
			"HealthCheckReturnCode":           flatten.Int64(application.HealthCheckReturnCode),
			"HealthCheckScheme":               flatten.String(application.HealthCheckScheme),
			"HealthCheckStartPeriod":          flatten.Int64(application.HealthCheckStartPeriod),
			"HealthCheckTimeout":              flatten.Int64(application.HealthCheckTimeout),
			"Id":                              flatten.Int64(application.Id),
			"InstallCommand":                  flatten.String(application.InstallCommand),
			"LimitsCpuShares":                 flatten.Int64(application.LimitsCpuShares),
			"LimitsCpus":                      flatten.String(application.LimitsCpus),
			"LimitsCpuset":                    flatten.String(application.LimitsCpuset),
			"LimitsMemory":                    flatten.String(application.LimitsMemory),
			"LimitsMemoryReservation":         flatten.String(application.LimitsMemoryReservation),
			"LimitsMemorySwap":                flatten.String(application.LimitsMemorySwap),
			"LimitsMemorySwappiness":          flatten.Int64(application.LimitsMemorySwappiness),
			"ManualWebhookSecretBitbucket":    flatten.String(application.ManualWebhookSecretBitbucket),
			"ManualWebhookSecretGitea":        flatten.String(application.ManualWebhookSecretGitea),
			"ManualWebhookSecretGithub":       flatten.String(application.ManualWebhookSecretGithub),
			"ManualWebhookSecretGitlab":       flatten.String(application.ManualWebhookSecretGitlab),
			"Name":                            flatten.String(application.Name),
			"PortsExposes":                    flatten.String(application.PortsExposes),
			"PortsMappings":                   flatten.String(application.PortsMappings),
			"PostDeploymentCommand":           flatten.String(application.PostDeploymentCommand),
			"PostDeploymentCommandContainer":  flatten.String(application.PostDeploymentCommandContainer),
			"PreDeploymentCommand":            flatten.String(application.PreDeploymentCommand),
			"PreDeploymentCommandContainer":   flatten.String(application.PreDeploymentCommandContainer),
			"PreviewUrlTemplate":              flatten.String(application.PreviewUrlTemplate),
			"PrivateKeyId":                    flatten.Int64(application.PrivateKeyId),
			"PublishDirectory":                flatten.String(application.PublishDirectory),
			"Redirect":                        flatten.String((*string)(application.Redirect)), // enum value
			"RepositoryProjectId":             flatten.Int64(application.RepositoryProjectId),
			"SourceId":                        flatten.Int64(application.SourceId),
			"StartCommand":                    flatten.String(application.StartCommand),
			"StaticImage":                     flatten.String(application.StaticImage),
			"Status":                          flatten.String(application.Status),
			"SwarmPlacementConstraints":       flatten.String(application.SwarmPlacementConstraints),
			"SwarmReplicas":                   flatten.Int64(application.SwarmReplicas),
			"UpdatedAt":                       flatten.Time(application.UpdatedAt),
			"Uuid":                            flatten.String(application.Uuid),
			"WatchPaths":                      flatten.String(application.WatchPaths),
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
