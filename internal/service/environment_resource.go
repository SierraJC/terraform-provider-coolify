package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/util"
)

var (
	_ resource.Resource                = &environmentResource{}
	_ resource.ResourceWithConfigure   = &environmentResource{}
	_ resource.ResourceWithImportState = &environmentResource{}
)

func NewEnvironmentResource() resource.Resource {
	return &environmentResource{}
}

type environmentResource struct {
	client *api.ClientWithResponses
}

type environmentResourceModel struct {
	ProjectUUID types.String `tfsdk:"project_uuid"`
	Name        types.String `tfsdk:"name"`
	UUID        types.String `tfsdk:"uuid"`
}

func (r *environmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (r *environmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a Coolify environment within a project. Environments organize resources (applications, services, databases) within a project. The auto-created 'production' environment cannot be managed via this resource.",
		Attributes: map[string]schema.Attribute{
			"project_uuid": schema.StringAttribute{
				Description: "UUID of the project this environment belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the environment (e.g., 'development', 'staging'). Cannot be changed after creation.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"uuid": schema.StringAttribute{
				Description: "UUID of the environment.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *environmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.client, resp)
}

func (r *environmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan environmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating environment", map[string]interface{}{
		"project_uuid": plan.ProjectUUID.ValueString(),
		"name":         plan.Name.ValueString(),
	})

	createResp, err := r.client.CreateEnvironmentWithResponse(ctx, plan.ProjectUUID.ValueString(), api.CreateEnvironmentJSONRequestBody{
		Name: plan.Name.ValueStringPointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating environment",
			err.Error(),
		)
		return
	}

	if createResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code creating environment",
			fmt.Sprintf("Received %s creating environment. Details: %s", createResp.Status(), createResp.Body),
		)
		return
	}

	if createResp.JSON201 == nil || createResp.JSON201.Uuid == nil {
		resp.Diagnostics.AddError(
			"Invalid API response",
			"API did not return environment UUID",
		)
		return
	}

	plan.UUID = types.StringValue(*createResp.JSON201.Uuid)

	tflog.Debug(ctx, "Created environment", map[string]interface{}{
		"uuid": plan.UUID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *environmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state environmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading environment", map[string]interface{}{
		"project_uuid": state.ProjectUUID.ValueString(),
		"uuid":         state.UUID.ValueString(),
	})

	readResp, err := r.client.GetEnvironmentByNameOrUuidWithResponse(ctx,
		state.ProjectUUID.ValueString(),
		state.UUID.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading environment",
			err.Error(),
		)
		return
	}

	if readResp.StatusCode() == http.StatusNotFound {
		tflog.Warn(ctx, "Environment not found, removing from state")
		resp.State.RemoveResource(ctx)
		return
	}

	if readResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading environment",
			fmt.Sprintf("Received %s reading environment. Details: %s", readResp.Status(), readResp.Body),
		)
		return
	}

	if readResp.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Invalid API response",
			"API did not return environment data",
		)
		return
	}

	// Update state with values from API
	if readResp.JSON200.Uuid != nil {
		state.UUID = types.StringValue(*readResp.JSON200.Uuid)
	}
	if readResp.JSON200.Name != nil {
		state.Name = types.StringValue(*readResp.JSON200.Name)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *environmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No update endpoint exists - all attributes have RequiresReplace plan modifiers
	resp.Diagnostics.AddError(
		"Update not supported",
		"Environment update is not supported by the Coolify API. Changes to project_uuid or name require recreation.",
	)
}

func (r *environmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state environmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting environment", map[string]interface{}{
		"project_uuid": state.ProjectUUID.ValueString(),
		"uuid":         state.UUID.ValueString(),
	})

	deleteResp, err := r.client.DeleteEnvironmentWithResponse(ctx,
		state.ProjectUUID.ValueString(),
		state.UUID.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting environment",
			err.Error(),
		)
		return
	}

	if deleteResp.StatusCode() == http.StatusNotFound {
		// Already deleted
		return
	}

	if deleteResp.StatusCode() == http.StatusBadRequest {
		resp.Diagnostics.AddError(
			"Cannot delete environment",
			"Environment contains resources and cannot be deleted. Remove all resources from this environment first.",
		)
		return
	}

	if deleteResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code deleting environment",
			fmt.Sprintf("Received %s deleting environment. Details: %s", deleteResp.Status(), deleteResp.Body),
		)
		return
	}

	tflog.Debug(ctx, "Deleted environment successfully")
}

func (r *environmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: project_uuid/environment_name_or_uuid
	parts := strings.Split(req.ID, "/")

	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format: project_uuid/environment_name_or_uuid",
		)
		return
	}

	projectUUID := parts[0]
	envIdentifier := parts[1]

	tflog.Debug(ctx, "Importing environment", map[string]interface{}{
		"project_uuid":     projectUUID,
		"env_identifier":   envIdentifier,
	})

	// Read the environment to get full details
	readResp, err := r.client.GetEnvironmentByNameOrUuidWithResponse(ctx, projectUUID, envIdentifier)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing environment",
			err.Error(),
		)
		return
	}

	if readResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Environment not found",
			fmt.Sprintf("Could not find environment '%s' in project '%s'. Status: %s", envIdentifier, projectUUID, readResp.Status()),
		)
		return
	}

	if readResp.JSON200 == nil || readResp.JSON200.Uuid == nil || readResp.JSON200.Name == nil {
		resp.Diagnostics.AddError(
			"Invalid API response",
			"API did not return complete environment data",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_uuid"), projectUUID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), *readResp.JSON200.Uuid)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), *readResp.JSON200.Name)...)
}
