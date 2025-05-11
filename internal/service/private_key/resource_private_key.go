package private_key

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/util"
)

var (
	_ resource.Resource                = &privateKeyResource{}
	_ resource.ResourceWithConfigure   = &privateKeyResource{}
	_ resource.ResourceWithImportState = &privateKeyResource{}
	_ resource.ResourceWithModifyPlan  = &privateKeyResource{}
)

func NewPrivateKeyResource() resource.Resource {
	return &privateKeyResource{}
}

type privateKeyResource struct {
	client *api.ClientWithResponses
}

func (r *privateKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_key"
}

func (r *privateKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create, read, update, and delete a Coolify private key resource.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Optional:   true,
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"id": schema.Int64Attribute{
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"is_git_related": schema.BoolAttribute{
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"name": schema.StringAttribute{
				Optional:   true,
				Computed:   true,
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"private_key": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"public_key": schema.StringAttribute{
				Computed:            true,
				Description:         "The public key of the private key.",
				MarkdownDescription: "The public key of the private key.",
			},
			"team_id": schema.Int64Attribute{
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"uuid": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"fingerprint": schema.StringAttribute{
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Description:         "The fingerprint of the private key.",
				MarkdownDescription: "The fingerprint of the private key.",
			},
			"created_at": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *privateKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.client, resp)
}

func (r *privateKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan privateKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating private key", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})
	createResp, err := r.client.CreatePrivateKeyWithResponse(ctx, api.CreatePrivateKeyJSONRequestBody{
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.Name.ValueStringPointer(),
		PrivateKey:  plan.PrivateKey.ValueString(),
	})

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

	data := r.readFromAPI(ctx, &resp.Diagnostics, *createResp.JSON201.Uuid)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *privateKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state privateKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading private key", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	data := r.readFromAPI(ctx, &resp.Diagnostics, state.Uuid.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *privateKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan privateKeyResourceModel
	var state privateKeyResourceModel

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

	tflog.Debug(ctx, "Updating private key", map[string]interface{}{
		"uuid": uuid,
	})
	updateResp, err := r.client.UpdatePrivateKeyWithResponse(ctx, uuid, api.UpdatePrivateKeyJSONRequestBody{
		Name:        plan.Name.ValueStringPointer(),
		Description: plan.Description.ValueStringPointer(),
		PrivateKey:  plan.PrivateKey.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating private key: uuid=%s", uuid),
			err.Error(),
		)
		return
	}

	if updateResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code updating private key",
			fmt.Sprintf("Received %s updating private key: uuid=%s. Details: %s", updateResp.Status(), uuid, updateResp.Body))
		return
	}

	data := r.readFromAPI(ctx, &resp.Diagnostics, uuid)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *privateKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state privateKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting private key", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	deleteResp, err := r.client.DeletePrivateKeyByUuidWithResponse(ctx, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete private key, got error: %s", err))
		return
	}

	if deleteResp.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code deleting private key",
			fmt.Sprintf("Received %s deleting private key: %s. Details: %s", deleteResp.Status(), state, deleteResp.Body))
		return
	}
}

func (r *privateKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *privateKeyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *privateKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() || plan == nil || state == nil {
		return
	}

	// If the private key is being updated, the fingerprint will change
	if !plan.PrivateKey.Equal(state.PrivateKey) {
		plan.Fingerprint = types.StringUnknown()
	}

	resp.Plan.Set(ctx, &plan)
}

func (r *privateKeyResource) readFromAPI(
	ctx context.Context,
	diags *diag.Diagnostics,
	uuid string,
) privateKeyResourceModel {
	readResp, err := r.client.GetPrivateKeyByUuidWithResponse(ctx, uuid)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading private key: uuid=%s", uuid),
			err.Error(),
		)
		return privateKeyResourceModel{}
	}

	if readResp.StatusCode() != http.StatusOK {
		diags.AddError(
			"Unexpected HTTP status code reading private key",
			fmt.Sprintf("Received %s for private key: uuid=%s. Details: %s", readResp.Status(), uuid, readResp.Body))
		return privateKeyResourceModel{}
	}

	return privateKeyResourceModel{}.FromAPI(readResp.JSON200)
}
