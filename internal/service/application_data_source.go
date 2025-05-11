package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/flatten"
	"terraform-provider-coolify/internal/provider/generated/datasource_application"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &applicationDataSource{}
var _ datasource.DataSourceWithConfigure = &applicationDataSource{}

func NewApplicationDataSource() datasource.DataSource {
	return &applicationDataSource{}
}

type applicationDataSource struct {
	client *api.ClientWithResponses
}

func (d *applicationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (d *applicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_application.ApplicationDataSourceSchema(ctx)
	resp.Schema.Description = "Get a Coolify application by `uuid`."

	// Mark sensitive attributes
	sensitiveAttrs := []string{"manual_webhook_secret_bitbucket", "manual_webhook_secret_gitea", "manual_webhook_secret_github", "manual_webhook_secret_gitlab", "http_basic_auth_password"}
	for _, attr := range sensitiveAttrs {
		makeDataSourceAttributeSensitive(resp.Schema.Attributes, attr)
	}
}

func (d *applicationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.client, resp)
}

func (d *applicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan datasource_application.ApplicationModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	applicationResp, err := d.client.GetApplicationByUuidWithResponse(ctx, plan.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading application", err.Error(),
		)
		return
	}

	if applicationResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading application",
			fmt.Sprintf("Received %s for application. Details: %s", applicationResp.Status(), string(applicationResp.Body)),
		)
		return
	}

	state := d.ApiToModel(ctx, &resp.Diagnostics, applicationResp.JSON200)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *applicationDataSource) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *api.Application,
) datasource_application.ApplicationModel {
	return datasource_application.ApplicationModel{
		BaseDirectory:                   flatten.String(response.BaseDirectory),
		BuildCommand:                    flatten.String(response.BuildCommand),
		BuildPack:                       flatten.String((*string)(response.BuildPack)), // enum value
		ComposeParsingVersion:           flatten.String(response.ComposeParsingVersion),
		ConfigHash:                      flatten.String(response.ConfigHash),
		CreatedAt:                       flatten.Time(response.CreatedAt),
		CustomDockerRunOptions:          flatten.String(response.CustomDockerRunOptions),
		CustomHealthcheckFound:          flatten.Bool(response.CustomHealthcheckFound),
		CustomLabels:                    flatten.String(response.CustomLabels),
		CustomNetworkAliases:            flatten.String(response.CustomNetworkAliases),
		CustomNginxConfiguration:        flatten.String(response.CustomNginxConfiguration),
		DeletedAt:                       flatten.Time(response.DeletedAt),
		Description:                     flatten.String(response.Description),
		DestinationId:                   flatten.Int64(response.DestinationId),
		DestinationType:                 flatten.String(response.DestinationType),
		DockerCompose:                   flatten.String(response.DockerCompose),
		DockerComposeCustomBuildCommand: flatten.String(response.DockerComposeCustomBuildCommand),
		DockerComposeCustomStartCommand: flatten.String(response.DockerComposeCustomStartCommand),
		DockerComposeDomains:            flatten.String(response.DockerComposeDomains),
		DockerComposeLocation:           flatten.String(response.DockerComposeLocation),
		DockerComposeRaw:                flatten.String(response.DockerComposeRaw),
		DockerRegistryImageName:         flatten.String(response.DockerRegistryImageName),
		DockerRegistryImageTag:          flatten.String(response.DockerRegistryImageTag),
		Dockerfile:                      flatten.String(response.Dockerfile),
		DockerfileLocation:              flatten.String(response.DockerfileLocation),
		DockerfileTargetBuild:           flatten.String(response.DockerfileTargetBuild),
		EnvironmentId:                   flatten.Int64(response.EnvironmentId),
		Fqdn:                            flatten.String(response.Fqdn),
		GitBranch:                       flatten.String(response.GitBranch),
		GitCommitSha:                    flatten.String(response.GitCommitSha),
		GitFullUrl:                      flatten.String(response.GitFullUrl),
		GitRepository:                   flatten.String(response.GitRepository),
		HealthCheckEnabled:              flatten.Bool(response.HealthCheckEnabled),
		HealthCheckHost:                 flatten.String(response.HealthCheckHost),
		HealthCheckInterval:             flatten.Int64(response.HealthCheckInterval),
		HealthCheckMethod:               flatten.String(response.HealthCheckMethod),
		HealthCheckPath:                 flatten.String(response.HealthCheckPath),
		HealthCheckPort:                 flatten.String(response.HealthCheckPort),
		HealthCheckResponseText:         flatten.String(response.HealthCheckResponseText),
		HealthCheckRetries:              flatten.Int64(response.HealthCheckRetries),
		HealthCheckReturnCode:           flatten.Int64(response.HealthCheckReturnCode),
		HealthCheckScheme:               flatten.String(response.HealthCheckScheme),
		HealthCheckStartPeriod:          flatten.Int64(response.HealthCheckStartPeriod),
		HealthCheckTimeout:              flatten.Int64(response.HealthCheckTimeout),
		HttpBasicAuthPassword:           flatten.String(response.HttpBasicAuthPassword),
		HttpBasicAuthUsername:           flatten.String(response.HttpBasicAuthUsername),
		Id:                              flatten.Int64(response.Id),
		InstallCommand:                  flatten.String(response.InstallCommand),
		IsHttpBasicAuthEnabled:          flatten.Bool(response.IsHttpBasicAuthEnabled),
		LimitsCpuShares:                 flatten.Int64(response.LimitsCpuShares),
		LimitsCpus:                      flatten.String(response.LimitsCpus),
		LimitsCpuset:                    flatten.String(response.LimitsCpuset),
		LimitsMemory:                    flatten.String(response.LimitsMemory),
		LimitsMemoryReservation:         flatten.String(response.LimitsMemoryReservation),
		LimitsMemorySwap:                flatten.String(response.LimitsMemorySwap),
		LimitsMemorySwappiness:          flatten.Int64(response.LimitsMemorySwappiness),
		ManualWebhookSecretBitbucket:    flatten.String(response.ManualWebhookSecretBitbucket),
		ManualWebhookSecretGitea:        flatten.String(response.ManualWebhookSecretGitea),
		ManualWebhookSecretGithub:       flatten.String(response.ManualWebhookSecretGithub),
		ManualWebhookSecretGitlab:       flatten.String(response.ManualWebhookSecretGitlab),
		Name:                            flatten.String(response.Name),
		PortsExposes:                    flatten.String(response.PortsExposes),
		PortsMappings:                   flatten.String(response.PortsMappings),
		PostDeploymentCommand:           flatten.String(response.PostDeploymentCommand),
		PostDeploymentCommandContainer:  flatten.String(response.PostDeploymentCommandContainer),
		PreDeploymentCommand:            flatten.String(response.PreDeploymentCommand),
		PreDeploymentCommandContainer:   flatten.String(response.PreDeploymentCommandContainer),
		PreviewUrlTemplate:              flatten.String(response.PreviewUrlTemplate),
		PrivateKeyId:                    flatten.Int64(response.PrivateKeyId),
		PublishDirectory:                flatten.String(response.PublishDirectory),
		Redirect:                        flatten.String((*string)(response.Redirect)), // enum value
		RepositoryProjectId:             flatten.Int64(response.RepositoryProjectId),
		SourceId:                        flatten.Int64(response.SourceId),
		StartCommand:                    flatten.String(response.StartCommand),
		StaticImage:                     flatten.String(response.StaticImage),
		Status:                          flatten.String(response.Status),
		SwarmPlacementConstraints:       flatten.String(response.SwarmPlacementConstraints),
		SwarmReplicas:                   flatten.Int64(response.SwarmReplicas),
		UpdatedAt:                       flatten.Time(response.UpdatedAt),
		Uuid:                            flatten.String(response.Uuid),
		WatchPaths:                      flatten.String(response.WatchPaths),
	}
}
