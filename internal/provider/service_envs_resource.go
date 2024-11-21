package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/generated/resource_service_envs"
	"terraform-provider-coolify/internal/provider/util"
)

var (
	_ resource.Resource                = &serviceEnvsResource{}
	_ resource.ResourceWithConfigure   = &serviceEnvsResource{}
	_ resource.ResourceWithImportState = &serviceEnvsResource{}
)

func NewServiceEnvsResource() resource.Resource {
	return &serviceEnvsResource{}
}

type serviceEnvsResource struct {
	providerData CoolifyProviderData
}

type serviceEnvsResourceModel struct {
	Uuid types.String                             `tfsdk:"uuid"`
	Env  []resource_service_envs.ServiceEnvsModel `tfsdk:"env"`
}

// Type alias for the anonymous struct used in the generated API code
type updateEnvsByServiceUuidJSONRequestBodyItem = struct {
	IsBuildTime *bool   `json:"is_build_time,omitempty"`
	IsLiteral   *bool   `json:"is_literal,omitempty"`
	IsMultiline *bool   `json:"is_multiline,omitempty"`
	IsPreview   *bool   `json:"is_preview,omitempty"`
	IsShownOnce *bool   `json:"is_shown_once,omitempty"`
	Key         *string `json:"key,omitempty"`
	Value       *string `json:"value,omitempty"`
}

func (r *serviceEnvsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_envs"
}

func (r *serviceEnvsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	codegenSchema := resource_service_envs.ServiceEnvsResourceSchema(ctx)
	// todo: upstream API bug, field 'is_preview' is not supported on services and shouldn't be used

	codegenSchema.Attributes["is_preview"] = schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		DeprecationMessage:  "This field is not supported on services and should not be used.",
		MarkdownDescription: "Not supported on services and should not be used.",
	}

	resp.Schema = schema.Schema{
		Description: "Create, read, update, and delete Service environment variables.",

		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the service.",
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

func (r *serviceEnvsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.providerData, resp)
}

func (r *serviceEnvsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serviceEnvsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating service envs", map[string]interface{}{
		"uuid": plan.Uuid.ValueString(),
	})

	uuid := plan.Uuid.ValueString()
	for i, env := range plan.Env {
		createResp, err := r.providerData.client.CreateEnvByServiceUuidWithResponse(ctx, uuid, api.CreateEnvByServiceUuidJSONRequestBody{
			IsBuildTime: env.IsBuildTime.ValueBoolPointer(),
			IsLiteral:   env.IsLiteral.ValueBoolPointer(),
			IsPreview:   env.IsPreview.ValueBoolPointer(),
			Key:         env.Key.ValueStringPointer(),
			Value:       env.Value.ValueStringPointer(),
		})

		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating service envs",
				err.Error(),
			)
			return
		}

		if createResp.StatusCode() != http.StatusCreated {
			resp.Diagnostics.AddError(
				"Unexpected HTTP status code creating service envs",
				fmt.Sprintf("Received %s creating service envs. Details: %s", createResp.Status(), createResp.Body),
			)
			return
		}

		// Set the UUID from the response
		plan.Env[i].Uuid = types.StringPointerValue(createResp.JSON201.Uuid)
	}

	data := r.readFromAPI(ctx, &resp.Diagnostics, uuid)
	data.Env = r.filterRelevantEnvs(plan.Env, data.Env)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serviceEnvsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serviceEnvsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading service envs", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	data := r.readFromAPI(ctx, &resp.Diagnostics, state.Uuid.ValueString())
	if len(state.Env) > 0 {
		data.Env = r.filterRelevantEnvs(state.Env, data.Env)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serviceEnvsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan serviceEnvsResourceModel
	var state serviceEnvsResourceModel

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
	tflog.Debug(ctx, "Updating service envs", map[string]interface{}{
		"uuid": uuid,
	})

	// Create a map of current state envs for fast lookup
	stateEnvs := make(map[string]resource_service_envs.ServiceEnvsModel)
	for _, env := range state.Env {
		key := fmt.Sprintf("%s-%t", env.Key.ValueString(), env.IsPreview.ValueBool())
		stateEnvs[key] = env
	}

	// Create a map of plan envs for fast lookup
	planEnvs := make(map[string]resource_service_envs.ServiceEnvsModel)
	for _, env := range plan.Env {
		key := fmt.Sprintf("%s-%t", env.Key.ValueString(), env.IsPreview.ValueBool())
		planEnvs[key] = env
	}

	// Delete envs that are in state but not in plan
	for key, env := range stateEnvs {
		if _, exists := planEnvs[key]; !exists {
			_, err := r.providerData.client.DeleteEnvByServiceUuidWithResponse(ctx, uuid, env.Uuid.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error deleting service env: key=%s, uuid=%s", key, uuid),
					err.Error(),
				)
				return
			}
		}
	}

	var bulkUpdateEnvs = []updateEnvsByServiceUuidJSONRequestBodyItem{}
	for _, env := range plan.Env {
		bulkUpdateEnvs = append(bulkUpdateEnvs, updateEnvsByServiceUuidJSONRequestBodyItem{
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
		updateResp, err := r.providerData.client.UpdateEnvsByServiceUuidWithResponse(ctx, uuid, api.UpdateEnvsByServiceUuidJSONRequestBody{
			Data: bulkUpdateEnvs,
		})

		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error updating service envs: uuid=%s", uuid),
				err.Error(),
			)
			return
		}

		if updateResp.StatusCode() != http.StatusCreated {
			resp.Diagnostics.AddError(
				"Unexpected HTTP status code updating service envs",
				fmt.Sprintf("Received %s updating service envs: uuid=%s. Details: %s", updateResp.Status(), uuid, updateResp.Body))
			return
		}
	}

	data := r.readFromAPI(ctx, &resp.Diagnostics, uuid)
	data.Env = r.filterRelevantEnvs(plan.Env, data.Env)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serviceEnvsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serviceEnvsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting service envs", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})

	for _, env := range state.Env {
		resp.Diagnostics.Append(r.deleteFromAPI(ctx, state.Uuid.ValueString(), env.Uuid.ValueString())...)
	}
}

func (r *serviceEnvsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

// MARK: Helper Functions

func (r *serviceEnvsResource) filterRelevantEnvs(
	stateEnvs []resource_service_envs.ServiceEnvsModel,
	apiEnvs []resource_service_envs.ServiceEnvsModel,
) []resource_service_envs.ServiceEnvsModel {
	apiEnvMap := make(map[string]resource_service_envs.ServiceEnvsModel)
	for _, env := range apiEnvs {
		key := fmt.Sprintf("%s-%t", env.Key.ValueString(), env.IsPreview.ValueBool())
		apiEnvMap[key] = env
	}

	var filteredEnvs []resource_service_envs.ServiceEnvsModel
	for _, env := range stateEnvs {
		key := fmt.Sprintf("%s-%t", env.Key.ValueString(), env.IsPreview.ValueBool())
		if apiEnv, exists := apiEnvMap[key]; exists {
			filteredEnvs = append(filteredEnvs, apiEnv)
		}
	}

	return filteredEnvs
}

func (r *serviceEnvsResource) deleteFromAPI(
	ctx context.Context,
	uuid string,
	envUuid string,
) (diags diag.Diagnostics) {
	_, err := r.providerData.client.DeleteEnvByServiceUuidWithResponse(ctx, uuid, envUuid)
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Unable to delete service envs, got error: %s", err))
	}
	return diags
}

func (r *serviceEnvsResource) readFromAPI(
	ctx context.Context,
	diags *diag.Diagnostics,
	uuid string,
) serviceEnvsResourceModel {
	readResp, err := r.providerData.client.ListEnvsByServiceUuidWithResponse(ctx, uuid)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading service envs: uuid=%s", uuid),
			err.Error(),
		)
		return serviceEnvsResourceModel{}
	}

	if readResp.StatusCode() != http.StatusOK {
		diags.AddError(
			"Unexpected HTTP status code reading service envs",
			fmt.Sprintf("Received %s for service envs: uuid=%s. Details: %s", readResp.Status(), uuid, readResp.Body))
		return serviceEnvsResourceModel{}
	}

	model := r.apiToModel(ctx, diags, readResp.JSON200)
	model.Uuid = types.StringValue(uuid)
	return model
}

func (r *serviceEnvsResource) apiToModel(
	_ context.Context,
	_ *diag.Diagnostics,
	response *[]api.EnvironmentVariable,
) serviceEnvsResourceModel {
	envs := make([]resource_service_envs.ServiceEnvsModel, len(*response))
	for i, env := range *response {
		envs[i] = resource_service_envs.ServiceEnvsModel{
			IsBuildTime: optionalBool(env.IsBuildTime),
			IsLiteral:   optionalBool(env.IsLiteral),
			IsMultiline: optionalBool(env.IsMultiline),
			IsPreview:   optionalBool(env.IsPreview),
			IsShownOnce: optionalBool(env.IsShownOnce),
			Key:         optionalString(env.Key),
			Uuid:        optionalString(env.Uuid),
			Value:       optionalString(env.Value),
		}
	}

	return serviceEnvsResourceModel{
		Uuid: types.StringUnknown(),
		Env:  envs,
	}
}
