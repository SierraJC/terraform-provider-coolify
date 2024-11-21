package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/generated/datasource_service"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &serviceDataSource{}
var _ datasource.DataSourceWithConfigure = &serviceDataSource{}

func NewServiceDataSource() datasource.DataSource {
	return &serviceDataSource{}
}

type serviceDataSource struct {
	providerData CoolifyProviderData
}

func (d *serviceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (d *serviceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_service.ServiceDataSourceSchema(ctx)
	resp.Schema.Description = "Get a Coolify service by `uuid`."

	// Mark sensitive attributes
	sensitiveAttrs := []string{"manual_webhook_secret_bitbucket", "manual_webhook_secret_gitea", "manual_webhook_secret_github", "manual_webhook_secret_gitlab"}
	for _, attr := range sensitiveAttrs {
		makeDataSourceAttributeSensitive(resp.Schema.Attributes, attr)
	}
}

func (d *serviceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.providerData, resp)
}

func (d *serviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan datasource_service.ServiceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	serviceResp, err := d.providerData.client.GetServiceByUuidWithResponse(ctx, plan.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading service", err.Error(),
		)
		return
	}

	if serviceResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading service",
			fmt.Sprintf("Received %s for service. Details: %s", serviceResp.Status(), string(serviceResp.Body)),
		)
		return
	}

	state := d.ApiToModel(ctx, &resp.Diagnostics, serviceResp.JSON200)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *serviceDataSource) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *api.Service,
) datasource_service.ServiceModel {
	return datasource_service.ServiceModel{
		ConfigHash:                      optionalString(response.ConfigHash),
		ConnectToDockerNetwork:          optionalBool(response.ConnectToDockerNetwork),
		CreatedAt:                       optionalString(response.CreatedAt),
		DeletedAt:                       optionalString(response.DeletedAt),
		Description:                     optionalString(response.Description),
		DestinationId:                   optionalInt64(response.DestinationId),
		DestinationType:                 optionalString(response.DestinationType),
		DockerCompose:                   optionalString(response.DockerCompose),
		DockerComposeRaw:                optionalString(response.DockerComposeRaw),
		EnvironmentId:                   optionalInt64(response.EnvironmentId),
		Id:                              optionalInt64(response.Id),
		IsContainerLabelEscapeEnabled:   optionalBool(response.IsContainerLabelEscapeEnabled),
		IsContainerLabelReadonlyEnabled: optionalBool(response.IsContainerLabelReadonlyEnabled),
		Name:                            optionalString(response.Name),
		ServerId:                        optionalInt64(response.ServerId),
		ServiceType:                     optionalString((*string)(response.ServiceType)), // enum value
		UpdatedAt:                       optionalString(response.UpdatedAt),
		Uuid:                            optionalString(response.Uuid),
	}
}
