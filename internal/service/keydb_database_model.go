package service

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/flatten"
)

type keydbDatabaseModel struct {
	commonDatabaseModel
	KeydbPassword types.String `tfsdk:"keydb_password"`
	KeydbConf     types.String `tfsdk:"keydb_conf"`
}

func (m keydbDatabaseModel) FromAPI(apiModel *api.Database, state keydbDatabaseModel) (keydbDatabaseModel, error) {
	apiModel.ValueByDiscriminator()
	db, err := apiModel.AsKeydbDatabase()
	if err != nil {
		return keydbDatabaseModel{}, err
	}

	return keydbDatabaseModel{
		commonDatabaseModel: commonDatabaseModel{}.FromAPI(apiModel, state.commonDatabaseModel),
		KeydbPassword:       flatten.String(db.KeydbPassword),
		KeydbConf:           flatten.String(db.KeydbConf),
	}, nil
}
