package service

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/flatten"
)

type redisDatabaseModel struct {
	commonDatabaseModel
	RedisPassword types.String `tfsdk:"redis_password"`
	RedisConf     types.String `tfsdk:"redis_conf"`
}

func (m redisDatabaseModel) FromAPI(apiModel *api.Database, state redisDatabaseModel) (redisDatabaseModel, error) {
	apiModel.ValueByDiscriminator()
	db, err := apiModel.AsRedisDatabase()
	if err != nil {
		return redisDatabaseModel{}, err
	}

	return redisDatabaseModel{
		commonDatabaseModel: commonDatabaseModel{}.FromAPI(apiModel, state.commonDatabaseModel),
		RedisPassword:       flatten.String(db.RedisPassword),
		RedisConf:           flatten.String(db.RedisConf),
	}, nil
}
