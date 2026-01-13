package service

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/flatten"
)

type clickhouseDatabaseModel struct {
	commonDatabaseModel
	ClickhouseAdminUser     types.String `tfsdk:"clickhouse_admin_user"`
	ClickhouseAdminPassword types.String `tfsdk:"clickhouse_admin_password"`
}

func (m clickhouseDatabaseModel) FromAPI(apiModel *api.Database, state clickhouseDatabaseModel) (clickhouseDatabaseModel, error) {
	apiModel.ValueByDiscriminator()
	db, err := apiModel.AsClickhouseDatabase()
	if err != nil {
		return clickhouseDatabaseModel{}, err
	}

	return clickhouseDatabaseModel{
		commonDatabaseModel:     commonDatabaseModel{}.FromAPI(apiModel, state.commonDatabaseModel),
		ClickhouseAdminUser:     flatten.String(db.ClickhouseAdminUser),
		ClickhouseAdminPassword: flatten.String(db.ClickhouseAdminPassword),
	}, nil
}
