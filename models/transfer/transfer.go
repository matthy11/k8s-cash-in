package transfer

type AccountInfo struct {
	AccountNumber string `firestore:"accountNumber"`
	AccountTypeID string `firestore:"bankAccountTypeId"`
	BankID        string `firestore:"bankId"`
	NationalID    string `firestore:"nationalId"`
}

type Transfer struct {
	AdditionalData                 map[string]interface{} `firestore:"additionalData"`
	Amount                         int64                  `firestore:"amount"`
	ChargeTransferPaymentRequestID string                 `firestore:"chargeTransferPaymentRequestId"`
	CreatedAt                      interface{}            `firestore:"createdAt"`
	DepositID                      string                 `firestore:"depositId"`
	DestinationAccountInfo         AccountInfo            `firestore:"destinationAccountInfo"`
	DestinationID                  string                 `firestore:"destinationId"`
	DestinationType                string                 `firestore:"destinationType"`
	Message                        string                 `firestore:"message"`
	OriginAccountInfo              AccountInfo            `firestore:"originAccountInfo"`
	Reversed                       bool                   `firestore:"reversed"`
	ReversedAt                     interface{}            `firestore:"reversedAt"`
	ReverseStatus                  string                 `firestore:"reverseStatus"`
	ReverseStatusData              map[string]interface{} `firestore:"reverseStatusData"`
	ReverseEventID                 string                 `firestore:"reverseEventId"`
	Status                         string                 `firestore:"status"`
	UnassignedReason               string                 `firestore:"unassignedReason"`
	UpdatedAt                      interface{}            `firestore:"updatedAt"`
}
