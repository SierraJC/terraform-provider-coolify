package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
)

var _ provider.Provider = (*CoolifyProvider)(nil)

type CoolifyProvider struct {
	version string
}

// type coolifyClient struct {
// 	endpoint string
// 	client   *api.APIClient
// }

type CoolifyProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CoolifyProvider{
			version: version,
		}
	}
}

func (p *CoolifyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Required:            os.Getenv("COOLIFY_ENDPOINT") == "",
				Description:         "The endpoint for the Coolify API",
				MarkdownDescription: "The endpoint for the Coolify API",
			},
			"token": schema.StringAttribute{
				Required:            os.Getenv("COOLIFY_TOKEN") == "",
				Sensitive:           true,
				Description:         "The API key for authenticating with Coolify",
				MarkdownDescription: "The API key for authenticating with Coolify",
			},
		},
	}
}

func (p *CoolifyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CoolifyProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiEndpoint string
	if !data.Endpoint.IsNull() {
		apiEndpoint = data.Endpoint.ValueString()
	} else if apiEndpointFromEnv, found := os.LookupEnv("COOLIFY_ENDPOINT"); found {
		apiEndpoint = apiEndpointFromEnv
	} else {
		apiEndpoint = api.DefaultServerURL
	}

	if apiEndpoint == "" {
		resp.Diagnostics.AddAttributeError(path.Root("endpoint"), "Failed to configure client", "No API Endpoint provided")
	}

	var apiToken string
	if !data.Token.IsNull() {
		apiToken = data.Token.ValueString()
	} else {
		if apiTokenFromEnv, found := os.LookupEnv("COOLIFY_TOKEN"); found {
			apiToken = apiTokenFromEnv
		}
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(path.Root("token"), "Failed to configure client", "No token provided")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client := api.NewAPIClient(apiEndpoint, apiToken)

	// GET /version
	versionResp, err := client.VersionWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to connect to Coolify API",
			err.Error(),
		)
		return
	}

	if versionResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code API client",
			fmt.Sprintf("Received %s creating API client. Details: %s", versionResp.Status(), versionResp.Body),
		)
		return
	}

	currentVersion := string(versionResp.Body)
	minVersion := "4.0.0-beta.318"

	if !isVersionCompatible(currentVersion, minVersion) {
		resp.Diagnostics.AddError(
			"Unsupported API version",
			fmt.Sprintf("The Coolify API version %s is not supported. The minimum supported version is %s", currentVersion, minVersion),
		)
		return
	}

	tflog.Info(ctx, "Successfully connected to Coolify API", map[string]interface{}{"version": currentVersion})

	providerData := &CoolifyProviderData{
		endpoint: apiEndpoint,
		client:   client,
	}

	resp.ResourceData = providerData
	resp.DataSourceData = providerData
}

func (p *CoolifyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "coolify"
	resp.Version = p.version
}

func (p *CoolifyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPrivateKeysDataSource,
	}
}

func (p *CoolifyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewServerResource,
		NewProjectResource,
		NewPrivateKeyResource,
	}
}
