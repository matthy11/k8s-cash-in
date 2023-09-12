package transferpaymentrequest

import (
	"time"
)

type TransferPaymentRequest struct {
	AdditionalData        map[string]interface{} `json:"additionalData"`
	Amount                int64                  `firestore:"amount"`
	CreatedAt             time.Time              `json:"createdAt"`
	Currency              string                 `json:"currency"`
	DestinationAccountID  string                 `json:"destinationAccountId"`
	DestinationNationalID string                 `json:"destinationNationalId"`
	ExpiresAt             time.Time              `json:"expiresAt"`
	ExpiresIn             int64                  `firestore:"expiresIn"`
	ID                    string                 // added on runtime
	OriginNationalID      string                 `json:"originNationalId"`
	ReceivedTransferID    string                 `json:"receivedTransferId"`
	TransactionID         string                 `json:"transactionId"`
	UpdatedAt             time.Time              `json:"updatedAt"`
}
