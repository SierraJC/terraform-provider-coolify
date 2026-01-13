package service

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coolify/internal/api"
	"terraform-provider-coolify/internal/flatten"
)

type dragonflyDatabaseModel struct {
	commonDatabaseModel
	DragonflyPassword types.String `tfsdk:"dragonfly_password"`
}

func (m dragonflyDatabaseModel) FromAPI(apiModel *api.Database, state dragonflyDatabaseModel) (dragonflyDatabaseModel, error) {
	apiModel.ValueByDiscriminator()
	db, err := apiModel.AsDragonflyDatabase()
	if err != nil {
		return dragonflyDatabaseModel{}, err
	}

	return dragonflyDatabaseModel{
		commonDatabaseModel: commonDatabaseModel{}.FromAPI(apiModel, state.commonDatabaseModel),
		DragonflyPassword:   flatten.String(db.DragonflyPassword),
	}, nil
}
