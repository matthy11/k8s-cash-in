package transferpaymentrequestsservice

import (
	"context"
	"fmt"
	"heypay-cash-in-server/internal/services/databaseservice"
	"heypay-cash-in-server/models/transferpaymentrequest"
	"sort"
	"time"

	"cloud.google.com/go/firestore"
)

var (
	FindPendingTransferPaymentRequest func(ctx context.Context, amount int64, destinationNationalID string, destinationAccountNumber string, originNationalID string, expiresAt time.Time, transaction *firestore.Transaction) (*transferpaymentrequest.TransferPaymentRequest, error)
	UpdateDocument                    func(ctx context.Context, documentID string, updateData []firestore.Update, config databaseservice.UpdateDocumentConfig) error
)

func init() {
	FindPendingTransferPaymentRequest = findPendingTransferPaymentRequest
	UpdateDocument = updateDocument
}

// findPendingTransferPaymentRequest finds the most recent pending transfer payment request.
func findPendingTransferPaymentRequest(ctx context.Context, amount int64, destinationNationalID string, destinationAccountNumber string, originNationalID string, expiresAt time.Time, transaction *firestore.Transaction) (*transferpaymentrequest.TransferPaymentRequest, error) {
	transferPaymentRequestsDocumentsList, err := GetDocumentsList(ctx, map[string]databaseservice.QueryFieldConfig{
		"expiresAt":                {WhereFilter: ">", Value: expiresAt},
		"originNationalId":         {WhereFilter: "==", Value: originNationalID},
		"destinationNationalId":    {WhereFilter: "==", Value: destinationNationalID},
		"destinationAccountNumber": {WhereFilter: "==", Value: destinationAccountNumber},
		"status":                   {WhereFilter: "==", Value: "pending"},
	}, &databaseservice.QueryOptions{Transaction: transaction})
	if err != nil {
		return nil, err
	}
	sort.Slice(transferPaymentRequestsDocumentsList, func(i, j int) bool {
		return transferPaymentRequestsDocumentsList[i].CreatedAt.After(transferPaymentRequestsDocumentsList[j].CreatedAt)
	})
	for _, transferPaymentRequestDocument := range transferPaymentRequestsDocumentsList {
		if transferPaymentRequestDocument.Amount == amount {
			return &transferPaymentRequestDocument, nil
		}
	}
	return nil, nil
}

// FindValidatedTransferPaymentRequest finds the validated transfer payment request.
func FindValidatedTransferPaymentRequest(ctx context.Context, transferID string, transaction *firestore.Transaction) (*transferpaymentrequest.TransferPaymentRequest, error) {
	transferPaymentRequestsDocumentsList, err := GetDocumentsList(ctx, map[string]databaseservice.QueryFieldConfig{
		"transferId": {WhereFilter: "==", Value: transferID},
		"status":     {WhereFilter: "==", Value: "validated"},
	}, &databaseservice.QueryOptions{Transaction: transaction})
	if err != nil {
		return nil, err
	}
	if len(transferPaymentRequestsDocumentsList) > 0 {
		return &transferPaymentRequestsDocumentsList[0], nil
	}
	return nil, nil
}

// TODO: Currently generic types are experimental, this should be refactored later on

// GetDocumentsList list transfer payment requests.
func GetDocumentsList(ctx context.Context, queryFieldConfigs map[string]databaseservice.QueryFieldConfig, queryOptions *databaseservice.QueryOptions) ([]transferpaymentrequest.TransferPaymentRequest, error) {
	snapshots, err := databaseservice.GetSnapshots(ctx, "transferPaymentRequests", queryFieldConfigs, queryOptions)
	if err != nil {
		return make([]transferpaymentrequest.TransferPaymentRequest, 0), err
	}
	var documentsList []transferpaymentrequest.TransferPaymentRequest
	for _, snapshot := range snapshots {
		var data transferpaymentrequest.TransferPaymentRequest
		err = snapshot.DataTo(&data)
		if err != nil {
			return nil, err
		}
		data.ID = snapshot.Ref.ID
		documentsList = append(documentsList, data)
	}
	return documentsList, nil
}

// updateDocument updates a transferPaymentRequests document by documentID with updateData
func updateDocument(ctx context.Context, documentID string, updateData []firestore.Update, config databaseservice.UpdateDocumentConfig) error {
	return databaseservice.UpdateDocument(ctx, fmt.Sprintf("transferPaymentRequests/%s", documentID), updateData, config)
}
