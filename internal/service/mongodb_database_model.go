package service

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/flatten"
)

type mongodbDatabaseModel struct {
	commonDatabaseModel
	MongoConf               types.String `tfsdk:"mongo_conf"`
	MongoInitdbRootUsername types.String `tfsdk:"mongo_initdb_root_username"`
	MongoInitdbRootPassword types.String `tfsdk:"mongo_initdb_root_password"`
	MongoInitdbDatabase     types.String `tfsdk:"mongo_initdb_database"`
}

func (m mongodbDatabaseModel) FromAPI(apiModel *api.Database, state mongodbDatabaseModel) (mongodbDatabaseModel, error) {
	apiModel.ValueByDiscriminator()
	db, err := apiModel.AsMongodbDatabase()
	if err != nil {
		return mongodbDatabaseModel{}, err
	}

	return mongodbDatabaseModel{
		commonDatabaseModel:     commonDatabaseModel{}.FromAPI(apiModel, state.commonDatabaseModel),
		MongoConf:               flatten.String(db.MongoConf),
		MongoInitdbRootUsername: flatten.String(db.MongoInitdbRootUsername),
		MongoInitdbRootPassword: flatten.String(db.MongoInitdbRootPassword),
		MongoInitdbDatabase:     flatten.String(db.MongoInitdbDatabase),
	}, nil
}
