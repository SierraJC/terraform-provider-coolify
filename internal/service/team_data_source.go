package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/util"
)

var _ datasource.DataSource = &teamDataSource{}
var _ datasource.DataSourceWithConfigure = &teamDataSource{}

type teamDataSourceModel = teamModel

func NewTeamDataSource() datasource.DataSource {
	return &teamDataSource{}
}

type teamDataSource struct {
	client *api.ClientWithResponses
}

func (d *teamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (d *teamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	notificationsDeprecationMessage := "No longer available via Coolify API. Will be removed in next update."

	resp.Schema = schema.Schema{
		Description: "Get a Coolify team by optional `id`. If no `id` is provided, the team associated with the current API key will be returned.",
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time the team was created.",
			},
			"custom_server_limit": schema.StringAttribute{
				Computed:    true,
				Description: "The custom server limit.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description of the team.",
			},
			"discord_enabled": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether Discord is enabled or not.",
			},
			"discord_notifications_database_backups": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send database backup notifications via Discord.",
			},
			"discord_notifications_deployments": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send deployment notifications via Discord.",
			},
			"discord_notifications_scheduled_tasks": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send scheduled task notifications via Discord.",
			},
			"discord_notifications_server_disk_usage": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send server disk usage notifications via Discord.",
			},
			"discord_notifications_status_changes": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send status change notifications via Discord.",
			},
			"discord_notifications_test": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send test notifications via Discord.",
			},
			"discord_webhook_url": schema.StringAttribute{
				Computed:           true,
				Sensitive:          true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The Discord webhook URL.",
			},
			"id": schema.Int64Attribute{
				Required:    false,
				Optional:    true,
				Description: "The unique identifier of the team.",
			},
			"members": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "The date when the user was created.",
						},
						"email": schema.StringAttribute{
							Computed:    true,
							Sensitive:   false, // todo: mark sensitive?
							Description: "The user email.",
						},
						"email_verified_at": schema.StringAttribute{
							Computed:    true,
							Description: "The date when the user email was verified.",
						},
						"force_password_reset": schema.BoolAttribute{
							Computed:    true,
							Description: "The flag to force the user to reset the password.",
						},
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "The user identifier in the database.",
						},
						"marketing_emails": schema.BoolAttribute{
							Computed:    true,
							Description: "The flag to receive marketing emails.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The user name.",
						},
						"two_factor_confirmed_at": schema.StringAttribute{
							Computed:    true,
							Description: "The date when the user two factor was confirmed.",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "The date when the user was updated.",
						},
					},
				},
				Computed:    true,
				Description: "The members of the team.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the team.",
			},
			"personal_team": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the team is personal or not.",
			},
			"resend_api_key": schema.StringAttribute{
				Computed:           true,
				Sensitive:          true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The resending API key.",
			},
			"resend_enabled": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to enable resending or not.",
			},
			"show_boarding": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to show the boarding screen or not.",
			},
			"smtp_enabled": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether SMTP is enabled or not.",
			},
			"smtp_encryption": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The SMTP encryption.",
			},
			"smtp_from_address": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The email address to send emails from.",
			},
			"smtp_from_name": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The name to send emails from.",
			},
			"smtp_host": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The SMTP host.",
			},
			"smtp_notifications_database_backups": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send database backup notifications via SMTP.",
			},
			"smtp_notifications_deployments": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send deployment notifications via SMTP.",
			},
			"smtp_notifications_scheduled_tasks": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send scheduled task notifications via SMTP.",
			},
			"smtp_notifications_server_disk_usage": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send server disk usage notifications via SMTP.",
			},
			"smtp_notifications_status_changes": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send status change notifications via SMTP.",
			},
			"smtp_notifications_test": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send test notifications via SMTP.",
			},
			"smtp_password": schema.StringAttribute{
				Computed:           true,
				Sensitive:          true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The SMTP password.",
			},
			"smtp_port": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The SMTP port.",
			},
			"smtp_recipients": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The email addresses to send emails to.",
			},
			"smtp_timeout": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The SMTP timeout.",
			},
			"smtp_username": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The SMTP username.",
			},
			"telegram_chat_id": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The Telegram chat ID.",
			},
			"telegram_enabled": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether Telegram is enabled or not.",
			},
			"telegram_notifications_database_backups": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send database backup notifications via Telegram.",
			},
			"telegram_notifications_database_backups_message_thread_id": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The Telegram database backup message thread ID.",
			},
			"telegram_notifications_deployments": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send deployment notifications via Telegram.",
			},
			"telegram_notifications_deployments_message_thread_id": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The Telegram deployment message thread ID.",
			},
			"telegram_notifications_scheduled_tasks": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send scheduled task notifications via Telegram.",
			},
			"telegram_notifications_scheduled_tasks_thread_id": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The Telegram scheduled task message thread ID.",
			},
			"telegram_notifications_status_changes": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send status change notifications via Telegram.",
			},
			"telegram_notifications_status_changes_message_thread_id": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The Telegram status change message thread ID.",
			},
			"telegram_notifications_test": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to send test notifications via Telegram.",
			},
			"telegram_notifications_test_message_thread_id": schema.StringAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The Telegram test message thread ID.",
			},
			"telegram_token": schema.StringAttribute{
				Computed:           true,
				Sensitive:          true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "The Telegram token.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time the team was last updated.",
			},
			"use_instance_email_settings": schema.BoolAttribute{
				Computed:           true,
				DeprecationMessage: notificationsDeprecationMessage,
				Description:        "Whether to use instance email settings or not.",
			},
		},
	}
}

func (d *teamDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	util.ProviderDataFromDataSourceConfigureRequest(req, &d.client, resp)
}

func (d *teamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan teamDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var team *api.Team
	var teamMembers *[]api.User

	if !plan.Id.IsNull() {
		// Get team by ID
		teamResp, err := d.client.GetTeamByIdWithResponse(ctx, int(plan.Id.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading team", err.Error(),
			)
			return
		}

		if teamResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError(
				"Unexpected HTTP status code reading team",
				fmt.Sprintf("Received %s for team. Details: %s", teamResp.Status(), string(teamResp.Body)),
			)
			return
		}

		team = teamResp.JSON200
		teamMembers = team.Members
	} else {
		// Get current team
		teamResp, err := d.client.GetCurrentTeamWithResponse(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading team", err.Error(),
			)
			return
		}

		if teamResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError(
				"Unexpected HTTP status code reading team",
				fmt.Sprintf("Received %s for team. Details: %s", teamResp.Status(), string(teamResp.Body)),
			)
			return
		}

		team = teamResp.JSON200
		teamMembers = team.Members
	}

	// If the API did not provide members, we need to fetch them separately
	// TODO: Coolify API inconsistency: Spec says members should be returned, but it is not.
	if teamMembers == nil {
		teamMembersResponse, err := d.client.GetMembersByTeamIdWithResponse(ctx, *team.Id)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading team members", err.Error(),
			)
			return
		}

		if teamMembersResponse.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError(
				"Unexpected HTTP status code reading team members",
				fmt.Sprintf("Received %s for team members. Details: %s", teamMembersResponse.Status(), string(teamMembersResponse.Body)),
			)
			return
		}
		team.Members = teamMembersResponse.JSON200
	}

	state := teamDataSourceModel{}.FromAPI(team)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
