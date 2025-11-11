package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/util"
)

var (
	_ resource.Resource                = &applicationResource{}
	_ resource.ResourceWithConfigure   = &applicationResource{}
	_ resource.ResourceWithImportState = &applicationResource{}
)

func NewApplicationResource() resource.Resource {
	return &applicationResource{}
}

type applicationResource struct {
	client *api.ClientWithResponses
}

func (r *applicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *applicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ApplicationModel{}.Schema(ctx)
	resp.Schema.Description = "Create, read, update, and delete a Coolify application resource."

	// Mark sensitive attributes
	sensitiveAttrs := []string{
		"manual_webhook_secret_bitbucket",
		"manual_webhook_secret_gitea",
		"manual_webhook_secret_github",
		"manual_webhook_secret_gitlab",
		"http_basic_auth_password",
	}
	for _, attr := range sensitiveAttrs {
		if err := makeResourceAttributeSensitive(resp.Schema.Attributes, attr); err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Failed to mark attribute as sensitive: %s", attr))
		}
	}
}

func (r *applicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.client, resp)
}

func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApplicationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate required fields based on source_type
	if diags := r.validateCreatePlan(ctx, plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	sourceType := ApplicationSourceType(plan.SourceType.ValueString())
	tflog.Debug(ctx, "Creating application", map[string]interface{}{
		"source_type": sourceType,
		"name":        plan.Name.ValueString(),
	})

	createBody, err := plan.ToAPICreate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error preparing application creation",
			err.Error(),
		)
		return
	}

	var uuid string
	switch sourceType {
	case ApplicationSourceTypePublic:
		body := createBody.(api.CreatePublicApplicationJSONRequestBody)
		apiResp, err := r.client.CreatePublicApplicationWithResponse(ctx, body)
		if err != nil {
			resp.Diagnostics.AddError("Error creating application", err.Error())
			return
		}
		if apiResp.StatusCode() != http.StatusCreated {
			resp.Diagnostics.AddError(
				"Unexpected HTTP status code creating application",
				fmt.Sprintf("Received %d creating application. Details: %s", apiResp.StatusCode(), string(apiResp.Body)),
			)
			return
		}
		if apiResp.JSON201 == nil || apiResp.JSON201.Uuid == nil {
			resp.Diagnostics.AddError(
				"Invalid response creating application",
				"Response did not contain application UUID",
			)
			return
		}
		uuid = *apiResp.JSON201.Uuid
	case ApplicationSourceTypePrivateGithubApp:
		body := createBody.(api.CreatePrivateGithubAppApplicationJSONRequestBody)
		apiResp, err := r.client.CreatePrivateGithubAppApplicationWithResponse(ctx, body)
		if err != nil {
			resp.Diagnostics.AddError("Error creating application", err.Error())
			return
		}
		if apiResp.StatusCode() != http.StatusCreated {
			resp.Diagnostics.AddError(
				"Unexpected HTTP status code creating application",
				fmt.Sprintf("Received %d creating application. Details: %s", apiResp.StatusCode(), string(apiResp.Body)),
			)
			return
		}
		if apiResp.JSON201 == nil || apiResp.JSON201.Uuid == nil {
			resp.Diagnostics.AddError(
				"Invalid response creating application",
				"Response did not contain application UUID",
			)
			return
		}
		uuid = *apiResp.JSON201.Uuid
	case ApplicationSourceTypePrivateDeployKey:
		body := createBody.(api.CreatePrivateDeployKeyApplicationJSONRequestBody)
		apiResp, err := r.client.CreatePrivateDeployKeyApplicationWithResponse(ctx, body)
		if err != nil {
			resp.Diagnostics.AddError("Error creating application", err.Error())
			return
		}
		if apiResp.StatusCode() != http.StatusCreated {
			resp.Diagnostics.AddError(
				"Unexpected HTTP status code creating application",
				fmt.Sprintf("Received %d creating application. Details: %s", apiResp.StatusCode(), string(apiResp.Body)),
			)
			return
		}
		if apiResp.JSON201 == nil || apiResp.JSON201.Uuid == nil {
			resp.Diagnostics.AddError(
				"Invalid response creating application",
				"Response did not contain application UUID",
			)
			return
		}
		uuid = *apiResp.JSON201.Uuid
	case ApplicationSourceTypeDockerfile:
		body := createBody.(api.CreateDockerfileApplicationJSONRequestBody)
		apiResp, err := r.client.CreateDockerfileApplicationWithResponse(ctx, body)
		if err != nil {
			resp.Diagnostics.AddError("Error creating application", err.Error())
			return
		}
		if apiResp.StatusCode() != http.StatusCreated {
			resp.Diagnostics.AddError(
				"Unexpected HTTP status code creating application",
				fmt.Sprintf("Received %d creating application. Details: %s", apiResp.StatusCode(), string(apiResp.Body)),
			)
			return
		}
		if apiResp.JSON201 == nil || apiResp.JSON201.Uuid == nil {
			resp.Diagnostics.AddError(
				"Invalid response creating application",
				"Response did not contain application UUID",
			)
			return
		}
		uuid = *apiResp.JSON201.Uuid
	case ApplicationSourceTypeDockerimage:
		body := createBody.(api.CreateDockerimageApplicationJSONRequestBody)
		apiResp, err := r.client.CreateDockerimageApplicationWithResponse(ctx, body)
		if err != nil {
			resp.Diagnostics.AddError("Error creating application", err.Error())
			return
		}
		if apiResp.StatusCode() != http.StatusCreated {
			resp.Diagnostics.AddError(
				"Unexpected HTTP status code creating application",
				fmt.Sprintf("Received %d creating application. Details: %s", apiResp.StatusCode(), string(apiResp.Body)),
			)
			return
		}
		if apiResp.JSON201 == nil || apiResp.JSON201.Uuid == nil {
			resp.Diagnostics.AddError(
				"Invalid response creating application",
				"Response did not contain application UUID",
			)
			return
		}
		uuid = *apiResp.JSON201.Uuid
	case ApplicationSourceTypeDockercompose:
		body := createBody.(api.CreateDockercomposeApplicationJSONRequestBody)
		apiResp, err := r.client.CreateDockercomposeApplicationWithResponse(ctx, body)
		if err != nil {
			resp.Diagnostics.AddError("Error creating application", err.Error())
			return
		}
		if apiResp.StatusCode() != http.StatusCreated {
			resp.Diagnostics.AddError(
				"Unexpected HTTP status code creating application",
				fmt.Sprintf("Received %d creating application. Details: %s", apiResp.StatusCode(), string(apiResp.Body)),
			)
			return
		}
		if apiResp.JSON201 == nil || apiResp.JSON201.Uuid == nil {
			resp.Diagnostics.AddError(
				"Invalid response creating application",
				"Response did not contain application UUID",
			)
			return
		}
		uuid = *apiResp.JSON201.Uuid
	default:
		resp.Diagnostics.AddError(
			"Unsupported source_type",
			fmt.Sprintf("source_type %s is not supported", sourceType),
		)
		return
	}

	data, ok := r.ReadFromAPI(ctx, &resp.Diagnostics, uuid, plan)
	if !ok {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ApplicationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading application", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})

	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	data, ok := r.ReadFromAPI(ctx, &resp.Diagnostics, state.Uuid.ValueString(), state)
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ApplicationModel
	var state ApplicationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()
	if uuid == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	tflog.Debug(ctx, "Updating application", map[string]interface{}{
		"uuid": uuid,
	})

	updateResp, err := r.client.UpdateApplicationByUuidWithResponse(ctx, uuid, plan.ToAPIUpdate())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating application: uuid=%s", uuid),
			err.Error(),
		)
		return
	}

	if updateResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code updating application",
			fmt.Sprintf("Received %d updating application: uuid=%s. Details: %s", updateResp.StatusCode(), uuid, string(updateResp.Body)),
		)
		return
	}

	data, ok := r.ReadFromAPI(ctx, &resp.Diagnostics, uuid, plan)
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ApplicationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()
	if uuid == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	tflog.Debug(ctx, "Deleting application", map[string]interface{}{
		"uuid": uuid,
	})

	deleteTrue := true
	deleteFalse := false
	deleteResp, err := r.client.DeleteApplicationByUuidWithResponse(ctx, uuid, &api.DeleteApplicationByUuidParams{
		DeleteConfigurations:    &deleteTrue,
		DeleteVolumes:           &deleteTrue,
		DockerCleanup:           &deleteTrue,
		DeleteConnectedNetworks: &deleteFalse,
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete application, got error: %s", err))
		return
	}

	if deleteResp.StatusCode() != http.StatusOK && deleteResp.StatusCode() != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code deleting application",
			fmt.Sprintf("Received %d deleting application: %s. Details: %s", deleteResp.StatusCode(), uuid, string(deleteResp.Body)),
		)
		return
	}
}

func (r *applicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}


// MARK: Helper Functions

func (r *applicationResource) ReadFromAPI(
	ctx context.Context,
	diags *diag.Diagnostics,
	uuid string,
	state ApplicationModel,
) (ApplicationModel, bool) {
	res, err := r.client.GetApplicationByUuidWithResponse(ctx, uuid)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading application: uuid=%s", uuid),
			err.Error(),
		)
		return ApplicationModel{}, false
	}

	if res.StatusCode() == http.StatusNotFound {
		return ApplicationModel{}, false
	}

	if res.StatusCode() != http.StatusOK {
		diags.AddError(
			"Unexpected HTTP status code reading application",
			fmt.Sprintf("Received %d for application: uuid=%s. Details: %s", res.StatusCode(), uuid, string(res.Body)),
		)
		return ApplicationModel{}, false
	}

	result := ApplicationModel{}.FromAPI(res.JSON200, state)
	return result, true
}

func (r *applicationResource) validateCreatePlan(ctx context.Context, plan ApplicationModel) diag.Diagnostics {
	var diags diag.Diagnostics
	sourceType := ApplicationSourceType(plan.SourceType.ValueString())

	switch sourceType {
	case ApplicationSourceTypePublic, ApplicationSourceTypePrivateGithubApp, ApplicationSourceTypePrivateDeployKey:
		if plan.GitRepository.IsNull() || plan.GitRepository.ValueString() == "" {
			diags.AddAttributeError(
				path.Root("git_repository"),
				"Missing required field",
				fmt.Sprintf("git_repository is required for source_type %s", sourceType),
			)
		}
		if plan.GitBranch.IsNull() || plan.GitBranch.ValueString() == "" {
			diags.AddAttributeError(
				path.Root("git_branch"),
				"Missing required field",
				fmt.Sprintf("git_branch is required for source_type %s", sourceType),
			)
		}
		if plan.BuildPack.IsNull() || plan.BuildPack.ValueString() == "" {
			diags.AddAttributeError(
				path.Root("build_pack"),
				"Missing required field",
				fmt.Sprintf("build_pack is required for source_type %s", sourceType),
			)
		}
		if plan.PortsExposes.IsNull() || plan.PortsExposes.ValueString() == "" {
			diags.AddAttributeError(
				path.Root("ports_exposes"),
				"Missing required field",
				fmt.Sprintf("ports_exposes is required for source_type %s", sourceType),
			)
		}
		if sourceType == ApplicationSourceTypePrivateGithubApp {
			if plan.GithubAppUuid.IsNull() || plan.GithubAppUuid.ValueString() == "" {
				diags.AddAttributeError(
					path.Root("github_app_uuid"),
					"Missing required field",
					"github_app_uuid is required for source_type private-github-app",
				)
			}
		}
		if sourceType == ApplicationSourceTypePrivateDeployKey {
			if plan.PrivateKeyUuid.IsNull() || plan.PrivateKeyUuid.ValueString() == "" {
				diags.AddAttributeError(
					path.Root("private_key_uuid"),
					"Missing required field",
					"private_key_uuid is required for source_type private-deploy-key",
				)
			}
		}
	case ApplicationSourceTypeDockerfile:
		if plan.Dockerfile.IsNull() || plan.Dockerfile.ValueString() == "" {
			diags.AddAttributeError(
				path.Root("dockerfile"),
				"Missing required field",
				"dockerfile is required for source_type dockerfile",
			)
		}
	case ApplicationSourceTypeDockerimage:
		if plan.DockerRegistryImageName.IsNull() || plan.DockerRegistryImageName.ValueString() == "" {
			diags.AddAttributeError(
				path.Root("docker_registry_image_name"),
				"Missing required field",
				"docker_registry_image_name is required for source_type dockerimage",
			)
		}
		if plan.PortsExposes.IsNull() || plan.PortsExposes.ValueString() == "" {
			diags.AddAttributeError(
				path.Root("ports_exposes"),
				"Missing required field",
				"ports_exposes is required for source_type dockerimage",
			)
		}
	case ApplicationSourceTypeDockercompose:
		if plan.DockerComposeRaw.IsNull() || plan.DockerComposeRaw.ValueString() == "" {
			diags.AddAttributeError(
				path.Root("docker_compose_raw"),
				"Missing required field",
				"docker_compose_raw is required for source_type dockercompose",
			)
		}
	default:
		diags.AddAttributeError(
			path.Root("source_type"),
			"Invalid source_type",
			fmt.Sprintf("source_type %s is not supported. Valid values: public, private-github-app, private-deploy-key, dockerfile, dockerimage, dockercompose", sourceType),
		)
	}

	return diags
}

