// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package resource_mysql_database

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func MysqlDatabaseResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Description of the database",
				MarkdownDescription: "Description of the database",
			},
			"destination_uuid": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "UUID of the destination if the server has multiple destinations",
				MarkdownDescription: "UUID of the destination if the server has multiple destinations",
			},
			"environment_name": schema.StringAttribute{
				Required:            true,
				Description:         "Name of the environment. You need to provide at least one of environment_name or environment_uuid.",
				MarkdownDescription: "Name of the environment. You need to provide at least one of environment_name or environment_uuid.",
			},
			"environment_uuid": schema.StringAttribute{
				Required:            true,
				Description:         "UUID of the environment. You need to provide at least one of environment_name or environment_uuid.",
				MarkdownDescription: "UUID of the environment. You need to provide at least one of environment_name or environment_uuid.",
			},
			"image": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Docker Image of the database",
				MarkdownDescription: "Docker Image of the database",
			},
			"instant_deploy": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Instant deploy the database",
				MarkdownDescription: "Instant deploy the database",
			},
			"is_public": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Is the database public?",
				MarkdownDescription: "Is the database public?",
			},
			"limits_cpu_shares": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "CPU shares of the database",
				MarkdownDescription: "CPU shares of the database",
			},
			"limits_cpus": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "CPU limit of the database",
				MarkdownDescription: "CPU limit of the database",
			},
			"limits_cpuset": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "CPU set of the database",
				MarkdownDescription: "CPU set of the database",
			},
			"limits_memory": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Memory limit of the database",
				MarkdownDescription: "Memory limit of the database",
			},
			"limits_memory_reservation": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Memory reservation of the database",
				MarkdownDescription: "Memory reservation of the database",
			},
			"limits_memory_swap": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Memory swap limit of the database",
				MarkdownDescription: "Memory swap limit of the database",
			},
			"limits_memory_swappiness": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Memory swappiness of the database",
				MarkdownDescription: "Memory swappiness of the database",
			},
			"mysql_conf": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "MySQL conf",
				MarkdownDescription: "MySQL conf",
			},
			"mysql_database": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "MySQL database",
				MarkdownDescription: "MySQL database",
			},
			"mysql_password": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "MySQL password",
				MarkdownDescription: "MySQL password",
			},
			"mysql_root_password": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "MySQL root password",
				MarkdownDescription: "MySQL root password",
			},
			"mysql_user": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "MySQL user",
				MarkdownDescription: "MySQL user",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Name of the database",
				MarkdownDescription: "Name of the database",
			},
			"project_uuid": schema.StringAttribute{
				Required:            true,
				Description:         "UUID of the project",
				MarkdownDescription: "UUID of the project",
			},
			"public_port": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Public port of the database",
				MarkdownDescription: "Public port of the database",
			},
			"server_uuid": schema.StringAttribute{
				Required:            true,
				Description:         "UUID of the server",
				MarkdownDescription: "UUID of the server",
			},
			"uuid": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "UUID of the database.",
				MarkdownDescription: "UUID of the database.",
			},
		},
	}
}

type MysqlDatabaseModel struct {
	Description             types.String `tfsdk:"description"`
	DestinationUuid         types.String `tfsdk:"destination_uuid"`
	EnvironmentName         types.String `tfsdk:"environment_name"`
	EnvironmentUuid         types.String `tfsdk:"environment_uuid"`
	Image                   types.String `tfsdk:"image"`
	InstantDeploy           types.Bool   `tfsdk:"instant_deploy"`
	IsPublic                types.Bool   `tfsdk:"is_public"`
	LimitsCpuShares         types.Int64  `tfsdk:"limits_cpu_shares"`
	LimitsCpus              types.String `tfsdk:"limits_cpus"`
	LimitsCpuset            types.String `tfsdk:"limits_cpuset"`
	LimitsMemory            types.String `tfsdk:"limits_memory"`
	LimitsMemoryReservation types.String `tfsdk:"limits_memory_reservation"`
	LimitsMemorySwap        types.String `tfsdk:"limits_memory_swap"`
	LimitsMemorySwappiness  types.Int64  `tfsdk:"limits_memory_swappiness"`
	MysqlConf               types.String `tfsdk:"mysql_conf"`
	MysqlDatabase           types.String `tfsdk:"mysql_database"`
	MysqlPassword           types.String `tfsdk:"mysql_password"`
	MysqlRootPassword       types.String `tfsdk:"mysql_root_password"`
	MysqlUser               types.String `tfsdk:"mysql_user"`
	Name                    types.String `tfsdk:"name"`
	ProjectUuid             types.String `tfsdk:"project_uuid"`
	PublicPort              types.Int64  `tfsdk:"public_port"`
	ServerUuid              types.String `tfsdk:"server_uuid"`
	Uuid                    types.String `tfsdk:"uuid"`
}
