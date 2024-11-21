package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/generated/datasource_application"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &applicationDataSource{}
var _ datasource.DataSourceWithConfigure = &applicationDataSource{}

func NewApplicationDataSource() datasource.DataSource {
	return &applicationDataSource{}
}

type applicationDataSource struct {
	providerData CoolifyProviderData
}

func (d *applicationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (d *applicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_application.ApplicationDataSourceSchema(ctx)
	resp.Schema.Description = "Get a Coolify application by `uuid`."

	// Mark sensitive attributes
	sensitiveAttrs := []string{"manual_webhook_secret_bitbucket", "manual_webhook_secret_gitea", "manual_webhook_secret_github", "manual_webhook_secret_gitlab"}
	for _, attr := range sensitiveAttrs {
		makeDataSourceAttributeSensitive(resp.Schema.Attributes, attr)
	}
}

func (d *applicationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.providerData, resp)
}

func (d *applicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan datasource_application.ApplicationModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	applicationResp, err := d.providerData.client.GetApplicationByUuidWithResponse(ctx, plan.Uuid.ValueString())
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
		BaseDirectory:                   optionalString(response.BaseDirectory),
		BuildCommand:                    optionalString(response.BuildCommand),
		BuildPack:                       optionalString((*string)(response.BuildPack)), // enum value
		ComposeParsingVersion:           optionalString(response.ComposeParsingVersion),
		ConfigHash:                      optionalString(response.ConfigHash),
		CreatedAt:                       optionalTime(response.CreatedAt),
		CustomDockerRunOptions:          optionalString(response.CustomDockerRunOptions),
		CustomHealthcheckFound:          optionalBool(response.CustomHealthcheckFound),
		CustomLabels:                    optionalString(response.CustomLabels),
		DeletedAt:                       optionalTime(response.DeletedAt),
		Description:                     optionalString(response.Description),
		DestinationId:                   optionalInt64(response.DestinationId),
		DestinationType:                 optionalString(response.DestinationType),
		DockerCompose:                   optionalString(response.DockerCompose),
		DockerComposeCustomBuildCommand: optionalString(response.DockerComposeCustomBuildCommand),
		DockerComposeCustomStartCommand: optionalString(response.DockerComposeCustomStartCommand),
		DockerComposeDomains:            optionalString(response.DockerComposeDomains),
		DockerComposeLocation:           optionalString(response.DockerComposeLocation),
		DockerComposeRaw:                optionalString(response.DockerComposeRaw),
		DockerRegistryImageName:         optionalString(response.DockerRegistryImageName),
		DockerRegistryImageTag:          optionalString(response.DockerRegistryImageTag),
		Dockerfile:                      optionalString(response.Dockerfile),
		DockerfileLocation:              optionalString(response.DockerfileLocation),
		DockerfileTargetBuild:           optionalString(response.DockerfileTargetBuild),
		EnvironmentId:                   optionalInt64(response.EnvironmentId),
		Fqdn:                            optionalString(response.Fqdn),
		GitBranch:                       optionalString(response.GitBranch),
		GitCommitSha:                    optionalString(response.GitCommitSha),
		GitFullUrl:                      optionalString(response.GitFullUrl),
		GitRepository:                   optionalString(response.GitRepository),
		HealthCheckEnabled:              optionalBool(response.HealthCheckEnabled),
		HealthCheckHost:                 optionalString(response.HealthCheckHost),
		HealthCheckInterval:             optionalInt64(response.HealthCheckInterval),
		HealthCheckMethod:               optionalString(response.HealthCheckMethod),
		HealthCheckPath:                 optionalString(response.HealthCheckPath),
		HealthCheckPort:                 optionalString(response.HealthCheckPort),
		HealthCheckResponseText:         optionalString(response.HealthCheckResponseText),
		HealthCheckRetries:              optionalInt64(response.HealthCheckRetries),
		HealthCheckReturnCode:           optionalInt64(response.HealthCheckReturnCode),
		HealthCheckScheme:               optionalString(response.HealthCheckScheme),
		HealthCheckStartPeriod:          optionalInt64(response.HealthCheckStartPeriod),
		HealthCheckTimeout:              optionalInt64(response.HealthCheckTimeout),
		Id:                              optionalInt64(response.Id),
		InstallCommand:                  optionalString(response.InstallCommand),
		LimitsCpuShares:                 optionalInt64(response.LimitsCpuShares),
		LimitsCpus:                      optionalString(response.LimitsCpus),
		LimitsCpuset:                    optionalString(response.LimitsCpuset),
		LimitsMemory:                    optionalString(response.LimitsMemory),
		LimitsMemoryReservation:         optionalString(response.LimitsMemoryReservation),
		LimitsMemorySwap:                optionalString(response.LimitsMemorySwap),
		LimitsMemorySwappiness:          optionalInt64(response.LimitsMemorySwappiness),
		ManualWebhookSecretBitbucket:    optionalString(response.ManualWebhookSecretBitbucket),
		ManualWebhookSecretGitea:        optionalString(response.ManualWebhookSecretGitea),
		ManualWebhookSecretGithub:       optionalString(response.ManualWebhookSecretGithub),
		ManualWebhookSecretGitlab:       optionalString(response.ManualWebhookSecretGitlab),
		Name:                            optionalString(response.Name),
		PortsExposes:                    optionalString(response.PortsExposes),
		PortsMappings:                   optionalString(response.PortsMappings),
		PostDeploymentCommand:           optionalString(response.PostDeploymentCommand),
		PostDeploymentCommandContainer:  optionalString(response.PostDeploymentCommandContainer),
		PreDeploymentCommand:            optionalString(response.PreDeploymentCommand),
		PreDeploymentCommandContainer:   optionalString(response.PreDeploymentCommandContainer),
		PreviewUrlTemplate:              optionalString(response.PreviewUrlTemplate),
		PrivateKeyId:                    optionalInt64(response.PrivateKeyId),
		PublishDirectory:                optionalString(response.PublishDirectory),
		Redirect:                        optionalString((*string)(response.Redirect)), // enum value
		RepositoryProjectId:             optionalInt64(response.RepositoryProjectId),
		SourceId:                        optionalInt64(response.SourceId),
		StartCommand:                    optionalString(response.StartCommand),
		StaticImage:                     optionalString(response.StaticImage),
		Status:                          optionalString(response.Status),
		SwarmPlacementConstraints:       optionalString(response.SwarmPlacementConstraints),
		SwarmReplicas:                   optionalInt64(response.SwarmReplicas),
		UpdatedAt:                       optionalTime(response.UpdatedAt),
		Uuid:                            optionalString(response.Uuid),
		WatchPaths:                      optionalString(response.WatchPaths),
	}
}
