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
	"terraform-provider-coolify/internal/provider/generated/datasource_projects"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &projectsDataSource{}
var _ datasource.DataSourceWithConfigure = &projectsDataSource{}

func NewProjectsDataSource() datasource.DataSource {
	return &projectsDataSource{}
}

type projectsDataSource struct {
	client *api.ClientWithResponses
}

type projectsDataSourceWithFilterModel struct {
	datasource_projects.ProjectsModel
	Filter []filter.BlockModel `tfsdk:"filter"`
}

var projectsFilterNames = []string{"id", "uuid", "name", "description"}

func (d *projectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *projectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_projects.ProjectsDataSourceSchema(ctx)
	resp.Schema.Description = "Get a list of Coolify projects."

	resp.Schema.Blocks = map[string]schema.Block{
		"filter": filter.CreateDatasourceFilter(projectsFilterNames),
	}
}

func (d *projectsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.client, resp)
}

func (d *projectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan projectsDataSourceWithFilterModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	listResponse, err := d.client.ListProjectsWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading projects", err.Error(),
		)
		return
	}

	if listResponse.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading projects",
			fmt.Sprintf("Received %s for projects. Details: %s", listResponse.Status(), listResponse.Body),
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

func (d *projectsDataSource) apiToModel(
	ctx context.Context,
	response *[]api.Project,
	filters []filter.BlockModel,
) (projectsDataSourceWithFilterModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var projects []attr.Value

	for _, project := range *response {
		attributes := map[string]attr.Value{
			"description": flatten.String(project.Description),
			"id":          flatten.Int64(project.Id),
			"name":        flatten.String(project.Name),
			"uuid":        flatten.String(project.Uuid),
		}

		if !filter.OnAttributes(attributes, filters) {
			continue
		}

		data, diag := datasource_projects.NewProjectsValue(
			datasource_projects.ProjectsValue{}.AttributeTypes(ctx),
			attributes)
		diags.Append(diag...)
		projects = append(projects, data)
	}

	dataSet, diag := types.SetValue(datasource_projects.ProjectsValue{}.Type(ctx), projects)
	diags.Append(diag...)

	return projectsDataSourceWithFilterModel{
		ProjectsModel: datasource_projects.ProjectsModel{
			Projects: dataSet,
		},
		Filter: filters,
	}, diags
}
