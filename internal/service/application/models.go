package application

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/flatten"
)

type commonApplicationModel struct {
	Description     types.String `tfsdk:"description"`
	DestinationUuid types.String `tfsdk:"destination_uuid"`
	EnvironmentName types.String `tfsdk:"environment_name"`
	EnvironmentUuid types.String `tfsdk:"environment_uuid"`
	InstantDeploy   types.Bool   `tfsdk:"instant_deploy"`
	Name            types.String `tfsdk:"name"`
	ProjectUuid     types.String `tfsdk:"project_uuid"`
	ServerUuid      types.String `tfsdk:"server_uuid"`
	Uuid            types.String `tfsdk:"uuid"`
}

func (m commonApplicationModel) CommonSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Description of the application",
			},
			"destination_uuid": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "UUID of the destination if the server has multiple destinations",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Default:       stringdefault.StaticString(""),
			},
			"environment_name": schema.StringAttribute{
				Required:      true,
				Description:   "Name of the environment",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"environment_uuid": schema.StringAttribute{
				Optional:      true, // todo: should change this to required and optional environment name
				Description:   "UUID of the environment. Will replace environment_name in future.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"instant_deploy": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Instant deploy the applciation.",
				Default:     booldefault.StaticBool(false),
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the application",
			},
			"project_uuid": schema.StringAttribute{
				Required:      true,
				Description:   "UUID of the project",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"server_uuid": schema.StringAttribute{
				Required:      true,
				Description:   "UUID of the server",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"uuid": schema.StringAttribute{
				Computed:      true,
				Description:   "UUID of the application",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (m commonApplicationModel) FromAPI(app *api.Application, state commonApplicationModel) commonApplicationModel {
	return commonApplicationModel{
		Uuid:            flatten.String(app.Uuid),
		Name:            flatten.String(app.Name),
		Description:     flatten.String(app.Description),
		ServerUuid:      state.ServerUuid, // Values not returned by API, so use the plan value
		ProjectUuid:     state.ProjectUuid,
		EnvironmentName: state.EnvironmentName,
		EnvironmentUuid: state.EnvironmentUuid,
		DestinationUuid: state.DestinationUuid,
		InstantDeploy:   state.InstantDeploy,
	}
}

type dockerComposeApplicationModel struct {
	commonApplicationModel
	DockerComposeRaw types.String `tfsdk:"dockercompose_raw"`
}

func (m dockerComposeApplicationModel) FromAPI(app *api.Application, state dockerComposeApplicationModel) dockerComposeApplicationModel {
	return dockerComposeApplicationModel{
		commonApplicationModel: commonApplicationModel{}.FromAPI(app, state.commonApplicationModel),
		DockerComposeRaw:       flatten.String(app.DockerComposeRaw),
	}
}
