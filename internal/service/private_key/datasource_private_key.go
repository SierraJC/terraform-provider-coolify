package private_key

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &privateKeyDataSource{}
var _ datasource.DataSourceWithConfigure = &privateKeyDataSource{}

func NewPrivateKeyDataSource() datasource.DataSource {
	return &privateKeyDataSource{}
}

type privateKeyDataSource struct {
	client *api.ClientWithResponses
}

func (d *privateKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_key"
}

func (d *privateKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get a single Coolify private key by UUID.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Computed: true,
			},
			"fingerprint": schema.StringAttribute{
				Computed:            true,
				Description:         "The fingerprint of the private key.",
				MarkdownDescription: "The fingerprint of the private key.",
			},
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"is_git_related": schema.BoolAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"private_key": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"public_key": schema.StringAttribute{
				Computed:            true,
				Description:         "The public key of the private key.",
				MarkdownDescription: "The public key of the private key.",
			},
			"team_id": schema.Int64Attribute{
				Computed: true,
			},
			"uuid": schema.StringAttribute{
				Required:            true,
				Description:         "Private Key UUID",
				MarkdownDescription: "Private Key UUID",
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *privateKeyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.client, resp)
}

func (d *privateKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan privateKeyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	privateKey, err := d.client.GetPrivateKeyByUuidWithResponse(ctx, plan.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading private key", err.Error(),
		)
		return
	}

	if privateKey.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading private key",
			fmt.Sprintf("Received %s for private key. Details: %s", privateKey.Status(), privateKey.Body),
		)
		return
	}

	state := privateKeyModel{}.FromAPI(privateKey.JSON200)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
