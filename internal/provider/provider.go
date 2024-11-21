package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
)

const (
	ENV_KEY_ENDPOINT = "COOLIFY_ENDPOINT"
	ENV_KEY_TOKEN    = "COOLIFY_TOKEN"

	DEFAULT_COOLIFY_ENDPOINT = "https://app.coolify.io/api/v1"
	MIN_COOLIFY_VERSION      = "4.0.0-beta.364"
)

// Ensure the implementation satisfies the expected interfaces.
var _ provider.Provider = &CoolifyProvider{}
var _ provider.ProviderWithFunctions = &CoolifyProvider{}
var _ provider.ProviderWithEphemeralResources = &CoolifyProvider{}

// CoolifyProvider defines the provider implementation.
type CoolifyProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type CoolifyProviderData struct {
	endpoint string
	client   *api.ClientWithResponses
}

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

func (p *CoolifyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "coolify"
	resp.Version = p.version
}

func (p *CoolifyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	hasEnvToken := os.Getenv(ENV_KEY_TOKEN) != ""
	resp.Schema = schema.Schema{
		MarkdownDescription: "" +
			"The \"coolify\" provider facilitates interaction with resources supported by [Coolify](https://coolify.io/) v" + MIN_COOLIFY_VERSION + " and later.\n\n" +
			"Before using this provider, you must configure it with your credentials, typically by setting the environment variable `" + ENV_KEY_TOKEN + "`.\n\n" +
			"For instructions on obtaining an API token, refer to Coolify's [API documentation](https://coolify.io/docs/api-reference/authorization#generate).",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "Coolify endpoint. If not set, checks env for `" + ENV_KEY_ENDPOINT + "`. Default: `" + DEFAULT_COOLIFY_ENDPOINT + "`.",
			},
			"token": schema.StringAttribute{
				Required:  !hasEnvToken,
				Optional:  hasEnvToken,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(api.TokenRegex, api.ErrInvalidToken.Error()),
				},
				Description: "Coolify token. If not set, checks env for `" + ENV_KEY_TOKEN + "`.",
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
		apiEndpoint = DEFAULT_COOLIFY_ENDPOINT
	}

	if apiEndpoint == "" {
		resp.Diagnostics.AddAttributeError(path.Root("endpoint"), "Failed to configure client", "No API Endpoint provided")
	}

	var apiToken string
	if !data.Token.IsNull() {
		apiToken = data.Token.ValueString()
	} else {
		if apiTokenFromEnv, found := os.LookupEnv(ENV_KEY_TOKEN); found {
			apiToken = apiTokenFromEnv
		}
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(path.Root("token"), "Failed to configure client", "No token provided")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := api.NewAPIClient(p.version, apiEndpoint, apiToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create API client",
			err.Error(),
		)
		return
	}

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

	if !isVersionCompatible(currentVersion, MIN_COOLIFY_VERSION) {
		resp.Diagnostics.AddError(
			"Unsupported API version",
			fmt.Sprintf("The Coolify API version %s is not supported. The minimum supported version is %s", currentVersion, MIN_COOLIFY_VERSION),
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

func (p *CoolifyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPrivateKeyResource,
		NewServerResource,
		NewProjectResource,
		NewApplicationEnvsResource,
		NewServiceEnvsResource,
	}
}

func (p *CoolifyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPrivateKeyDataSource,
		NewPrivateKeysDataSource,
		NewTeamDataSource,
		NewTeamsDataSource,
		NewServerDataSource,
		NewServersDataSource,
		NewServerResourcesDataSource,
		NewServerDomainsDataSource,
		NewProjectDataSource,
		NewProjectsDataSource,
		NewApplicationDataSource,
		NewServiceDataSource,
	}
}

func (p *CoolifyProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *CoolifyProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}
