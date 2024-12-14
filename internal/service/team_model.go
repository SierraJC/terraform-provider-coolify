package service

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/filter"
	"terraform-provider-coolify/internal/flatten"
)

type teamModel struct {
	CreatedAt                                           types.String      `tfsdk:"created_at"`
	CustomServerLimit                                   types.String      `tfsdk:"custom_server_limit"`
	Description                                         types.String      `tfsdk:"description"`
	DiscordEnabled                                      types.Bool        `tfsdk:"discord_enabled"`
	DiscordNotificationsDatabaseBackups                 types.Bool        `tfsdk:"discord_notifications_database_backups"`
	DiscordNotificationsDeployments                     types.Bool        `tfsdk:"discord_notifications_deployments"`
	DiscordNotificationsScheduledTasks                  types.Bool        `tfsdk:"discord_notifications_scheduled_tasks"`
	DiscordNotificationsServerDiskUsage                 types.Bool        `tfsdk:"discord_notifications_server_disk_usage"`
	DiscordNotificationsStatusChanges                   types.Bool        `tfsdk:"discord_notifications_status_changes"`
	DiscordNotificationsTest                            types.Bool        `tfsdk:"discord_notifications_test"`
	DiscordWebhookUrl                                   types.String      `tfsdk:"discord_webhook_url"`
	Id                                                  types.Int64       `tfsdk:"id"`
	Members                                             []teamMemberModel `tfsdk:"members"`
	Name                                                types.String      `tfsdk:"name"`
	PersonalTeam                                        types.Bool        `tfsdk:"personal_team"`
	ResendApiKey                                        types.String      `tfsdk:"resend_api_key"`
	ResendEnabled                                       types.Bool        `tfsdk:"resend_enabled"`
	ShowBoarding                                        types.Bool        `tfsdk:"show_boarding"`
	SmtpEnabled                                         types.Bool        `tfsdk:"smtp_enabled"`
	SmtpEncryption                                      types.String      `tfsdk:"smtp_encryption"`
	SmtpFromAddress                                     types.String      `tfsdk:"smtp_from_address"`
	SmtpFromName                                        types.String      `tfsdk:"smtp_from_name"`
	SmtpHost                                            types.String      `tfsdk:"smtp_host"`
	SmtpNotificationsDatabaseBackups                    types.Bool        `tfsdk:"smtp_notifications_database_backups"`
	SmtpNotificationsDeployments                        types.Bool        `tfsdk:"smtp_notifications_deployments"`
	SmtpNotificationsScheduledTasks                     types.Bool        `tfsdk:"smtp_notifications_scheduled_tasks"`
	SmtpNotificationsServerDiskUsage                    types.Bool        `tfsdk:"smtp_notifications_server_disk_usage"`
	SmtpNotificationsStatusChanges                      types.Bool        `tfsdk:"smtp_notifications_status_changes"`
	SmtpNotificationsTest                               types.Bool        `tfsdk:"smtp_notifications_test"`
	SmtpPassword                                        types.String      `tfsdk:"smtp_password"`
	SmtpPort                                            types.String      `tfsdk:"smtp_port"`
	SmtpRecipients                                      types.String      `tfsdk:"smtp_recipients"`
	SmtpTimeout                                         types.String      `tfsdk:"smtp_timeout"`
	SmtpUsername                                        types.String      `tfsdk:"smtp_username"`
	TelegramChatId                                      types.String      `tfsdk:"telegram_chat_id"`
	TelegramEnabled                                     types.Bool        `tfsdk:"telegram_enabled"`
	TelegramNotificationsDatabaseBackups                types.Bool        `tfsdk:"telegram_notifications_database_backups"`
	TelegramNotificationsDatabaseBackupsMessageThreadId types.String      `tfsdk:"telegram_notifications_database_backups_message_thread_id"`
	TelegramNotificationsDeployments                    types.Bool        `tfsdk:"telegram_notifications_deployments"`
	TelegramNotificationsDeploymentsMessageThreadId     types.String      `tfsdk:"telegram_notifications_deployments_message_thread_id"`
	TelegramNotificationsScheduledTasks                 types.Bool        `tfsdk:"telegram_notifications_scheduled_tasks"`
	TelegramNotificationsScheduledTasksThreadId         types.String      `tfsdk:"telegram_notifications_scheduled_tasks_thread_id"`
	TelegramNotificationsStatusChanges                  types.Bool        `tfsdk:"telegram_notifications_status_changes"`
	TelegramNotificationsStatusChangesMessageThreadId   types.String      `tfsdk:"telegram_notifications_status_changes_message_thread_id"`
	TelegramNotificationsTest                           types.Bool        `tfsdk:"telegram_notifications_test"`
	TelegramNotificationsTestMessageThreadId            types.String      `tfsdk:"telegram_notifications_test_message_thread_id"`
	TelegramToken                                       types.String      `tfsdk:"telegram_token"`
	UpdatedAt                                           types.String      `tfsdk:"updated_at"`
	UseInstanceEmailSettings                            types.Bool        `tfsdk:"use_instance_email_settings"`
}

type teamMemberModel struct {
	CreatedAt            types.String `tfsdk:"created_at"`
	Email                types.String `tfsdk:"email"`
	EmailVerifiedAt      types.String `tfsdk:"email_verified_at"`
	ForcePasswordReset   types.Bool   `tfsdk:"force_password_reset"`
	Id                   types.Int64  `tfsdk:"id"`
	MarketingEmails      types.Bool   `tfsdk:"marketing_emails"`
	Name                 types.String `tfsdk:"name"`
	TwoFactorConfirmedAt types.String `tfsdk:"two_factor_confirmed_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

var _ filter.FilterableStructModel = teamModel{}

func (m teamModel) FromAPI(apiModel *api.Team) teamModel {
	var members []teamMemberModel

	if apiModel.Members != nil {
		members = make([]teamMemberModel, len(*apiModel.Members))
		for i, member := range *apiModel.Members {
			members[i] = teamMemberModel{
				CreatedAt:            flatten.String(member.CreatedAt),
				Email:                flatten.String(member.Email),
				EmailVerifiedAt:      flatten.String(member.EmailVerifiedAt),
				ForcePasswordReset:   flatten.Bool(member.ForcePasswordReset),
				Id:                   flatten.Int64(member.Id),
				MarketingEmails:      flatten.Bool(member.MarketingEmails),
				Name:                 flatten.String(member.Name),
				TwoFactorConfirmedAt: flatten.String(member.TwoFactorConfirmedAt),
				UpdatedAt:            flatten.String(member.UpdatedAt),
			}
		}
	}

	return teamModel{
		CreatedAt:                            flatten.String(apiModel.CreatedAt),
		CustomServerLimit:                    flatten.String(apiModel.CustomServerLimit),
		Description:                          flatten.String(apiModel.Description),
		DiscordEnabled:                       flatten.Bool(apiModel.DiscordEnabled),
		DiscordNotificationsDatabaseBackups:  flatten.Bool(apiModel.DiscordNotificationsDatabaseBackups),
		DiscordNotificationsDeployments:      flatten.Bool(apiModel.DiscordNotificationsDeployments),
		DiscordNotificationsScheduledTasks:   flatten.Bool(apiModel.DiscordNotificationsScheduledTasks),
		DiscordNotificationsServerDiskUsage:  flatten.Bool(apiModel.DiscordNotificationsServerDiskUsage),
		DiscordNotificationsStatusChanges:    flatten.Bool(apiModel.DiscordNotificationsStatusChanges),
		DiscordNotificationsTest:             flatten.Bool(apiModel.DiscordNotificationsTest),
		DiscordWebhookUrl:                    flatten.String(apiModel.DiscordWebhookUrl),
		Id:                                   flatten.Int64(apiModel.Id),
		Members:                              members,
		Name:                                 flatten.String(apiModel.Name),
		PersonalTeam:                         flatten.Bool(apiModel.PersonalTeam),
		ResendApiKey:                         flatten.String(apiModel.ResendApiKey),
		ResendEnabled:                        flatten.Bool(apiModel.ResendEnabled),
		ShowBoarding:                         flatten.Bool(apiModel.ShowBoarding),
		SmtpEnabled:                          flatten.Bool(apiModel.SmtpEnabled),
		SmtpEncryption:                       flatten.String(apiModel.SmtpEncryption),
		SmtpFromAddress:                      flatten.String(apiModel.SmtpFromAddress),
		SmtpFromName:                         flatten.String(apiModel.SmtpFromName),
		SmtpHost:                             flatten.String(apiModel.SmtpHost),
		SmtpNotificationsDatabaseBackups:     flatten.Bool(apiModel.SmtpNotificationsDatabaseBackups),
		SmtpNotificationsDeployments:         flatten.Bool(apiModel.SmtpNotificationsDeployments),
		SmtpNotificationsScheduledTasks:      flatten.Bool(apiModel.SmtpNotificationsScheduledTasks),
		SmtpNotificationsServerDiskUsage:     flatten.Bool(apiModel.SmtpNotificationsServerDiskUsage),
		SmtpNotificationsStatusChanges:       flatten.Bool(apiModel.SmtpNotificationsStatusChanges),
		SmtpNotificationsTest:                flatten.Bool(apiModel.SmtpNotificationsTest),
		SmtpPassword:                         flatten.String(apiModel.SmtpPassword),
		SmtpPort:                             flatten.String(apiModel.SmtpPort),
		SmtpRecipients:                       flatten.String(apiModel.SmtpRecipients),
		SmtpTimeout:                          flatten.String(apiModel.SmtpTimeout),
		SmtpUsername:                         flatten.String(apiModel.SmtpUsername),
		TelegramChatId:                       flatten.String(apiModel.TelegramChatId),
		TelegramEnabled:                      flatten.Bool(apiModel.TelegramEnabled),
		TelegramNotificationsDatabaseBackups: flatten.Bool(apiModel.TelegramNotificationsDatabaseBackups),
		TelegramNotificationsDatabaseBackupsMessageThreadId: flatten.String(apiModel.TelegramNotificationsDatabaseBackupsMessageThreadId),
		TelegramNotificationsDeployments:                    flatten.Bool(apiModel.TelegramNotificationsDeployments),
		TelegramNotificationsDeploymentsMessageThreadId:     flatten.String(apiModel.TelegramNotificationsDeploymentsMessageThreadId),
		TelegramNotificationsScheduledTasks:                 flatten.Bool(apiModel.TelegramNotificationsScheduledTasks),
		TelegramNotificationsScheduledTasksThreadId:         flatten.String(apiModel.TelegramNotificationsScheduledTasksThreadId),
		TelegramNotificationsStatusChanges:                  flatten.Bool(apiModel.TelegramNotificationsStatusChanges),
		TelegramNotificationsStatusChangesMessageThreadId:   flatten.String(apiModel.TelegramNotificationsStatusChangesMessageThreadId),
		TelegramNotificationsTest:                           flatten.Bool(apiModel.TelegramNotificationsTest),
		TelegramNotificationsTestMessageThreadId:            flatten.String(apiModel.TelegramNotificationsTestMessageThreadId),
		TelegramToken:                                       flatten.String(apiModel.TelegramToken),
		UpdatedAt:                                           flatten.String(apiModel.UpdatedAt),
		UseInstanceEmailSettings:                            flatten.Bool(apiModel.UseInstanceEmailSettings),
	}
}

var teamsFilterNames = []string{"name", "description", "id", "personal_team"}

func (m teamModel) FilterAttributes() map[string]attr.Value {
	return map[string]attr.Value{
		"description":   m.Description,
		"id":            m.Id,
		"name":          m.Name,
		"personal_team": m.PersonalTeam,
	}
}
