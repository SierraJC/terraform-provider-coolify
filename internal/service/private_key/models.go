package private_key

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/filter"
	"terraform-provider-coolify/internal/flatten"
)

type privateKeyModel struct {
	Description  types.String `tfsdk:"description"`
	Fingerprint  types.String `tfsdk:"fingerprint"`
	Id           types.Int64  `tfsdk:"id"`
	IsGitRelated types.Bool   `tfsdk:"is_git_related"`
	Name         types.String `tfsdk:"name"`
	PrivateKey   types.String `tfsdk:"private_key"`
	PublicKey    types.String `tfsdk:"public_key"`
	TeamId       types.Int64  `tfsdk:"team_id"`
	Uuid         types.String `tfsdk:"uuid"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

type privateKeyResourceModel = privateKeyModel
type privateKeyDataSourceModel = privateKeyModel
type privateKeysDataSourceModel struct {
	PrivateKeys []privateKeyDataSourceModel `tfsdk:"private_keys"`
	Filter      []filter.BlockModel         `tfsdk:"filter"`
}

func (m privateKeyModel) FromAPI(apiModel *api.PrivateKey) privateKeyModel {
	return privateKeyModel{
		Description:  flatten.String(apiModel.Description),
		Fingerprint:  flatten.String(apiModel.Fingerprint),
		Id:           flatten.Int64(apiModel.Id),
		IsGitRelated: flatten.Bool(apiModel.IsGitRelated),
		Name:         flatten.String(apiModel.Name),
		PrivateKey:   flatten.String(apiModel.PrivateKey),
		PublicKey:    flatten.String(apiModel.PublicKey),
		TeamId:       flatten.Int64(apiModel.TeamId),
		Uuid:         flatten.String(apiModel.Uuid),
		CreatedAt:    flatten.String(apiModel.CreatedAt),
		UpdatedAt:    flatten.String(apiModel.UpdatedAt),
	}
}

var privateKeysFilterNames = []string{"name", "description", "team_id", "is_git_related"}

func (m privateKeyModel) FilterAttributes() map[string]attr.Value {
	return map[string]attr.Value{
		"description":    m.Description,
		"is_git_related": m.IsGitRelated,
		"name":           m.Name,
		"team_id":        m.TeamId,
	}
}
