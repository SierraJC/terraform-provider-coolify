package service

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/flatten"
	sutil "terraform-provider-coolify/internal/service/util"
)

type ServiceModel struct {
	Uuid                   types.String   `tfsdk:"uuid"`
	Name                   types.String   `tfsdk:"name"`
	Description            types.String   `tfsdk:"description"`
	DestinationUuid        types.String   `tfsdk:"destination_uuid"`
	EnvironmentName        types.String   `tfsdk:"environment_name"`
	EnvironmentUuid        types.String   `tfsdk:"environment_uuid"`
	ProjectUuid            types.String   `tfsdk:"project_uuid"`
	ServerUuid             types.String   `tfsdk:"server_uuid"`
	InstantDeploy          types.Bool     `tfsdk:"instant_deploy"`
	Compose                types.String   `tfsdk:"compose"`
	Status                 types.String   `tfsdk:"status"`
	ServerStatus           types.String   `tfsdk:"server_status"`
	WaitForDeployment      types.Bool     `tfsdk:"wait_for_deployment"`
	Timeouts               timeouts.Value `tfsdk:"timeouts"`
}

func (m ServiceModel) Schema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Create, read, update, and delete a Coolify service resource.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:      true,
				Description:   "UUID of the service.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the service.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Description of the service.",
			},
			"destination_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "UUID of the destination. Optional - if omitted, Coolify will auto-select the first available destination on the server. Required only when the server has multiple destinations.",
			},
			"environment_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the environment.",
			},
			"environment_uuid": schema.StringAttribute{
				Optional:    true, // todo: should change this to required and optional environment name
				Description: "UUID of the environment. Will replace environment_name in future.",
			},
			"instant_deploy": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Instant deploy the service.",
				Default:     booldefault.StaticBool(false),
			},
			"project_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the project.",
			},
			"server_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the server.",
			},
			"compose": schema.StringAttribute{
				Required:            true,
				Description:         "The Docker Compose raw content.",
				MarkdownDescription: "The Docker Compose raw content.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The aggregate status of the service from its applications and databases. Format: status:health or status:health:excluded (e.g., running:healthy, degraded:unhealthy).",
			},
			"server_status": schema.StringAttribute{
				Computed:    true,
				Description: "The functional status of the server where the service is running.",
			},
			"wait_for_deployment": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Wait for service deployment (status becomes running:healthy) during creation and updates. Defaults to false.",
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Update: true,
			}),
		},
	}
}

func (m ServiceModel) FromAPI(service *api.Service, state ServiceModel) ServiceModel {
	return ServiceModel{
		Uuid:            flatten.String(service.Uuid),
		Name:            flatten.String(service.Name),
		Description:     flatten.String(service.Description),
		Status:          flatten.String(service.Status),
		ServerStatus:    flatten.String(service.ServerStatus),
		ServerUuid:      state.ServerUuid, // Values not returned by API, so use the plan value
		ProjectUuid:     state.ProjectUuid,
		EnvironmentName: state.EnvironmentName,
		EnvironmentUuid: state.EnvironmentUuid,
		DestinationUuid: state.DestinationUuid,
		InstantDeploy:   state.InstantDeploy,
		Compose:         state.Compose,
		WaitForDeployment: state.WaitForDeployment,
		Timeouts:        state.Timeouts,
	}
}

func (m ServiceModel) ToAPICreate() api.CreateServiceJSONRequestBody {
	body := api.CreateServiceJSONRequestBody{
		Name:             m.Name.ValueStringPointer(),
		Description:      m.Description.ValueStringPointer(),
		EnvironmentName:  m.EnvironmentName.ValueString(),
		EnvironmentUuid:  m.EnvironmentUuid.ValueString(),
		ProjectUuid:      m.ProjectUuid.ValueString(),
		ServerUuid:       m.ServerUuid.ValueString(),
		InstantDeploy:    m.InstantDeploy.ValueBoolPointer(),
		DockerComposeRaw: sutil.Base64EncodeAttr(m.Compose),
	}
	// Only send destination_uuid if it's actually set
	if !m.DestinationUuid.IsNull() && !m.DestinationUuid.IsUnknown() && m.DestinationUuid.ValueString() != "" {
		body.DestinationUuid = m.DestinationUuid.ValueStringPointer()
	}
	return body
}
func (m ServiceModel) ToAPIUpdate() api.UpdateServiceByUuidJSONRequestBody {
	body := api.UpdateServiceByUuidJSONRequestBody{
		Name:             m.Name.ValueStringPointer(),
		Description:      m.Description.ValueStringPointer(),
		EnvironmentName:  m.EnvironmentName.ValueString(),
		EnvironmentUuid:  m.EnvironmentUuid.ValueString(),
		ProjectUuid:      m.ProjectUuid.ValueString(),
		ServerUuid:       m.ServerUuid.ValueString(),
		InstantDeploy:    m.InstantDeploy.ValueBoolPointer(),
		DockerComposeRaw: *sutil.Base64EncodeAttr(m.Compose),
	}
	// Only send destination_uuid if it's actually set
	if !m.DestinationUuid.IsNull() && !m.DestinationUuid.IsUnknown() && m.DestinationUuid.ValueString() != "" {
		body.DestinationUuid = m.DestinationUuid.ValueStringPointer()
	}
	return body
}
