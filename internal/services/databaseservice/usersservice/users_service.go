package usersservice

import (
	"context"
	"heypay-cash-in-server/internal/services/databaseservice"
	"heypay-cash-in-server/models/user"
)

var (
	GetDocumentsList func(ctx context.Context, queryFieldConfigs map[string]databaseservice.QueryFieldConfig, queryOptions *databaseservice.QueryOptions) ([]user.User, error)
)

func init() {
	GetDocumentsList = getDocumentsList
}

// TODO: Currently generic types are experimental, this should be refactored later on

func getDocumentsList(ctx context.Context, queryFieldConfigs map[string]databaseservice.QueryFieldConfig, queryOptions *databaseservice.QueryOptions) ([]user.User, error) {
	snapshots, err := databaseservice.GetSnapshots(ctx, "users", queryFieldConfigs, queryOptions)
	if err != nil {
		return make([]user.User, 0), err
	}
	var documentsList []user.User
	for _, snapshot := range snapshots {
		var data user.User
		err = snapshot.DataTo(&data)
		if err != nil {
			return nil, err
		}
		data.ID = snapshot.Ref.ID
		documentsList = append(documentsList, data)
	}
	return documentsList, nil
}
