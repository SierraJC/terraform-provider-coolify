package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/expand"
	"terraform-provider-coolify/internal/flatten"
)

type ApplicationSourceType string

const (
	ApplicationSourceTypePublic           ApplicationSourceType = "public"
	ApplicationSourceTypePrivateGithubApp ApplicationSourceType = "private-github-app"
	ApplicationSourceTypePrivateDeployKey ApplicationSourceType = "private-deploy-key"
	ApplicationSourceTypeDockerfile       ApplicationSourceType = "dockerfile"
	ApplicationSourceTypeDockerimage      ApplicationSourceType = "dockerimage"
	ApplicationSourceTypeDockercompose   ApplicationSourceType = "dockercompose"
)

type ApplicationModel struct {
	Uuid types.String `tfsdk:"uuid"`

	// Source type discriminant
	SourceType types.String `tfsdk:"source_type"`

	// Common fields
	ProjectUuid     types.String `tfsdk:"project_uuid"`
	ServerUuid      types.String `tfsdk:"server_uuid"`
	EnvironmentName types.String `tfsdk:"environment_name"`
	EnvironmentUuid types.String `tfsdk:"environment_uuid"`
	DestinationUuid types.String `tfsdk:"destination_uuid"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Domains         types.String `tfsdk:"domains"`
	InstantDeploy   types.Bool   `tfsdk:"instant_deploy"`

	// Git-based source fields (public, private-github-app, private-deploy-key)
	GitRepository types.String `tfsdk:"git_repository"`
	GitBranch     types.String `tfsdk:"git_branch"`
	BuildPack     types.String `tfsdk:"build_pack"`
	PortsExposes  types.String `tfsdk:"ports_exposes"`

	// Private GitHub App specific
	GithubAppUuid types.String `tfsdk:"github_app_uuid"`

	// Private Deploy Key specific
	PrivateKeyUuid types.String `tfsdk:"private_key_uuid"`

	// Dockerfile specific
	Dockerfile types.String `tfsdk:"dockerfile"`

	// Docker Image specific
	DockerRegistryImageName types.String `tfsdk:"docker_registry_image_name"`
	DockerRegistryImageTag  types.String `tfsdk:"docker_registry_image_tag"`

	// Docker Compose specific
	DockerComposeRaw types.String `tfsdk:"docker_compose_raw"`

	// Optional common fields
	BaseDirectory                   types.String `tfsdk:"base_directory"`
	BuildCommand                    types.String `tfsdk:"build_command"`
	StartCommand                    types.String `tfsdk:"start_command"`
	InstallCommand                  types.String `tfsdk:"install_command"`
	PublishDirectory                types.String `tfsdk:"publish_directory"`
	PortsMappings                   types.String `tfsdk:"ports_mappings"`
	GitCommitSha                    types.String `tfsdk:"git_commit_sha"`
	IsStatic                        types.Bool   `tfsdk:"is_static"`
	StaticImage                     types.String `tfsdk:"static_image"`
	HealthCheckEnabled              types.Bool   `tfsdk:"health_check_enabled"`
	HealthCheckPath                 types.String `tfsdk:"health_check_path"`
	HealthCheckPort                 types.String `tfsdk:"health_check_port"`
	HealthCheckHost                 types.String `tfsdk:"health_check_host"`
	HealthCheckMethod               types.String `tfsdk:"health_check_method"`
	HealthCheckReturnCode           types.Int64  `tfsdk:"health_check_return_code"`
	HealthCheckScheme               types.String `tfsdk:"health_check_scheme"`
	HealthCheckResponseText         types.String `tfsdk:"health_check_response_text"`
	HealthCheckInterval             types.Int64  `tfsdk:"health_check_interval"`
	HealthCheckTimeout              types.Int64  `tfsdk:"health_check_timeout"`
	HealthCheckRetries              types.Int64  `tfsdk:"health_check_retries"`
	HealthCheckStartPeriod          types.Int64  `tfsdk:"health_check_start_period"`
	LimitsMemory                    types.String `tfsdk:"limits_memory"`
	LimitsMemorySwap                types.String `tfsdk:"limits_memory_swap"`
	LimitsMemorySwappiness          types.Int64  `tfsdk:"limits_memory_swappiness"`
	LimitsMemoryReservation         types.String `tfsdk:"limits_memory_reservation"`
	LimitsCpus                      types.String `tfsdk:"limits_cpus"`
	LimitsCpuset                    types.String `tfsdk:"limits_cpuset"`
	LimitsCpuShares                 types.Int64  `tfsdk:"limits_cpu_shares"`
	CustomLabels                    types.String `tfsdk:"custom_labels"`
	CustomDockerRunOptions          types.String `tfsdk:"custom_docker_run_options"`
	PostDeploymentCommand           types.String `tfsdk:"post_deployment_command"`
	PostDeploymentCommandContainer  types.String `tfsdk:"post_deployment_command_container"`
	PreDeploymentCommand           types.String `tfsdk:"pre_deployment_command"`
	PreDeploymentCommandContainer  types.String `tfsdk:"pre_deployment_command_container"`
	ManualWebhookSecretGithub       types.String `tfsdk:"manual_webhook_secret_github"`
	ManualWebhookSecretGitlab       types.String `tfsdk:"manual_webhook_secret_gitlab"`
	ManualWebhookSecretBitbucket    types.String `tfsdk:"manual_webhook_secret_bitbucket"`
	ManualWebhookSecretGitea        types.String `tfsdk:"manual_webhook_secret_gitea"`
	Redirect                        types.String `tfsdk:"redirect"`
	UseBuildServer                  types.Bool   `tfsdk:"use_build_server"`
	IsHttpBasicAuthEnabled          types.Bool   `tfsdk:"is_http_basic_auth_enabled"`
	HttpBasicAuthUsername           types.String `tfsdk:"http_basic_auth_username"`
	HttpBasicAuthPassword           types.String `tfsdk:"http_basic_auth_password"`
	DockerComposeLocation            types.String `tfsdk:"docker_compose_location"`
	DockerComposeCustomStartCommand  types.String `tfsdk:"docker_compose_custom_start_command"`
	DockerComposeCustomBuildCommand  types.String `tfsdk:"docker_compose_custom_build_command"`
	WatchPaths                       types.String `tfsdk:"watch_paths"`
}

func (m ApplicationModel) Schema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description: "Create, read, update, and delete a Coolify application resource.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:      true,
				Description:   "UUID of the application.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"source_type": schema.StringAttribute{
				Required:    true,
				Description: "Type of application source. One of: public, private-github-app, private-deploy-key, dockerfile, dockerimage, dockercompose",
			},
			"project_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the project.",
			},
			"server_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the server.",
			},
			"environment_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the environment.",
			},
			"environment_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "UUID of the environment. Will replace environment_name in future.",
			},
			"destination_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "UUID of the destination if the server has multiple destinations.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the application.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Description of the application.",
			},
			"domains": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Application domains.",
			},
			"instant_deploy": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Instant deploy the application.",
				Default:     booldefault.StaticBool(false),
			},
			// Git-based fields
			"git_repository": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Git repository URL. Required for public, private-github-app, and private-deploy-key source types.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"git_branch": schema.StringAttribute{
				Optional:    true,
				Description: "Git branch. Required for public, private-github-app, and private-deploy-key source types.",
			},
			"build_pack": schema.StringAttribute{
				Optional:    true,
				Description: "Build pack type (nixpacks, static, dockerfile, dockercompose). Required for public, private-github-app, and private-deploy-key source types.",
			},
			"ports_exposes": schema.StringAttribute{
				Optional:    true,
				Description: "Ports to expose. Required for public, private-github-app, private-deploy-key, and dockerimage source types.",
			},
			// Private GitHub App specific
			"github_app_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "GitHub App UUID. Required for private-github-app source type.",
			},
			// Private Deploy Key specific
			"private_key_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "Private key UUID. Required for private-deploy-key source type.",
			},
			// Dockerfile specific
			"dockerfile": schema.StringAttribute{
				Optional:    true,
				Description: "Dockerfile content. Required for dockerfile source type.",
			},
			// Docker Image specific
			"docker_registry_image_name": schema.StringAttribute{
				Optional:    true,
				Description: "Docker registry image name. Required for dockerimage source type.",
			},
			"docker_registry_image_tag": schema.StringAttribute{
				Optional:    true,
				Description: "Docker registry image tag.",
			},
			// Docker Compose specific
			"docker_compose_raw": schema.StringAttribute{
				Optional:    true,
				Description: "Docker Compose raw content. Required for dockercompose source type.",
			},
			// Optional common fields
			"base_directory": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Base directory for all commands.",
				PlanModifiers: []planmodifier.String{
					UseStateForUnknownUnlessNullString(),
				},
			},
			"build_command":                     schema.StringAttribute{Optional: true, Description: "Build command."},
			"start_command":                     schema.StringAttribute{Optional: true, Description: "Start command."},
			"install_command":                   schema.StringAttribute{Optional: true, Description: "Install command."},
			"publish_directory":                 schema.StringAttribute{Optional: true, Description: "Publish directory."},
			"ports_mappings":                    schema.StringAttribute{Optional: true, Description: "Ports mappings."},
			"git_commit_sha": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Git commit SHA.",
				PlanModifiers: []planmodifier.String{
					UseStateForUnknownUnlessNullString(),
				},
			},
			"is_static":                         schema.BoolAttribute{Optional: true, Description: "Flag to indicate if the application is static."},
			"static_image": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Static image (e.g., nginx:alpine).",
				PlanModifiers: []planmodifier.String{
					UseStateForUnknownUnlessNullString(),
				},
			},
			"health_check_enabled": schema.BoolAttribute{Optional: true, Description: "Health check enabled."},
			"health_check_path": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Health check path.",
				PlanModifiers: []planmodifier.String{
					UseStateForUnknownUnlessNullString(),
				},
			},
			"health_check_port": schema.StringAttribute{Optional: true, Description: "Health check port."},
			"health_check_host": schema.StringAttribute{Optional: true, Description: "Health check host."},
			"health_check_method": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Health check method.",
				PlanModifiers: []planmodifier.String{
					UseStateForUnknownUnlessNullString(),
				},
			},
			"health_check_return_code": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Health check return code.",
				PlanModifiers: []planmodifier.Int64{
					UseStateForUnknownUnlessNullInt64(),
				},
			},
			"health_check_scheme": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Health check scheme.",
				PlanModifiers: []planmodifier.String{
					UseStateForUnknownUnlessNullString(),
				},
			},
			"health_check_response_text": schema.StringAttribute{Optional: true, Description: "Health check response text."},
			"health_check_interval": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Health check interval in seconds.",
				PlanModifiers: []planmodifier.Int64{
					UseStateForUnknownUnlessNullInt64(),
				},
			},
			"health_check_timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Health check timeout in seconds.",
				PlanModifiers: []planmodifier.Int64{
					UseStateForUnknownUnlessNullInt64(),
				},
			},
			"health_check_retries": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Health check retries count.",
				PlanModifiers: []planmodifier.Int64{
					UseStateForUnknownUnlessNullInt64(),
				},
			},
			"health_check_start_period": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Health check start period in seconds.",
				PlanModifiers: []planmodifier.Int64{
					UseStateForUnknownUnlessNullInt64(),
				},
			},
			"limits_memory":                      schema.StringAttribute{Optional: true, Description: "Memory limit."},
			"limits_memory_swap": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Memory swap limit.",
				PlanModifiers: []planmodifier.String{
					UseStateForUnknownUnlessNullString(),
				},
			},
			"limits_memory_swappiness": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Memory swappiness.",
				PlanModifiers: []planmodifier.Int64{
					UseStateForUnknownUnlessNullInt64(),
				},
			},
			"limits_memory_reservation": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Memory reservation.",
				PlanModifiers: []planmodifier.String{
					UseStateForUnknownUnlessNullString(),
				},
			},
			"limits_cpus":                        schema.StringAttribute{Optional: true, Description: "CPU limit."},
			"limits_cpuset":                      schema.StringAttribute{Optional: true, Description: "CPU set."},
			"limits_cpu_shares": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "CPU shares.",
				PlanModifiers: []planmodifier.Int64{
					UseStateForUnknownUnlessNullInt64(),
				},
			},
			"custom_labels": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Custom labels.",
				PlanModifiers: []planmodifier.String{
					UseStateForUnknownUnlessNullString(),
				},
			},
			"custom_docker_run_options":          schema.StringAttribute{Optional: true, Description: "Custom docker run options."},
			"post_deployment_command":            schema.StringAttribute{Optional: true, Description: "Post deployment command."},
			"post_deployment_command_container": schema.StringAttribute{Optional: true, Description: "Post deployment command container."},
			"pre_deployment_command":              schema.StringAttribute{Optional: true, Description: "Pre deployment command."},
			"pre_deployment_command_container":    schema.StringAttribute{Optional: true, Description: "Pre deployment command container."},
			"manual_webhook_secret_github":        schema.StringAttribute{Optional: true, Sensitive: true, Description: "Manual webhook secret for Github."},
			"manual_webhook_secret_gitlab":        schema.StringAttribute{Optional: true, Sensitive: true, Description: "Manual webhook secret for Gitlab."},
			"manual_webhook_secret_bitbucket":    schema.StringAttribute{Optional: true, Sensitive: true, Description: "Manual webhook secret for Bitbucket."},
			"manual_webhook_secret_gitea":         schema.StringAttribute{Optional: true, Sensitive: true, Description: "Manual webhook secret for Gitea."},
			"redirect":                            schema.StringAttribute{Optional: true, Description: "How to set redirect with Traefik / Caddy. www<->non-www."},
			"use_build_server":                    schema.BoolAttribute{Optional: true, Description: "Use build server."},
			"is_http_basic_auth_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "HTTP Basic Authentication enabled.",
				PlanModifiers: []planmodifier.Bool{
					UseStateForUnknownUnlessNullBool(),
				},
			},
			"http_basic_auth_username":            schema.StringAttribute{Optional: true, Description: "Username for HTTP Basic Authentication"},
			"http_basic_auth_password":            schema.StringAttribute{Optional: true, Sensitive: true, Description: "Password for HTTP Basic Authentication"},
			"docker_compose_location": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Docker Compose location.",
				PlanModifiers: []planmodifier.String{
					UseStateForUnknownUnlessNullString(),
				},
			},
			"docker_compose_custom_start_command": schema.StringAttribute{Optional: true, Description: "Docker Compose custom start command."},
			"docker_compose_custom_build_command": schema.StringAttribute{Optional: true, Description: "Docker Compose custom build command."},
			"watch_paths":                         schema.StringAttribute{Optional: true, Description: "Watch paths."},
		},
	}
}

func preserveGitRepository(stateVal types.String, apiVal types.String) types.String {
	if stateVal.IsNull() || stateVal.IsUnknown() {
		return apiVal
	}
	stateStr := stateVal.ValueString()
	apiStr := apiVal.ValueString()
	if stateStr != "" && apiStr != "" && stateStr != apiStr {
		if (strings.HasPrefix(stateStr, "http://") || strings.HasPrefix(stateStr, "https://") || strings.HasPrefix(stateStr, "git@")) &&
			!strings.Contains(apiStr, "://") && !strings.HasPrefix(apiStr, "git@") {
			return stateVal
		}
		return stateVal
	}
	return apiVal
}

func (m ApplicationModel) FromAPI(app *api.Application, state ApplicationModel) ApplicationModel {
	apiGitRepo := flatten.String(app.GitRepository)
	preservedGitRepo := preserveGitRepository(state.GitRepository, apiGitRepo)

	return ApplicationModel{
		Uuid:                          flatten.String(app.Uuid),
		SourceType:                    state.SourceType,
		ProjectUuid:                   state.ProjectUuid,
		ServerUuid:                    state.ServerUuid,
		EnvironmentName:                state.EnvironmentName,
		EnvironmentUuid:                state.EnvironmentUuid,
		DestinationUuid:                state.DestinationUuid,
		Name:                          flatten.String(app.Name),
		Description:                   flatten.String(app.Description),
		Domains:                       flatten.String(app.Fqdn),
		InstantDeploy:                 state.InstantDeploy,
		GitRepository:                 preservedGitRepo,
		GitBranch:                     flatten.String(app.GitBranch),
		BuildPack:                     flatten.String((*string)(app.BuildPack)),
		PortsExposes:                  flatten.String(app.PortsExposes),
		GithubAppUuid:                  state.GithubAppUuid,
		PrivateKeyUuid:                state.PrivateKeyUuid,
		Dockerfile:                    flatten.String(app.Dockerfile),
		DockerRegistryImageName:        flatten.String(app.DockerRegistryImageName),
		DockerRegistryImageTag:         flatten.String(app.DockerRegistryImageTag),
		DockerComposeRaw:               flatten.String(app.DockerComposeRaw),
		BaseDirectory:                  flatten.String(app.BaseDirectory),
		BuildCommand:                   flatten.String(app.BuildCommand),
		StartCommand:                   flatten.String(app.StartCommand),
		InstallCommand:                 flatten.String(app.InstallCommand),
		PublishDirectory:               flatten.String(app.PublishDirectory),
		PortsMappings:                  flatten.String(app.PortsMappings),
		GitCommitSha:                   flatten.String(app.GitCommitSha),
		IsStatic:                       state.IsStatic,
		StaticImage:                    flatten.String(app.StaticImage),
		HealthCheckEnabled:             flatten.Bool(app.HealthCheckEnabled),
		HealthCheckPath:                flatten.String(app.HealthCheckPath),
		HealthCheckPort:                flatten.String(app.HealthCheckPort),
		HealthCheckHost:                flatten.String(app.HealthCheckHost),
		HealthCheckMethod:              flatten.String(app.HealthCheckMethod),
		HealthCheckReturnCode:          flatten.Int64(app.HealthCheckReturnCode),
		HealthCheckScheme:               flatten.String(app.HealthCheckScheme),
		HealthCheckResponseText:        flatten.String(app.HealthCheckResponseText),
		HealthCheckInterval:            flatten.Int64(app.HealthCheckInterval),
		HealthCheckTimeout:             flatten.Int64(app.HealthCheckTimeout),
		HealthCheckRetries:             flatten.Int64(app.HealthCheckRetries),
		HealthCheckStartPeriod:         flatten.Int64(app.HealthCheckStartPeriod),
		LimitsMemory:                   flatten.String(app.LimitsMemory),
		LimitsMemorySwap:                flatten.String(app.LimitsMemorySwap),
		LimitsMemorySwappiness:         flatten.Int64(app.LimitsMemorySwappiness),
		LimitsMemoryReservation:        flatten.String(app.LimitsMemoryReservation),
		LimitsCpus:                     flatten.String(app.LimitsCpus),
		LimitsCpuset:                   flatten.String(app.LimitsCpuset),
		LimitsCpuShares:                flatten.Int64(app.LimitsCpuShares),
		CustomLabels:                   flatten.String(app.CustomLabels),
		CustomDockerRunOptions:         flatten.String(app.CustomDockerRunOptions),
		PostDeploymentCommand:          flatten.String(app.PostDeploymentCommand),
		PostDeploymentCommandContainer: flatten.String(app.PostDeploymentCommandContainer),
		PreDeploymentCommand:           flatten.String(app.PreDeploymentCommand),
		PreDeploymentCommandContainer: flatten.String(app.PreDeploymentCommandContainer),
		ManualWebhookSecretGithub:      state.ManualWebhookSecretGithub,
		ManualWebhookSecretGitlab:      state.ManualWebhookSecretGitlab,
		ManualWebhookSecretBitbucket:   state.ManualWebhookSecretBitbucket,
		ManualWebhookSecretGitea:        state.ManualWebhookSecretGitea,
		Redirect:                       flatten.String((*string)(app.Redirect)),
		UseBuildServer:                 state.UseBuildServer,
		IsHttpBasicAuthEnabled:         flatten.Bool(app.IsHttpBasicAuthEnabled),
		HttpBasicAuthUsername:          flatten.String(app.HttpBasicAuthUsername),
		HttpBasicAuthPassword:           state.HttpBasicAuthPassword,
		DockerComposeLocation:           flatten.String(app.DockerComposeLocation),
		DockerComposeCustomStartCommand: flatten.String(app.DockerComposeCustomStartCommand),
		DockerComposeCustomBuildCommand: flatten.String(app.DockerComposeCustomBuildCommand),
		WatchPaths:                     flatten.String(app.WatchPaths),
	}
}

// ToAPICreate routes to the appropriate API creation method based on source_type
func (m ApplicationModel) ToAPICreate() (interface{}, error) {
	sourceType := ApplicationSourceType(m.SourceType.ValueString())

	switch sourceType {
	case ApplicationSourceTypePublic:
		return m.toCreatePublicApplication(), nil
	case ApplicationSourceTypePrivateGithubApp:
		return m.toCreatePrivateGithubAppApplication(), nil
	case ApplicationSourceTypePrivateDeployKey:
		return m.toCreatePrivateDeployKeyApplication(), nil
	case ApplicationSourceTypeDockerfile:
		return m.toCreateDockerfileApplication(), nil
	case ApplicationSourceTypeDockerimage:
		return m.toCreateDockerimageApplication(), nil
	case ApplicationSourceTypeDockercompose:
		return m.toCreateDockercomposeApplication(), nil
	default:
		return nil, fmt.Errorf("unsupported source_type: %s", sourceType)
	}
}

func validateRedirect[T ~string](redirect *string) *T {
	if redirect == nil || *redirect == "" {
		return nil
	}
	validRedirects := map[string]bool{"both": true, "non-www": true, "www": true}
	if !validRedirects[*redirect] {
		return nil
	}
	enumVal := T(*redirect)
	return &enumVal
}

func (m ApplicationModel) toCreatePublicApplication() api.CreatePublicApplicationJSONRequestBody {
	buildPack := api.CreatePublicApplicationJSONBodyBuildPack(m.BuildPack.ValueString())
	redirect := expand.StringOrNil(m.Redirect)
	redirectEnum := validateRedirect[api.CreatePublicApplicationJSONBodyRedirect](redirect)
	staticImage := expand.String(m.StaticImage)
	var staticImageEnum *api.CreatePublicApplicationJSONBodyStaticImage
	if staticImage != nil {
		staticImageEnumVal := api.CreatePublicApplicationJSONBodyStaticImage(*staticImage)
		staticImageEnum = &staticImageEnumVal
	}

	return api.CreatePublicApplicationJSONRequestBody{
		ProjectUuid:                    m.ProjectUuid.ValueString(),
		ServerUuid:                     m.ServerUuid.ValueString(),
		EnvironmentName:                m.EnvironmentName.ValueString(),
		EnvironmentUuid:                m.EnvironmentUuid.ValueString(),
		GitRepository:                  m.GitRepository.ValueString(),
		GitBranch:                      m.GitBranch.ValueString(),
		BuildPack:                      buildPack,
		PortsExposes:                   m.PortsExposes.ValueString(),
		DestinationUuid:                expand.StringOrNil(m.DestinationUuid),
		Name:                           expand.String(m.Name),
		Description:                    expand.String(m.Description),
		Domains:                        expand.String(m.Domains),
		GitCommitSha:                   expand.String(m.GitCommitSha),
		DockerRegistryImageName:        expand.String(m.DockerRegistryImageName),
		DockerRegistryImageTag:         expand.String(m.DockerRegistryImageTag),
		IsStatic:                       expand.Bool(m.IsStatic),
		StaticImage:                    staticImageEnum,
		InstallCommand:                 expand.String(m.InstallCommand),
		BuildCommand:                   expand.String(m.BuildCommand),
		StartCommand:                   expand.String(m.StartCommand),
		PortsMappings:                  expand.String(m.PortsMappings),
		BaseDirectory:                  expand.String(m.BaseDirectory),
		PublishDirectory:               expand.String(m.PublishDirectory),
		HealthCheckEnabled:             expand.Bool(m.HealthCheckEnabled),
		HealthCheckPath:                expand.String(m.HealthCheckPath),
		HealthCheckPort:                expand.String(m.HealthCheckPort),
		HealthCheckHost:                expand.StringOrNil(m.HealthCheckHost),
		HealthCheckMethod:              expand.String(m.HealthCheckMethod),
		HealthCheckReturnCode:          expand.Int64(m.HealthCheckReturnCode),
		HealthCheckScheme:              expand.String(m.HealthCheckScheme),
		HealthCheckResponseText:        expand.String(m.HealthCheckResponseText),
		HealthCheckInterval:            expand.Int64(m.HealthCheckInterval),
		HealthCheckTimeout:             expand.Int64(m.HealthCheckTimeout),
		HealthCheckRetries:             expand.Int64(m.HealthCheckRetries),
		HealthCheckStartPeriod:         expand.Int64(m.HealthCheckStartPeriod),
		LimitsMemory:                   expand.String(m.LimitsMemory),
		LimitsMemorySwap:               expand.String(m.LimitsMemorySwap),
		LimitsMemorySwappiness:        expand.Int64(m.LimitsMemorySwappiness),
		LimitsMemoryReservation:        expand.String(m.LimitsMemoryReservation),
		LimitsCpus:                     expand.String(m.LimitsCpus),
		LimitsCpuset:                   expand.String(m.LimitsCpuset),
		LimitsCpuShares:               expand.Int64(m.LimitsCpuShares),
		CustomLabels:                   expand.String(m.CustomLabels),
		CustomDockerRunOptions:         expand.String(m.CustomDockerRunOptions),
		PostDeploymentCommand:          expand.String(m.PostDeploymentCommand),
		PostDeploymentCommandContainer: expand.String(m.PostDeploymentCommandContainer),
		PreDeploymentCommand:           expand.String(m.PreDeploymentCommand),
		PreDeploymentCommandContainer:  expand.String(m.PreDeploymentCommandContainer),
		ManualWebhookSecretGithub:      expand.String(m.ManualWebhookSecretGithub),
		ManualWebhookSecretGitlab:      expand.String(m.ManualWebhookSecretGitlab),
		ManualWebhookSecretBitbucket:    expand.String(m.ManualWebhookSecretBitbucket),
		ManualWebhookSecretGitea:       expand.String(m.ManualWebhookSecretGitea),
		Redirect:                       redirectEnum,
		InstantDeploy:                  expand.Bool(m.InstantDeploy),
		Dockerfile:                     expand.String(m.Dockerfile),
		DockerComposeLocation:          expand.String(m.DockerComposeLocation),
		DockerComposeRaw:               expand.String(m.DockerComposeRaw),
		DockerComposeCustomStartCommand: expand.String(m.DockerComposeCustomStartCommand),
		DockerComposeCustomBuildCommand: expand.String(m.DockerComposeCustomBuildCommand),
		WatchPaths:                     expand.String(m.WatchPaths),
		UseBuildServer:                 expand.Bool(m.UseBuildServer),
		IsHttpBasicAuthEnabled:         expand.Bool(m.IsHttpBasicAuthEnabled),
		HttpBasicAuthUsername:          expand.String(m.HttpBasicAuthUsername),
		HttpBasicAuthPassword:          expand.String(m.HttpBasicAuthPassword),
	}
}

func (m ApplicationModel) toCreatePrivateGithubAppApplication() api.CreatePrivateGithubAppApplicationJSONRequestBody {
	buildPack := api.CreatePrivateGithubAppApplicationJSONBodyBuildPack(m.BuildPack.ValueString())
	redirect := expand.StringOrNil(m.Redirect)
	redirectEnum := validateRedirect[api.CreatePrivateGithubAppApplicationJSONBodyRedirect](redirect)
	staticImage := expand.String(m.StaticImage)
	var staticImageEnum *api.CreatePrivateGithubAppApplicationJSONBodyStaticImage
	if staticImage != nil {
		staticImageEnumVal := api.CreatePrivateGithubAppApplicationJSONBodyStaticImage(*staticImage)
		staticImageEnum = &staticImageEnumVal
	}

	return api.CreatePrivateGithubAppApplicationJSONRequestBody{
		ProjectUuid:                    m.ProjectUuid.ValueString(),
		ServerUuid:                     m.ServerUuid.ValueString(),
		EnvironmentName:                m.EnvironmentName.ValueString(),
		EnvironmentUuid:                m.EnvironmentUuid.ValueString(),
		GithubAppUuid:                  m.GithubAppUuid.ValueString(),
		GitRepository:                  m.GitRepository.ValueString(),
		GitBranch:                      m.GitBranch.ValueString(),
		BuildPack:                      buildPack,
		PortsExposes:                   m.PortsExposes.ValueString(),
		DestinationUuid:                expand.StringOrNil(m.DestinationUuid),
		Name:                           expand.String(m.Name),
		Description:                    expand.String(m.Description),
		Domains:                        expand.String(m.Domains),
		GitCommitSha:                   expand.String(m.GitCommitSha),
		DockerRegistryImageName:        expand.String(m.DockerRegistryImageName),
		DockerRegistryImageTag:         expand.String(m.DockerRegistryImageTag),
		IsStatic:                       expand.Bool(m.IsStatic),
		StaticImage:                    staticImageEnum,
		InstallCommand:                 expand.String(m.InstallCommand),
		BuildCommand:                   expand.String(m.BuildCommand),
		StartCommand:                   expand.String(m.StartCommand),
		PortsMappings:                  expand.String(m.PortsMappings),
		BaseDirectory:                  expand.String(m.BaseDirectory),
		PublishDirectory:               expand.String(m.PublishDirectory),
		HealthCheckEnabled:             expand.Bool(m.HealthCheckEnabled),
		HealthCheckPath:                expand.String(m.HealthCheckPath),
		HealthCheckPort:                expand.String(m.HealthCheckPort),
		HealthCheckHost:                expand.StringOrNil(m.HealthCheckHost),
		HealthCheckMethod:              expand.String(m.HealthCheckMethod),
		HealthCheckReturnCode:          expand.Int64(m.HealthCheckReturnCode),
		HealthCheckScheme:              expand.String(m.HealthCheckScheme),
		HealthCheckResponseText:        expand.String(m.HealthCheckResponseText),
		HealthCheckInterval:            expand.Int64(m.HealthCheckInterval),
		HealthCheckTimeout:             expand.Int64(m.HealthCheckTimeout),
		HealthCheckRetries:             expand.Int64(m.HealthCheckRetries),
		HealthCheckStartPeriod:         expand.Int64(m.HealthCheckStartPeriod),
		LimitsMemory:                   expand.String(m.LimitsMemory),
		LimitsMemorySwap:               expand.String(m.LimitsMemorySwap),
		LimitsMemorySwappiness:        expand.Int64(m.LimitsMemorySwappiness),
		LimitsMemoryReservation:        expand.String(m.LimitsMemoryReservation),
		LimitsCpus:                     expand.String(m.LimitsCpus),
		LimitsCpuset:                   expand.String(m.LimitsCpuset),
		LimitsCpuShares:                expand.Int64(m.LimitsCpuShares),
		CustomLabels:                   expand.String(m.CustomLabels),
		CustomDockerRunOptions:         expand.String(m.CustomDockerRunOptions),
		PostDeploymentCommand:          expand.String(m.PostDeploymentCommand),
		PostDeploymentCommandContainer: expand.String(m.PostDeploymentCommandContainer),
		PreDeploymentCommand:           expand.String(m.PreDeploymentCommand),
		PreDeploymentCommandContainer:  expand.String(m.PreDeploymentCommandContainer),
		ManualWebhookSecretGithub:      expand.String(m.ManualWebhookSecretGithub),
		ManualWebhookSecretGitlab:      expand.String(m.ManualWebhookSecretGitlab),
		ManualWebhookSecretBitbucket:   expand.String(m.ManualWebhookSecretBitbucket),
		ManualWebhookSecretGitea:       expand.String(m.ManualWebhookSecretGitea),
		Redirect:                       redirectEnum,
		InstantDeploy:                  expand.Bool(m.InstantDeploy),
		Dockerfile:                     expand.String(m.Dockerfile),
		DockerComposeLocation:          expand.String(m.DockerComposeLocation),
		DockerComposeRaw:               expand.String(m.DockerComposeRaw),
		DockerComposeCustomStartCommand: expand.String(m.DockerComposeCustomStartCommand),
		DockerComposeCustomBuildCommand: expand.String(m.DockerComposeCustomBuildCommand),
		WatchPaths:                     expand.String(m.WatchPaths),
		UseBuildServer:                 expand.Bool(m.UseBuildServer),
		IsHttpBasicAuthEnabled:         expand.Bool(m.IsHttpBasicAuthEnabled),
		HttpBasicAuthUsername:          expand.String(m.HttpBasicAuthUsername),
		HttpBasicAuthPassword:           expand.String(m.HttpBasicAuthPassword),
	}
}

func (m ApplicationModel) toCreatePrivateDeployKeyApplication() api.CreatePrivateDeployKeyApplicationJSONRequestBody {
	buildPack := api.CreatePrivateDeployKeyApplicationJSONBodyBuildPack(m.BuildPack.ValueString())
	redirect := expand.StringOrNil(m.Redirect)
	redirectEnum := validateRedirect[api.CreatePrivateDeployKeyApplicationJSONBodyRedirect](redirect)
	staticImage := expand.String(m.StaticImage)
	var staticImageEnum *api.CreatePrivateDeployKeyApplicationJSONBodyStaticImage
	if staticImage != nil {
		staticImageEnumVal := api.CreatePrivateDeployKeyApplicationJSONBodyStaticImage(*staticImage)
		staticImageEnum = &staticImageEnumVal
	}

	return api.CreatePrivateDeployKeyApplicationJSONRequestBody{
		ProjectUuid:                    m.ProjectUuid.ValueString(),
		ServerUuid:                     m.ServerUuid.ValueString(),
		EnvironmentName:                m.EnvironmentName.ValueString(),
		EnvironmentUuid:                m.EnvironmentUuid.ValueString(),
		PrivateKeyUuid:                  m.PrivateKeyUuid.ValueString(),
		GitRepository:                  m.GitRepository.ValueString(),
		GitBranch:                      m.GitBranch.ValueString(),
		BuildPack:                      buildPack,
		PortsExposes:                   m.PortsExposes.ValueString(),
		DestinationUuid:                expand.StringOrNil(m.DestinationUuid),
		Name:                           expand.String(m.Name),
		Description:                    expand.String(m.Description),
		Domains:                        expand.String(m.Domains),
		GitCommitSha:                   expand.String(m.GitCommitSha),
		DockerRegistryImageName:        expand.String(m.DockerRegistryImageName),
		DockerRegistryImageTag:         expand.String(m.DockerRegistryImageTag),
		IsStatic:                       expand.Bool(m.IsStatic),
		StaticImage:                    staticImageEnum,
		InstallCommand:                 expand.String(m.InstallCommand),
		BuildCommand:                   expand.String(m.BuildCommand),
		StartCommand:                   expand.String(m.StartCommand),
		PortsMappings:                  expand.String(m.PortsMappings),
		BaseDirectory:                  expand.String(m.BaseDirectory),
		PublishDirectory:               expand.String(m.PublishDirectory),
		HealthCheckEnabled:             expand.Bool(m.HealthCheckEnabled),
		HealthCheckPath:                expand.String(m.HealthCheckPath),
		HealthCheckPort:                expand.String(m.HealthCheckPort),
		HealthCheckHost:                expand.StringOrNil(m.HealthCheckHost),
		HealthCheckMethod:              expand.String(m.HealthCheckMethod),
		HealthCheckReturnCode:          expand.Int64(m.HealthCheckReturnCode),
		HealthCheckScheme:              expand.String(m.HealthCheckScheme),
		HealthCheckResponseText:        expand.String(m.HealthCheckResponseText),
		HealthCheckInterval:            expand.Int64(m.HealthCheckInterval),
		HealthCheckTimeout:             expand.Int64(m.HealthCheckTimeout),
		HealthCheckRetries:             expand.Int64(m.HealthCheckRetries),
		HealthCheckStartPeriod:         expand.Int64(m.HealthCheckStartPeriod),
		LimitsMemory:                   expand.String(m.LimitsMemory),
		LimitsMemorySwap:               expand.String(m.LimitsMemorySwap),
		LimitsMemorySwappiness:        expand.Int64(m.LimitsMemorySwappiness),
		LimitsMemoryReservation:        expand.String(m.LimitsMemoryReservation),
		LimitsCpus:                     expand.String(m.LimitsCpus),
		LimitsCpuset:                   expand.String(m.LimitsCpuset),
		LimitsCpuShares:                expand.Int64(m.LimitsCpuShares),
		CustomLabels:                   expand.String(m.CustomLabels),
		CustomDockerRunOptions:         expand.String(m.CustomDockerRunOptions),
		PostDeploymentCommand:          expand.String(m.PostDeploymentCommand),
		PostDeploymentCommandContainer: expand.String(m.PostDeploymentCommandContainer),
		PreDeploymentCommand:           expand.String(m.PreDeploymentCommand),
		PreDeploymentCommandContainer:  expand.String(m.PreDeploymentCommandContainer),
		ManualWebhookSecretGithub:      expand.String(m.ManualWebhookSecretGithub),
		ManualWebhookSecretGitlab:      expand.String(m.ManualWebhookSecretGitlab),
		ManualWebhookSecretBitbucket:   expand.String(m.ManualWebhookSecretBitbucket),
		ManualWebhookSecretGitea:       expand.String(m.ManualWebhookSecretGitea),
		Redirect:                       redirectEnum,
		InstantDeploy:                  expand.Bool(m.InstantDeploy),
		Dockerfile:                     expand.String(m.Dockerfile),
		DockerComposeLocation:          expand.String(m.DockerComposeLocation),
		DockerComposeRaw:               expand.String(m.DockerComposeRaw),
		DockerComposeCustomStartCommand: expand.String(m.DockerComposeCustomStartCommand),
		DockerComposeCustomBuildCommand: expand.String(m.DockerComposeCustomBuildCommand),
		WatchPaths:                     expand.String(m.WatchPaths),
		UseBuildServer:                 expand.Bool(m.UseBuildServer),
		IsHttpBasicAuthEnabled:         expand.Bool(m.IsHttpBasicAuthEnabled),
		HttpBasicAuthUsername:          expand.String(m.HttpBasicAuthUsername),
		HttpBasicAuthPassword:           expand.String(m.HttpBasicAuthPassword),
	}
}

func (m ApplicationModel) toCreateDockerfileApplication() api.CreateDockerfileApplicationJSONRequestBody {
	buildPack := expand.String(m.BuildPack)
	var buildPackEnum *api.CreateDockerfileApplicationJSONBodyBuildPack
	if buildPack != nil {
		buildPackEnumVal := api.CreateDockerfileApplicationJSONBodyBuildPack(*buildPack)
		buildPackEnum = &buildPackEnumVal
	}
	redirect := expand.StringOrNil(m.Redirect)
	redirectEnum := validateRedirect[api.CreateDockerfileApplicationJSONBodyRedirect](redirect)

	return api.CreateDockerfileApplicationJSONRequestBody{
		ProjectUuid:                    m.ProjectUuid.ValueString(),
		ServerUuid:                     m.ServerUuid.ValueString(),
		EnvironmentName:                m.EnvironmentName.ValueString(),
		EnvironmentUuid:                m.EnvironmentUuid.ValueString(),
		Dockerfile:                     m.Dockerfile.ValueString(),
		BuildPack:                      buildPackEnum,
		PortsExposes:                   expand.String(m.PortsExposes),
		DestinationUuid:                expand.StringOrNil(m.DestinationUuid),
		Name:                           expand.String(m.Name),
		Description:                    expand.String(m.Description),
		Domains:                        expand.String(m.Domains),
		DockerRegistryImageName:        expand.String(m.DockerRegistryImageName),
		DockerRegistryImageTag:         expand.String(m.DockerRegistryImageTag),
		PortsMappings:                  expand.String(m.PortsMappings),
		BaseDirectory:                  expand.String(m.BaseDirectory),
		HealthCheckEnabled:             expand.Bool(m.HealthCheckEnabled),
		HealthCheckPath:                expand.String(m.HealthCheckPath),
		HealthCheckPort:                expand.String(m.HealthCheckPort),
		HealthCheckHost:                expand.StringOrNil(m.HealthCheckHost),
		HealthCheckMethod:              expand.String(m.HealthCheckMethod),
		HealthCheckReturnCode:          expand.Int64(m.HealthCheckReturnCode),
		HealthCheckScheme:              expand.String(m.HealthCheckScheme),
		HealthCheckResponseText:        expand.String(m.HealthCheckResponseText),
		HealthCheckInterval:            expand.Int64(m.HealthCheckInterval),
		HealthCheckTimeout:             expand.Int64(m.HealthCheckTimeout),
		HealthCheckRetries:             expand.Int64(m.HealthCheckRetries),
		HealthCheckStartPeriod:         expand.Int64(m.HealthCheckStartPeriod),
		LimitsMemory:                   expand.String(m.LimitsMemory),
		LimitsMemorySwap:               expand.String(m.LimitsMemorySwap),
		LimitsMemorySwappiness:        expand.Int64(m.LimitsMemorySwappiness),
		LimitsMemoryReservation:        expand.String(m.LimitsMemoryReservation),
		LimitsCpus:                     expand.String(m.LimitsCpus),
		LimitsCpuset:                   expand.String(m.LimitsCpuset),
		LimitsCpuShares:                expand.Int64(m.LimitsCpuShares),
		CustomLabels:                   expand.String(m.CustomLabels),
		CustomDockerRunOptions:         expand.String(m.CustomDockerRunOptions),
		PostDeploymentCommand:          expand.String(m.PostDeploymentCommand),
		PostDeploymentCommandContainer: expand.String(m.PostDeploymentCommandContainer),
		PreDeploymentCommand:           expand.String(m.PreDeploymentCommand),
		PreDeploymentCommandContainer:  expand.String(m.PreDeploymentCommandContainer),
		ManualWebhookSecretGithub:      expand.String(m.ManualWebhookSecretGithub),
		ManualWebhookSecretGitlab:      expand.String(m.ManualWebhookSecretGitlab),
		ManualWebhookSecretBitbucket:   expand.String(m.ManualWebhookSecretBitbucket),
		ManualWebhookSecretGitea:       expand.String(m.ManualWebhookSecretGitea),
		Redirect:                       redirectEnum,
		InstantDeploy:                  expand.Bool(m.InstantDeploy),
		UseBuildServer:                 expand.Bool(m.UseBuildServer),
	}
}

func (m ApplicationModel) toCreateDockerimageApplication() api.CreateDockerimageApplicationJSONRequestBody {
	redirect := expand.StringOrNil(m.Redirect)
	redirectEnum := validateRedirect[api.CreateDockerimageApplicationJSONBodyRedirect](redirect)

	return api.CreateDockerimageApplicationJSONRequestBody{
		ProjectUuid:                    m.ProjectUuid.ValueString(),
		ServerUuid:                     m.ServerUuid.ValueString(),
		EnvironmentName:                m.EnvironmentName.ValueString(),
		EnvironmentUuid:                m.EnvironmentUuid.ValueString(),
		DockerRegistryImageName:        m.DockerRegistryImageName.ValueString(),
		DockerRegistryImageTag:         expand.String(m.DockerRegistryImageTag),
		PortsExposes:                   m.PortsExposes.ValueString(),
		DestinationUuid:                expand.StringOrNil(m.DestinationUuid),
		Name:                           expand.String(m.Name),
		Description:                    expand.String(m.Description),
		Domains:                        expand.String(m.Domains),
		PortsMappings:                  expand.String(m.PortsMappings),
		HealthCheckEnabled:             expand.Bool(m.HealthCheckEnabled),
		HealthCheckPath:                expand.String(m.HealthCheckPath),
		HealthCheckPort:                expand.String(m.HealthCheckPort),
		HealthCheckHost:                expand.StringOrNil(m.HealthCheckHost),
		HealthCheckMethod:              expand.String(m.HealthCheckMethod),
		HealthCheckReturnCode:          expand.Int64(m.HealthCheckReturnCode),
		HealthCheckScheme:              expand.String(m.HealthCheckScheme),
		HealthCheckResponseText:        expand.String(m.HealthCheckResponseText),
		HealthCheckInterval:            expand.Int64(m.HealthCheckInterval),
		HealthCheckTimeout:             expand.Int64(m.HealthCheckTimeout),
		HealthCheckRetries:             expand.Int64(m.HealthCheckRetries),
		HealthCheckStartPeriod:         expand.Int64(m.HealthCheckStartPeriod),
		LimitsMemory:                   expand.String(m.LimitsMemory),
		LimitsMemorySwap:               expand.String(m.LimitsMemorySwap),
		LimitsMemorySwappiness:        expand.Int64(m.LimitsMemorySwappiness),
		LimitsMemoryReservation:        expand.String(m.LimitsMemoryReservation),
		LimitsCpus:                     expand.String(m.LimitsCpus),
		LimitsCpuset:                   expand.String(m.LimitsCpuset),
		LimitsCpuShares:                expand.Int64(m.LimitsCpuShares),
		CustomLabels:                   expand.String(m.CustomLabels),
		CustomDockerRunOptions:         expand.String(m.CustomDockerRunOptions),
		PostDeploymentCommand:          expand.String(m.PostDeploymentCommand),
		PostDeploymentCommandContainer: expand.String(m.PostDeploymentCommandContainer),
		PreDeploymentCommand:           expand.String(m.PreDeploymentCommand),
		PreDeploymentCommandContainer:  expand.String(m.PreDeploymentCommandContainer),
		ManualWebhookSecretGithub:      expand.String(m.ManualWebhookSecretGithub),
		ManualWebhookSecretGitlab:      expand.String(m.ManualWebhookSecretGitlab),
		ManualWebhookSecretBitbucket:   expand.String(m.ManualWebhookSecretBitbucket),
		ManualWebhookSecretGitea:       expand.String(m.ManualWebhookSecretGitea),
		Redirect:                       redirectEnum,
		InstantDeploy:                  expand.Bool(m.InstantDeploy),
	}
}

func (m ApplicationModel) toCreateDockercomposeApplication() api.CreateDockercomposeApplicationJSONRequestBody {
	return api.CreateDockercomposeApplicationJSONRequestBody{
		ProjectUuid:                    m.ProjectUuid.ValueString(),
		ServerUuid:                     m.ServerUuid.ValueString(),
		EnvironmentName:                m.EnvironmentName.ValueString(),
		EnvironmentUuid:                m.EnvironmentUuid.ValueString(),
		DockerComposeRaw:               m.DockerComposeRaw.ValueString(),
		DestinationUuid:                expand.StringOrNil(m.DestinationUuid),
		Name:                           expand.String(m.Name),
		Description:                    expand.String(m.Description),
		InstantDeploy:                  expand.Bool(m.InstantDeploy),
		UseBuildServer:                 expand.Bool(m.UseBuildServer),
	}
}

func (m ApplicationModel) ToAPIUpdate() api.UpdateApplicationByUuidJSONRequestBody {
	buildPack := expand.String(m.BuildPack)
	var buildPackEnum *api.UpdateApplicationByUuidJSONBodyBuildPack
	if buildPack != nil {
		buildPackEnumVal := api.UpdateApplicationByUuidJSONBodyBuildPack(*buildPack)
		buildPackEnum = &buildPackEnumVal
	}
	redirect := expand.StringOrNil(m.Redirect)
	redirectEnum := validateRedirect[api.UpdateApplicationByUuidJSONBodyRedirect](redirect)

	return api.UpdateApplicationByUuidJSONRequestBody{
		Description:                    expand.String(m.Description),
		DestinationUuid:                expand.StringOrNil(m.DestinationUuid),
		Domains:                        expand.String(m.Domains),
		Name:                           expand.String(m.Name),
		ProjectUuid:                    expand.String(m.ProjectUuid),
		ServerUuid:                     expand.String(m.ServerUuid),
		EnvironmentName:                expand.String(m.EnvironmentName),
		BuildPack:                      buildPackEnum,
		BaseDirectory:                  expand.String(m.BaseDirectory),
		BuildCommand:                   expand.String(m.BuildCommand),
		StartCommand:                   expand.String(m.StartCommand),
		InstallCommand:                 expand.String(m.InstallCommand),
		PublishDirectory:               expand.String(m.PublishDirectory),
		PortsMappings:                  expand.String(m.PortsMappings),
		PortsExposes:                   expand.String(m.PortsExposes),
		GitCommitSha:                   expand.String(m.GitCommitSha),
		GitBranch:                      expand.String(m.GitBranch),
		GitRepository:                  expand.String(m.GitRepository),
		GithubAppUuid:                  expand.String(m.GithubAppUuid),
		HealthCheckEnabled:             expand.Bool(m.HealthCheckEnabled),
		HealthCheckPath:                expand.String(m.HealthCheckPath),
		HealthCheckPort:                expand.String(m.HealthCheckPort),
		HealthCheckHost:                expand.StringOrNil(m.HealthCheckHost),
		HealthCheckMethod:              expand.String(m.HealthCheckMethod),
		HealthCheckReturnCode:          expand.Int64(m.HealthCheckReturnCode),
		HealthCheckScheme:              expand.String(m.HealthCheckScheme),
		HealthCheckResponseText:        expand.String(m.HealthCheckResponseText),
		HealthCheckInterval:            expand.Int64(m.HealthCheckInterval),
		HealthCheckTimeout:             expand.Int64(m.HealthCheckTimeout),
		HealthCheckRetries:             expand.Int64(m.HealthCheckRetries),
		HealthCheckStartPeriod:         expand.Int64(m.HealthCheckStartPeriod),
		LimitsMemory:                   expand.String(m.LimitsMemory),
		LimitsMemorySwap:               expand.String(m.LimitsMemorySwap),
		LimitsMemorySwappiness:        expand.Int64(m.LimitsMemorySwappiness),
		LimitsMemoryReservation:        expand.String(m.LimitsMemoryReservation),
		LimitsCpus:                     expand.String(m.LimitsCpus),
		LimitsCpuset:                   expand.String(m.LimitsCpuset),
		LimitsCpuShares:                expand.Int64(m.LimitsCpuShares),
		CustomLabels:                   expand.String(m.CustomLabels),
		CustomDockerRunOptions:         expand.String(m.CustomDockerRunOptions),
		PostDeploymentCommand:          expand.String(m.PostDeploymentCommand),
		PostDeploymentCommandContainer: expand.String(m.PostDeploymentCommandContainer),
		PreDeploymentCommand:           expand.String(m.PreDeploymentCommand),
		PreDeploymentCommandContainer:  expand.String(m.PreDeploymentCommandContainer),
		ManualWebhookSecretGithub:      expand.String(m.ManualWebhookSecretGithub),
		ManualWebhookSecretGitlab:      expand.String(m.ManualWebhookSecretGitlab),
		ManualWebhookSecretBitbucket:   expand.String(m.ManualWebhookSecretBitbucket),
		ManualWebhookSecretGitea:       expand.String(m.ManualWebhookSecretGitea),
		Redirect:                       redirectEnum,
		InstantDeploy:                  expand.Bool(m.InstantDeploy),
		Dockerfile:                     expand.String(m.Dockerfile),
		DockerComposeLocation:          expand.String(m.DockerComposeLocation),
		DockerComposeRaw:               expand.String(m.DockerComposeRaw),
		DockerComposeCustomStartCommand: expand.String(m.DockerComposeCustomStartCommand),
		DockerComposeCustomBuildCommand: expand.String(m.DockerComposeCustomBuildCommand),
		DockerComposeDomains:           nil,
		WatchPaths:                     expand.String(m.WatchPaths),
		UseBuildServer:                 expand.Bool(m.UseBuildServer),
		DockerRegistryImageName:        expand.String(m.DockerRegistryImageName),
		DockerRegistryImageTag:         expand.String(m.DockerRegistryImageTag),
	}
}

