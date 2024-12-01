package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
)

type mysqlDatabaseModel struct {
	commonDatabaseModel
	MysqlConf         types.String `tfsdk:"mysql_conf"`
	MysqlDatabase     types.String `tfsdk:"mysql_database"`
	MysqlPassword     types.String `tfsdk:"mysql_password"`
	MysqlRootPassword types.String `tfsdk:"mysql_root_password"`
	MysqlUser         types.String `tfsdk:"mysql_user"`
}

func (m mysqlDatabaseModel) FromAPI(apiModel *api.Database, state mysqlDatabaseModel) (mysqlDatabaseModel, error) {
	apiModel.ValueByDiscriminator()
	db, err := apiModel.AsMysqlDatabase()
	if err != nil {
		return mysqlDatabaseModel{}, err
	}

	return mysqlDatabaseModel{
		commonDatabaseModel: commonDatabaseModel{}.FromAPI(apiModel, state.commonDatabaseModel),
		MysqlConf:           optionalString(db.MysqlConf),
		MysqlDatabase:       optionalString(db.MysqlDatabase),
		MysqlPassword:       optionalString(db.MysqlPassword),
		MysqlRootPassword:   optionalString(db.MysqlRootPassword),
		MysqlUser:           optionalString(db.MysqlUser),
	}, nil
}
