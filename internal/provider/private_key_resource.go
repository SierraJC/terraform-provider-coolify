package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/generated/resource_private_key"
	"terraform-provider-coolify/internal/provider/util"
)

var _ resource.Resource = (*privateKeyResource)(nil)

var (
	_ resource.Resource                = (*privateKeyResource)(nil)
	_ resource.ResourceWithConfigure   = (*privateKeyResource)(nil)
	_ resource.ResourceWithImportState = (*privateKeyResource)(nil)
)

// type (
// 	CreateRequestBody = api.Eb4780acaa990c594cdbe8ffa80b4fb0JSONRequestBody
// 	CreateResponse    = api.Eb4780acaa990c594cdbe8ffa80b4fb0Response
// )

// func (r *privateKeyResource) createPrivateKey(ctx context.Context, body CreateRequestBody) (*CreateResponse, error) {
// 	return r.providerData.client.Eb4780acaa990c594cdbe8ffa80b4fb0WithResponse(ctx, body)
// }

func NewPrivateKeyResource() resource.Resource {
	return &privateKeyResource{}
}

type privateKeyResource struct {
	providerData CoolifyProviderData
}

func (r *privateKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_key"
}

func (r *privateKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_private_key.PrivateKeyResourceSchema(ctx)
}

func (r *privateKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.providerData, resp)
}

func (r *privateKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state resource_private_key.PrivateKeyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := api.CreatePrivateKeyJSONRequestBody{
		PrivateKey: state.PrivateKey.ValueString(),
	}
	if state.Description.ValueString() != "" {
		body.Description = state.Description.ValueStringPointer()
	}
	if state.Name.ValueString() != "" {
		body.Name = state.Name.ValueStringPointer()
	}

	tflog.Debug(ctx, "Creating private key")
	createResp, err := r.providerData.client.CreatePrivateKeyWithResponse(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating private key",
			err.Error(),
		)
		return
	}

	if createResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code creating private key",
			fmt.Sprintf("Received %s creating private key. Details: %s", createResp.Status(), createResp.Body),
		)
		return
	}

	assignStr(createResp.JSON201.Uuid, &state.Uuid)

	readResp, err := r.providerData.client.GetPrivateKeyByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading private key: uuid=%s", state.Uuid.ValueString()),
			err.Error(),
		)
		return
	}

	if readResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading private key",
			fmt.Sprintf("Received %s for private key: uuid=%s. Details: %s", readResp.Status(), state.Uuid.ValueString(), readResp.Body))
		return
	}

	assignStr(readResp.JSON200.CreatedAt, &state.CreatedAt)
	assignStr(readResp.JSON200.Description, &state.Description)
	assignInt(readResp.JSON200.Id, &state.Id)
	assignBool(readResp.JSON200.IsGitRelated, &state.IsGitRelated)
	assignStr(readResp.JSON200.Name, &state.Name)
	assignStr(readResp.JSON200.PrivateKey, &state.PrivateKey)
	assignInt(readResp.JSON200.TeamId, &state.TeamId)
	assignStr(readResp.JSON200.UpdatedAt, &state.UpdatedAt)
	assignStr(readResp.JSON200.Uuid, &state.Uuid)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *privateKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resource_private_key.PrivateKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading security key: uuid=%s", state.Uuid.ValueString()))
	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	readResp, err := r.providerData.client.GetPrivateKeyByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading private key: uuid=%s", state.Uuid.ValueString()),
			err.Error(),
		)
		return
	}

	if readResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading private key",
			fmt.Sprintf("Received %s for private key: uuid=%s. Details: %s", readResp.Status(), state.Uuid.ValueString(), readResp.Body))
		return
	}

	assignStr(readResp.JSON200.CreatedAt, &state.CreatedAt)
	assignStr(readResp.JSON200.Description, &state.Description)
	assignInt(readResp.JSON200.Id, &state.Id)
	assignBool(readResp.JSON200.IsGitRelated, &state.IsGitRelated)
	assignStr(readResp.JSON200.Name, &state.Name)
	assignStr(readResp.JSON200.PrivateKey, &state.PrivateKey)
	assignInt(readResp.JSON200.TeamId, &state.TeamId)
	assignStr(readResp.JSON200.UpdatedAt, &state.UpdatedAt)
	assignStr(readResp.JSON200.Uuid, &state.Uuid)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *privateKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state resource_private_key.PrivateKeyModel
	var plan resource_private_key.PrivateKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	// Update API call logic
	body := api.UpdatePrivateKeyJSONRequestBody{
		PrivateKey: plan.PrivateKey.ValueString(),
	}
	if state.Description.ValueString() != "" {
		body.Description = plan.Description.ValueStringPointer()
	}
	if state.Name.ValueString() != "" {
		body.Name = plan.Name.ValueStringPointer()
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating private key: uuid=%s", state.Uuid.ValueString()))
	updateResp, err := r.providerData.client.UpdatePrivateKeyWithResponse(ctx, state.Uuid.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating private key: uuid=%s", state.Uuid.ValueString()),
			err.Error(),
		)
		return
	}

	if updateResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code updating private key",
			fmt.Sprintf("Received %s updating private key: uuid=%s. Details: %s", updateResp.Status(), state.Uuid.ValueString(), updateResp.Body))
		return
	}

	assignStr(updateResp.JSON201.Uuid, &state.Uuid)

	readResp, err := r.providerData.client.GetPrivateKeyByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading private key: uuid=%s", state.Uuid.ValueString()),
			err.Error(),
		)
		return
	}

	if readResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code reading private key",
			fmt.Sprintf("Received %s for private key: uuid=%s. Details: %s", readResp.Status(), state.Uuid.ValueString(), readResp.Body))
		return
	}

	assignStr(readResp.JSON200.CreatedAt, &state.CreatedAt)
	assignStr(readResp.JSON200.Description, &state.Description)
	assignInt(readResp.JSON200.Id, &state.Id)
	assignBool(readResp.JSON200.IsGitRelated, &state.IsGitRelated)
	assignStr(readResp.JSON200.Name, &state.Name)
	assignStr(readResp.JSON200.PrivateKey, &state.PrivateKey)
	assignInt(readResp.JSON200.TeamId, &state.TeamId)
	assignStr(readResp.JSON200.UpdatedAt, &state.UpdatedAt)
	assignStr(readResp.JSON200.Uuid, &state.Uuid)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *privateKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resource_private_key.PrivateKeyModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	deleteResp, err := r.providerData.client.DeletePrivateKeyByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete private key, got error: %s", err))
		return
	}

	if deleteResp.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code deleting firewall",
			fmt.Sprintf("Received %s deleting firewall: %s. Details: %s", deleteResp.Status(), state, deleteResp.Body))
		return
	}
}

func (r *privateKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

// ---------------------------------------------------------------------- //
