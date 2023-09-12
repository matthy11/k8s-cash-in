package databaseservice

import (
	"context"
	"fmt"
	"heypay-cash-in-server/internal/services/firebaseservice"
	"strings"

	"cloud.google.com/go/firestore"
)

// TODO: Currently generic types are experimental, this should be refactored later on

type QueryFieldConfig struct {
	WhereFilter string
	Value       interface{}
}

type OrderByConfig struct {
	Field     string
	Direction firestore.Direction
}

type QueryOptions struct {
	Limit       *int
	OrderBy     *[]OrderByConfig
	ParentID    *string
	Transaction *firestore.Transaction
}

type UpdateDocumentConfig struct {
	Transaction *firestore.Transaction
}

// TODO: This will be refactored when generic types become available and we implement this logic as a class service
func getPath(collection string, queryOptions *QueryOptions) string {
	if strings.Contains(collection, "/") {
		parentCollection := collection[:strings.Index(collection, "/")]
		childCollection := collection[strings.Index(collection, "/")+1:]
		if queryOptions != nil && queryOptions.ParentID != nil {
			return fmt.Sprintf("%s/%s/%s", parentCollection, *queryOptions.ParentID, childCollection)
		}
	}
	return collection
}

func GetSnapshots(ctx context.Context, collection string, queryFieldConfigs map[string]QueryFieldConfig, queryOptions *QueryOptions) ([]*firestore.DocumentSnapshot, error) {
	path := getPath(collection, queryOptions)
	collectionRef := firebaseservice.Db.Collection(path)
	var query *firestore.Query
	for field, config := range queryFieldConfigs {
		if query == nil {
			query = new(firestore.Query)
			*query = collectionRef.Where(field, config.WhereFilter, config.Value)
			continue
		}
		*query = query.Where(field, config.WhereFilter, config.Value)
	}
	if query == nil {
		return make([]*firestore.DocumentSnapshot, 0), nil
	}
	if queryOptions != nil {
		if queryOptions.Limit != nil {
			*query = query.Limit(*queryOptions.Limit)
		}
		if queryOptions.OrderBy != nil {
			for _, orderByConfig := range *queryOptions.OrderBy {
				*query = query.OrderBy(orderByConfig.Field, orderByConfig.Direction)
			}
		}
	}
	var snapshots []*firestore.DocumentSnapshot
	var err error
	if queryOptions != nil && queryOptions.Transaction != nil {
		snapshots, err = queryOptions.Transaction.Documents(query).GetAll()
	} else {
		snapshots, err = query.Documents(ctx).GetAll()
	}
	if len(snapshots) == 0 && err == nil {
		return make([]*firestore.DocumentSnapshot, 0), nil
	}
	if err != nil {
		return nil, err
	}
	return snapshots, nil
}

// UpdateDocument updates a document located on path with updateData
func UpdateDocument(ctx context.Context, path string, updateData []firestore.Update, config UpdateDocumentConfig) error {
	documentRef := firebaseservice.Db.Doc(path)
	if config.Transaction != nil {
		return config.Transaction.Update(documentRef, updateData)
	}
	_, err := documentRef.Update(ctx, updateData)
	return err
}
