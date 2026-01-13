package service

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/flatten"
)

type mariadbDatabaseModel struct {
	commonDatabaseModel
	MariadbConf         types.String `tfsdk:"mariadb_conf"`
	MariadbDatabase     types.String `tfsdk:"mariadb_database"`
	MariadbPassword     types.String `tfsdk:"mariadb_password"`
	MariadbRootPassword types.String `tfsdk:"mariadb_root_password"`
	MariadbUser         types.String `tfsdk:"mariadb_user"`
}

func (m mariadbDatabaseModel) FromAPI(apiModel *api.Database, state mariadbDatabaseModel) (mariadbDatabaseModel, error) {
	apiModel.ValueByDiscriminator()
	db, err := apiModel.AsMariadbDatabase()
	if err != nil {
		return mariadbDatabaseModel{}, err
	}

	return mariadbDatabaseModel{
		commonDatabaseModel: commonDatabaseModel{}.FromAPI(apiModel, state.commonDatabaseModel),
		MariadbConf:         flatten.String(db.MariadbConf),
		MariadbDatabase:     flatten.String(db.MariadbDatabase),
		MariadbPassword:     flatten.String(db.MariadbPassword),
		MariadbRootPassword: flatten.String(db.MariadbRootPassword),
		MariadbUser:         flatten.String(db.MariadbUser),
	}, nil
}
