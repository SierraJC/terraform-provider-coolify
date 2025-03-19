package service

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/flatten"
	sutil "terraform-provider-coolify/internal/service/util"
)

type ServiceModel struct {
	Uuid            types.String `tfsdk:"uuid"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	DestinationUuid types.String `tfsdk:"destination_uuid"`
	EnvironmentName types.String `tfsdk:"environment_name"`
	EnvironmentUuid types.String `tfsdk:"environment_uuid"`
	ProjectUuid     types.String `tfsdk:"project_uuid"`
	ServerUuid      types.String `tfsdk:"server_uuid"`
	InstantDeploy   types.Bool   `tfsdk:"instant_deploy"`
	Compose         types.String `tfsdk:"compose"`
}

func (m ServiceModel) Schema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Create, read, update, and delete a Coolify service resource.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Name of the service.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"uuid": schema.StringAttribute{
				Computed:      true,
				Description:   "UUID of the service.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"description": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Description of the service.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"destination_uuid": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "UUID of the destination if the server has multiple destinations.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Default:       stringdefault.StaticString(""),
			},
			"environment_name": schema.StringAttribute{
				Required:      true,
				Description:   "Name of the environment.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"environment_uuid": schema.StringAttribute{
				Optional:      true, // todo: should change this to required and optional environment name
				Description:   "UUID of the environment. Will replace environment_name in future.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"instant_deploy": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Instant deploy the service.",
				Default:       booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"project_uuid": schema.StringAttribute{
				Required:      true,
				Description:   "UUID of the project.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"server_uuid": schema.StringAttribute{
				Required:      true,
				Description:   "UUID of the server.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"compose": schema.StringAttribute{
				Required:            true,
				Description:         "The Docker Compose raw content.",
				MarkdownDescription: "The Docker Compose raw content.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (m ServiceModel) FromAPI(service *api.Service, state ServiceModel) ServiceModel {
	return ServiceModel{
		Uuid:            flatten.String(service.Uuid),
		Name:            flatten.String(service.Name),
		Description:     flatten.String(service.Description),
		ServerUuid:      state.ServerUuid, // Values not returned by API, so use the plan value
		ProjectUuid:     state.ProjectUuid,
		EnvironmentName: state.EnvironmentName,
		EnvironmentUuid: state.EnvironmentUuid,
		DestinationUuid: state.DestinationUuid,
		InstantDeploy:   state.InstantDeploy,
		Compose:         state.Compose,
	}
}

func (m ServiceModel) ToAPICreate() api.CreateServiceJSONRequestBody {
	c := false
	return api.CreateServiceJSONRequestBody{
		Name:                   m.Name.ValueStringPointer(),
		Description:            m.Description.ValueStringPointer(),
		DestinationUuid:        m.DestinationUuid.ValueStringPointer(),
		EnvironmentName:        m.EnvironmentName.ValueString(),
		EnvironmentUuid:        m.EnvironmentUuid.ValueString(),
		InstantDeploy:          m.InstantDeploy.ValueBoolPointer(),
		ProjectUuid:            m.ProjectUuid.ValueString(),
		ServerUuid:             m.ServerUuid.ValueString(),
		DockerComposeRaw:       *sutil.Base64EncodeAttr(m.Compose),
		ConnectToDockerNetwork: &c,
	}
}
