package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/flatten"
	"terraform-provider-coolify/internal/provider/generated/resource_application_envs"
	"terraform-provider-coolify/internal/provider/util"
)

var (
	_ resource.Resource                = &applicationEnvsResource{}
	_ resource.ResourceWithConfigure   = &applicationEnvsResource{}
	_ resource.ResourceWithImportState = &applicationEnvsResource{}
)

func NewApplicationEnvsResource() resource.Resource {
	return &applicationEnvsResource{}
}

type applicationEnvsResource struct {
	client *api.ClientWithResponses
}

type applicationEnvsResourceModel struct {
	Uuid types.String                                     `tfsdk:"uuid"`
	Env  []resource_application_envs.ApplicationEnvsModel `tfsdk:"env"`
}

// Type alias for the anonymous struct used in the generated API code
type updateEnvsByApplicationUuidJSONRequestBodyItem = struct {
	IsBuildTime *bool   `json:"is_build_time,omitempty"`
	IsLiteral   *bool   `json:"is_literal,omitempty"`
	IsMultiline *bool   `json:"is_multiline,omitempty"`
	IsPreview   *bool   `json:"is_preview,omitempty"`
	IsShownOnce *bool   `json:"is_shown_once,omitempty"`
	Key         *string `json:"key,omitempty"`
	Value       *string `json:"value,omitempty"`
}

func (r *applicationEnvsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_envs"
}

func (r *applicationEnvsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	codegenSchema := resource_application_envs.ApplicationEnvsResourceSchema(ctx)

	resp.Schema = schema.Schema{
		Description: "Create, read, update, and delete Application environment variables.",

		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the application.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"env": schema.ListNestedBlock{
				MarkdownDescription: "Environment variable to set.",
				NestedObject: schema.NestedBlockObject{
					Attributes: codegenSchema.Attributes,
				},
			},
		},
	}

	makeResourceAttributeRequired(codegenSchema.Attributes, "key")
	makeResourceAttributeRequired(codegenSchema.Attributes, "value")
}

func (r *applicationEnvsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.client, resp)
}

func (r *applicationEnvsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan applicationEnvsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating application envs", map[string]interface{}{
		"uuid": plan.Uuid.ValueString(),
	})

	uuid := plan.Uuid.ValueString()
	for i, env := range plan.Env {
		createResp, err := r.client.CreateEnvByApplicationUuidWithResponse(ctx, uuid, api.CreateEnvByApplicationUuidJSONRequestBody{
			IsBuildTime: env.IsBuildTime.ValueBoolPointer(),
			IsLiteral:   env.IsLiteral.ValueBoolPointer(),
			IsPreview:   env.IsPreview.ValueBoolPointer(),
			Key:         env.Key.ValueStringPointer(),
			Value:       env.Value.ValueStringPointer(),
		})

		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating application envs",
				err.Error(),
			)
			return
		}

		if createResp.StatusCode() != http.StatusCreated {
			resp.Diagnostics.AddError(
				"Unexpected HTTP status code creating application envs",
				fmt.Sprintf("Received %s creating application envs. Details: %s", createResp.Status(), createResp.Body),
			)
			return
		}

		// Set the UUID from the response
		plan.Env[i].Uuid = types.StringPointerValue(createResp.JSON201.Uuid)
	}

	data, _ := r.readFromAPI(ctx, &resp.Diagnostics, uuid)
	data.Env = r.filterRelevantEnvs(plan.Env, data.Env)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *applicationEnvsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state applicationEnvsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading application envs", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	data, ok := r.readFromAPI(ctx, &resp.Diagnostics, state.Uuid.ValueString())
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}
	if len(state.Env) > 0 {
		data.Env = r.filterRelevantEnvs(state.Env, data.Env)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *applicationEnvsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan applicationEnvsResourceModel
	var state applicationEnvsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := plan.Uuid.ValueString()

	// Update API call logic
	tflog.Debug(ctx, "Updating application envs", map[string]interface{}{
		"uuid": uuid,
	})

	// Create a map of current state envs for fast lookup
	stateEnvs := make(map[string]resource_application_envs.ApplicationEnvsModel)
	for _, env := range state.Env {
		key := fmt.Sprintf("%s-%t", env.Key.ValueString(), env.IsPreview.ValueBool())
		stateEnvs[key] = env
	}

	// Create a map of plan envs for fast lookup
	planEnvs := make(map[string]resource_application_envs.ApplicationEnvsModel)
	for _, env := range plan.Env {
		key := fmt.Sprintf("%s-%t", env.Key.ValueString(), env.IsPreview.ValueBool())
		planEnvs[key] = env
	}

	// Delete envs that are in state but not in plan
	for key, env := range stateEnvs {
		if _, exists := planEnvs[key]; !exists {
			_, err := r.client.DeleteEnvByApplicationUuidWithResponse(ctx, uuid, env.Uuid.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error deleting application env: key=%s, uuid=%s", key, uuid),
					err.Error(),
				)
				return
			}
		}
	}

	var bulkUpdateEnvs = []updateEnvsByApplicationUuidJSONRequestBodyItem{}
	for _, env := range plan.Env {
		bulkUpdateEnvs = append(bulkUpdateEnvs, updateEnvsByApplicationUuidJSONRequestBodyItem{
			IsBuildTime: env.IsBuildTime.ValueBoolPointer(),
			IsLiteral:   env.IsLiteral.ValueBoolPointer(),
			IsPreview:   env.IsPreview.ValueBoolPointer(),
			Key:         env.Key.ValueStringPointer(),
			Value:       env.Value.ValueStringPointer(),
			IsMultiline: env.IsMultiline.ValueBoolPointer(),
			IsShownOnce: env.IsShownOnce.ValueBoolPointer(),
		})
	}

	if len(bulkUpdateEnvs) > 0 {
		updateResp, err := r.client.UpdateEnvsByApplicationUuidWithResponse(ctx, uuid, api.UpdateEnvsByApplicationUuidJSONRequestBody{
			Data: bulkUpdateEnvs,
		})

		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error updating application envs: uuid=%s", uuid),
				err.Error(),
			)
			return
		}

		if updateResp.StatusCode() != http.StatusCreated {
			resp.Diagnostics.AddError(
				"Unexpected HTTP status code updating application envs",
				fmt.Sprintf("Received %s updating application envs: uuid=%s. Details: %s", updateResp.Status(), uuid, updateResp.Body))
			return
		}
	}

	data, ok := r.readFromAPI(ctx, &resp.Diagnostics, uuid)
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}
	data.Env = r.filterRelevantEnvs(plan.Env, data.Env)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *applicationEnvsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state applicationEnvsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting application envs", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})

	for _, env := range state.Env {
		resp.Diagnostics.Append(r.deleteFromAPI(ctx, state.Uuid.ValueString(), env.Uuid.ValueString())...)
	}
}

func (r *applicationEnvsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

// MARK: Helper Functions

func (r *applicationEnvsResource) filterRelevantEnvs(
	stateEnvs []resource_application_envs.ApplicationEnvsModel,
	apiEnvs []resource_application_envs.ApplicationEnvsModel,
) []resource_application_envs.ApplicationEnvsModel {
	apiEnvMap := make(map[string]resource_application_envs.ApplicationEnvsModel)
	for _, env := range apiEnvs {
		key := fmt.Sprintf("%s-%t", env.Key.ValueString(), env.IsPreview.ValueBool())
		apiEnvMap[key] = env
	}

	var filteredEnvs []resource_application_envs.ApplicationEnvsModel
	for _, env := range stateEnvs {
		key := fmt.Sprintf("%s-%t", env.Key.ValueString(), env.IsPreview.ValueBool())
		if apiEnv, exists := apiEnvMap[key]; exists {
			filteredEnvs = append(filteredEnvs, apiEnv)
		}
	}

	return filteredEnvs
}

func (r *applicationEnvsResource) deleteFromAPI(
	ctx context.Context,
	uuid string,
	envUuid string,
) (diags diag.Diagnostics) {
	_, err := r.client.DeleteEnvByApplicationUuidWithResponse(ctx, uuid, envUuid)
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Unable to delete application envs, got error: %s", err))
	}
	return diags
}

func (r *applicationEnvsResource) readFromAPI(
	ctx context.Context,
	diags *diag.Diagnostics,
	uuid string,
) (applicationEnvsResourceModel, bool) {
	readResp, err := r.client.ListEnvsByApplicationUuidWithResponse(ctx, uuid)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading application envs: uuid=%s", uuid),
			err.Error(),
		)
		return applicationEnvsResourceModel{}, false
	}

	if readResp.StatusCode() == http.StatusNotFound {
		return applicationEnvsResourceModel{}, false
	}

	if readResp.StatusCode() != http.StatusOK {
		diags.AddError(
			"Unexpected HTTP status code reading application envs",
			fmt.Sprintf("Received %s for application envs: uuid=%s. Details: %s", readResp.Status(), uuid, readResp.Body))
		return applicationEnvsResourceModel{}, false
	}

	model := r.apiToModel(ctx, diags, readResp.JSON200)
	model.Uuid = types.StringValue(uuid)
	return model, true
}

func (r *applicationEnvsResource) apiToModel(
	_ context.Context,
	_ *diag.Diagnostics,
	response *[]api.EnvironmentVariable,
) applicationEnvsResourceModel {
	envs := make([]resource_application_envs.ApplicationEnvsModel, len(*response))
	for i, env := range *response {
		envs[i] = resource_application_envs.ApplicationEnvsModel{
			IsBuildTime: flatten.Bool(env.IsBuildTime),
			IsLiteral:   flatten.Bool(env.IsLiteral),
			IsMultiline: flatten.Bool(env.IsMultiline),
			IsPreview:   flatten.Bool(env.IsPreview),
			IsShownOnce: flatten.Bool(env.IsShownOnce),
			Key:         flatten.String(env.Key),
			Uuid:        flatten.String(env.Uuid),
			Value:       flatten.String(env.Value),
		}
	}

	return applicationEnvsResourceModel{
		Uuid: types.StringUnknown(),
		Env:  envs,
	}
}
