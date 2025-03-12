package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/provider/util"
	sutil "terraform-provider-coolify/internal/service/util"
)

var (
	_ resource.Resource                = &postgresqlDatabaseResource{}
	_ resource.ResourceWithConfigure   = &postgresqlDatabaseResource{}
	_ resource.ResourceWithImportState = &postgresqlDatabaseResource{}
	_ resource.ResourceWithModifyPlan  = &postgresqlDatabaseResource{}
)

type postgresqlDatabaseResourceModel = postgresqlDatabaseModel

func NewPostgresqlDatabaseResource() resource.Resource {
	return &postgresqlDatabaseResource{}
}

type postgresqlDatabaseResource struct {
	client *api.ClientWithResponses
}

func (r *postgresqlDatabaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql_database"
}

func (r *postgresqlDatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	commonSchema := commonDatabaseModel{}.CommonSchema(ctx)
	postgresqlSchema := schema.Schema{
		Description: "Create, read, update, and delete a Coolify database (PostgreSQL) resource.",
		Attributes: map[string]schema.Attribute{
			"postgres_conf": schema.StringAttribute{
				Optional:    true,
				Description: "PostgreSQL conf",
			},
			"postgres_db": schema.StringAttribute{
				Required:    true,
				Description: "PostgreSQL database",
			},
			"postgres_host_auth_method": schema.StringAttribute{
				Optional:    true,
				Description: "PostgreSQL host auth method",
			},
			"postgres_initdb_args": schema.StringAttribute{
				Optional:    true,
				Description: "PostgreSQL initdb args",
			},
			"postgres_password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "PostgreSQL password",
			},
			"postgres_user": schema.StringAttribute{
				Required:    true,
				Description: "PostgreSQL user",
			},
		},
	}

	resp.Schema = sutil.MergeResourceSchemas(commonSchema, postgresqlSchema)
}

func (r *postgresqlDatabaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	util.ProviderDataFromResourceConfigureRequest(req, &r.client, resp)
}

func (r *postgresqlDatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan postgresqlDatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating postgresql database", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	createResp, err := r.client.CreateDatabasePostgresqlWithResponse(ctx, api.CreateDatabasePostgresqlJSONRequestBody{
		Description:     plan.Description.ValueStringPointer(),
		Name:            plan.Name.ValueStringPointer(),
		DestinationUuid: plan.DestinationUuid.ValueStringPointer(),
		EnvironmentName: plan.EnvironmentName.ValueString(),
		EnvironmentUuid: plan.EnvironmentUuid.ValueString(),
		Image:           plan.Image.ValueStringPointer(),
		InstantDeploy:   plan.InstantDeploy.ValueBoolPointer(),
		IsPublic:        plan.IsPublic.ValueBoolPointer(),
		LimitsCpuShares: func() *int {
			if plan.LimitsCpuShares.IsUnknown() || plan.LimitsCpuShares.IsNull() {
				return nil
			}
			value := int(*plan.LimitsCpuShares.ValueInt64Pointer())
			return &value
		}(),
		LimitsCpus:              plan.LimitsCpus.ValueStringPointer(),
		LimitsCpuset:            plan.LimitsCpuset.ValueStringPointer(),
		LimitsMemory:            plan.LimitsMemory.ValueStringPointer(),
		LimitsMemoryReservation: plan.LimitsMemoryReservation.ValueStringPointer(),
		LimitsMemorySwap:        plan.LimitsMemorySwap.ValueStringPointer(),
		LimitsMemorySwappiness: func() *int {
			if plan.LimitsMemorySwappiness.IsUnknown() || plan.LimitsMemorySwappiness.IsNull() {
				return nil
			}
			value := int(*plan.LimitsMemorySwappiness.ValueInt64Pointer())
			return &value
		}(),
		PostgresConf:           sutil.Base64EncodeAttr(plan.PostgresConf),
		PostgresDb:             plan.PostgresDb.ValueStringPointer(),
		PostgresHostAuthMethod: plan.PostgresHostAuthMethod.ValueStringPointer(),
		PostgresInitdbArgs:     plan.PostgresInitdbArgs.ValueStringPointer(),
		PostgresPassword:       plan.PostgresPassword.ValueStringPointer(),
		PostgresUser:           plan.PostgresUser.ValueStringPointer(),
		ProjectUuid:            plan.ProjectUuid.ValueString(),
		PublicPort: func() *int {
			if plan.PublicPort.IsUnknown() || plan.PublicPort.IsNull() {
				return nil
			}
			value := int(*plan.PublicPort.ValueInt64Pointer())
			return &value
		}(),
		ServerUuid: plan.ServerUuid.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating postgresql database",
			err.Error(),
		)
		return
	}

	if createResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code creating postgresql database",
			fmt.Sprintf("Received %s creating postgresql database. Details: %s", createResp.Status(), createResp.Body),
		)
		return
	}

	data := r.ReadFromAPI(ctx, &resp.Diagnostics, createResp.JSON201.Uuid, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
func (r *postgresqlDatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state postgresqlDatabaseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading postgresql database", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Invalid State", "No UUID found in state")
		return
	}

	data := r.ReadFromAPI(ctx, &resp.Diagnostics, state.Uuid.ValueString(), state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *postgresqlDatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan postgresqlDatabaseResourceModel
	var state postgresqlDatabaseResourceModel

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

	// Update API call logic
	tflog.Debug(ctx, "Updating postgresql database", map[string]interface{}{
		"uuid": uuid,
	})

	updateResp, err := r.client.UpdateDatabaseByUuidWithResponse(ctx, uuid, api.UpdateDatabaseByUuidJSONRequestBody{
		Description: plan.Description.ValueStringPointer(),
		Image:       plan.Image.ValueStringPointer(),
		IsPublic:    plan.IsPublic.ValueBoolPointer(),
		LimitsCpuShares: func() *int {
			if plan.LimitsCpuShares.IsUnknown() || plan.LimitsCpuShares.IsNull() {
				return nil
			}
			value := int(*plan.LimitsCpuShares.ValueInt64Pointer())
			return &value
		}(),
		LimitsCpus:              plan.LimitsCpus.ValueStringPointer(),
		LimitsCpuset:            plan.LimitsCpuset.ValueStringPointer(),
		LimitsMemory:            plan.LimitsMemory.ValueStringPointer(),
		LimitsMemoryReservation: plan.LimitsMemoryReservation.ValueStringPointer(),
		LimitsMemorySwap:        plan.LimitsMemorySwap.ValueStringPointer(),
		LimitsMemorySwappiness: func() *int {
			if plan.LimitsMemorySwappiness.IsUnknown() || plan.LimitsMemorySwappiness.IsNull() {
				return nil
			}
			value := int(*plan.LimitsMemorySwappiness.ValueInt64Pointer())
			return &value
		}(),
		Name:                   plan.Name.ValueStringPointer(),
		PostgresConf:           sutil.Base64EncodeAttr(plan.PostgresConf),
		PostgresDb:             plan.PostgresDb.ValueStringPointer(),
		PostgresHostAuthMethod: plan.PostgresHostAuthMethod.ValueStringPointer(),
		PostgresInitdbArgs:     plan.PostgresInitdbArgs.ValueStringPointer(),
		PostgresPassword:       plan.PostgresPassword.ValueStringPointer(),
		PostgresUser:           plan.PostgresUser.ValueStringPointer(),
		PublicPort: func() *int {
			if plan.PublicPort.IsUnknown() || plan.PublicPort.IsNull() {
				return nil
			}
			value := int(*plan.PublicPort.ValueInt64Pointer())
			return &value
		}(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating postgresql database: uuid=%s", uuid),
			err.Error(),
		)
		return
	}

	if updateResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code updating postgresql database",
			fmt.Sprintf("Received %s updating postgresql database: uuid=%s. Details: %s", updateResp.Status(), uuid, updateResp.Body))
		return
	}

	if plan.InstantDeploy.ValueBool() {
		r.client.RestartDatabaseByUuid(ctx, uuid)
	}

	data := r.ReadFromAPI(ctx, &resp.Diagnostics, uuid, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *postgresqlDatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state postgresqlDatabaseResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting postgresql database", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
	deleteResp, err := r.client.DeleteDatabaseByUuidWithResponse(ctx, state.Uuid.ValueString(), &api.DeleteDatabaseByUuidParams{
		DeleteConfigurations:    types.BoolValue(true).ValueBoolPointer(),
		DeleteVolumes:           types.BoolValue(true).ValueBoolPointer(),
		DockerCleanup:           types.BoolValue(true).ValueBoolPointer(),
		DeleteConnectedNetworks: types.BoolValue(false).ValueBoolPointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete postgresql database, got error: %s", err))
		return
	}

	if deleteResp.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unexpected HTTP status code deleting postgresql database",
			fmt.Sprintf("Received %s deleting postgresql database: %s. Details: %s", deleteResp.Status(), state, deleteResp.Body))
		return
	}
}

func (r *postgresqlDatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ids := strings.Split(req.ID, "/")
	if len(ids) != 4 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID should be in the format: <server_uuid>/<project_uuid>/<environment_name>/<database_uuid>",
		)
		return
	}

	serverUuid, projectUuid, environmentName, uuid := ids[0], ids[1], ids[2], ids[3]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("server_uuid"), serverUuid)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_uuid"), projectUuid)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_name"), environmentName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), uuid)...)
}

func (r *postgresqlDatabaseResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *postgresqlDatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() || plan == nil || state == nil {
		return
	}

	// If the username, password, or db change, the internal URL will change
	if !(plan.PostgresUser.Equal(state.PostgresUser) &&
		plan.PostgresPassword.Equal(state.PostgresPassword) &&
		plan.PostgresDb.Equal(state.PostgresDb)) {
		plan.InternalDbUrl = types.StringUnknown()
	}

	resp.Plan.Set(ctx, &plan)
}

// MARK: Helper functions

func (r *postgresqlDatabaseResource) ReadFromAPI(
	ctx context.Context,
	diags *diag.Diagnostics,
	uuid string,
	state postgresqlDatabaseResourceModel,
) postgresqlDatabaseResourceModel {
	readResp, err := r.client.GetDatabaseByUuidWithResponse(ctx, uuid)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Error reading postgresql database: uuid=%s", uuid),
			err.Error(),
		)
		return postgresqlDatabaseResourceModel{}
	}

	if readResp.StatusCode() != http.StatusOK {
		diags.AddError(
			"Unexpected HTTP status code reading postgresql database",
			fmt.Sprintf("Received %s for postgresql database: uuid=%s. Details: %s", readResp.Status(), uuid, readResp.Body))
		return postgresqlDatabaseResourceModel{}
	}

	result, err := postgresqlDatabaseResourceModel{}.FromAPI(readResp.JSON200, state)
	if err != nil {
		diags.AddError("Error converting API response to model", err.Error())
		return postgresqlDatabaseResourceModel{}
	}

	return result
}
