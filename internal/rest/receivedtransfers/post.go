package receivedtransfers

import (
	"context"
	"encoding/json"
	"fmt"
	"heypay-cash-in-server/constants"
	"heypay-cash-in-server/internal/services/databaseservice/transferpaymentrequestsservice"
	"heypay-cash-in-server/internal/services/firebaseservice"
	"heypay-cash-in-server/utils"
	. "heypay-cash-in-server/utils/types"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Request struct {
	Amount                   *int64  `json:"amount"`
	DestinationAccountNumber *string `json:"destinationAccountNumber"`
	DestinationAccountTypeID *string `json:"destinationAccountTypeId"`
	DestinationRut           *string `json:"destinationRut"`
	Message                  *string `json:"message"`
	OriginAccountNumber      *string `json:"originAccountNumber"`
	OriginAccountTypeID      *string `json:"originAccountTypeId"`
	OriginBankID             *string `json:"originBankId"`
	OriginRut                *string `json:"originRut"`
	ReceivedAt               *string `json:"receivedAt"`
	TraceNumber              *string `json:"traceNumber"`
}
type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

// PostReceivedTransfers - Creates a transfer to be processed.
func PostReceivedTransfers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	utils.ChangeStep(&ctx, "decode-auth-token", "14")
	utils.Info(ctx, "Decoding request auth token")
	decodedToken, err := firebaseservice.DecodeToken(ctx, r.Header)
	if err != nil {
		if strings.HasPrefix(err.Error(), "auth") {
			utils.Errormf(ctx, Map{"eventName": "access-token-expired"}, "Access token expired")
			utils.HttpRespondError(ctx, r, w, http.StatusUnauthorized, "Access token expired.", "access-token-expired")
			return
		}
		utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error, "code": "auth-token-decode-failed"}, "Error could not decode auth token - reason: %v", err)
		utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
		return
	}

	utils.ChangeStep(&ctx, "verify-authentication", "15")
	firebaseservice.PrintDecodedToken(ctx, *decodedToken)
	if decodedToken.Claims["admin"] == nil && decodedToken.Claims["cash_in"] == nil {
		utils.Warnmf(ctx, Map{"eventName": "unauthorized-attempt-to-create-received-transfer"}, "Warning decodedToken claims doesn't meet validations, unauthorized attempt to create received transfer: %+v", decodedToken)
		utils.HttpRespondError(ctx, r, w, http.StatusUnauthorized, "Unauthorized.", "unauthorized")
		return
	}

	var request Request
	utils.ChangeStep(&ctx, "decode-request-body", "01")
	utils.Info(ctx, "Decoding request body")
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error, "code": "request-body-decode-failed"}, "Error decoding request body - reason: %v", err)
		utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
		return
	}

	utils.ChangeStep(&ctx, "check-request-parameters", "02")
	utils.Info(ctx, "Checking request parameters")
	if request.Amount == nil ||
		request.DestinationAccountNumber == nil ||
		request.DestinationAccountTypeID == nil ||
		request.DestinationRut == nil ||
		request.OriginAccountNumber == nil ||
		request.OriginAccountTypeID == nil ||
		request.OriginBankID == nil ||
		request.OriginRut == nil ||
		request.ReceivedAt == nil ||
		request.TraceNumber == nil {
		utils.Errormf(ctx, Map{"eventName": "invalid-parameters"}, "Error request parameters are not valid")
		utils.HttpRespondError(ctx, r, w, http.StatusBadRequest, "Invalid request parameters", "invalid-parameters")
		return
	}

	*request.Amount = *request.Amount / 100
	*request.OriginRut = utils.RemoveStartingZeros(*request.OriginRut)
	*request.OriginAccountNumber = utils.RemoveStartingZeros(*request.OriginAccountNumber)
	if value, exists := constants.BankAccountTypeID[*request.OriginAccountTypeID]; exists {
		*request.OriginAccountTypeID = value
	}
	originBankID := *request.OriginBankID
	if value, exists := constants.BankID[*request.OriginBankID]; exists {
		*request.OriginBankID = value
	}
	*request.DestinationRut = utils.RemoveStartingZeros(*request.DestinationRut)
	*request.DestinationAccountNumber = utils.RemoveStartingZeros(*request.DestinationAccountNumber)
	if value, exists := constants.BankAccountTypeID[*request.DestinationAccountTypeID]; exists {
		*request.DestinationAccountTypeID = value
	}
	accountRut := (*request.DestinationAccountNumber)[3:]
	accountRut += utils.GenerateDigitVerifier(accountRut)

	message := ""
	if request.Message != nil {
		message = *request.Message
	}

	utils.Infomf(ctx, Map{
		"eventName":                "new-received-transfer",
		"amount":                   *request.Amount,
		"message":                  message,
		"originRut":                *request.OriginRut,
		"originBankId":             *request.OriginBankID,
		"originAccountNumber":      *request.OriginAccountNumber,
		"originAccountTypeId":      *request.OriginAccountTypeID,
		"destinationRut":           *request.DestinationRut,
		"destinationAccountNumber": *request.DestinationAccountNumber,
		"destinationAccountTypeId": *request.DestinationAccountTypeID,
		"receivedAt":               *request.ReceivedAt,
		"traceNumber":              *request.TraceNumber,
	}, "New received transfer: {Amount: %v, Message: %v, OriginRut: %v, OriginBankID: %v, OriginAccountNumber: %v, OriginAccountTypeID: %v, DestinationRut: %v, DestinationAccountNumber: %v, DestionationAccountTypeId: %v, ReceivedAt: %v, TraceNumber: %v}",
		*request.Amount, message, *request.OriginRut, *request.OriginBankID, *request.OriginAccountNumber, *request.OriginAccountTypeID, *request.DestinationRut, *request.DestinationAccountNumber, *request.DestinationAccountTypeID, *request.ReceivedAt, *request.TraceNumber)

	transferID := fmt.Sprintf("%s*%s*%s", *request.ReceivedAt, *request.TraceNumber, originBankID)
	utils.Infof(ctx, "transferID: %v", transferID)
	utils.ChangeStep(&ctx, "start-transaction", "16")
	err = firebaseservice.Db.RunTransaction(ctx, func(txCtx context.Context, transaction *firestore.Transaction) error {
		utils.ChangeStep(&ctx, "find-transfer", "17")
		transferRef := firebaseservice.Db.Doc(fmt.Sprintf("transfers/%s", transferID))
		transferSnapshot, err := transaction.Get(transferRef)
		if err != nil {
			if status.Code(err) != codes.NotFound {
				return &utils.HandledError{
					CapturedError: err,
					Code:          "get-transfer-failed",
					Message:       fmt.Sprintf("Error getting transfer - reason: %v", err),
				}
			}
		}
		if !transferSnapshot.Exists() {
			utils.ChangeStep(&ctx, "find-transfer-payment-requests", "18")
			var transferPaymentRequestID string
			transferPaymentRequestDocument, err := transferpaymentrequestsservice.FindValidatedTransferPaymentRequest(ctx, transferID, transaction)
			if err != nil {
				return &utils.HandledError{
					CapturedError: err,
					Code:          "find-transfer-payment-request-failed",
					Message:       fmt.Sprintf("Error finding transfer payment requests - reason: %v", err),
				}
			}
			if transferPaymentRequestDocument != nil {
				transferPaymentRequestID = transferPaymentRequestDocument.ID
				err := transaction.Update(firebaseservice.Db.Doc(fmt.Sprintf("transferPaymentRequests/%s", transferPaymentRequestDocument.ID)),
					[]firestore.Update{{Path: "status", Value: "done"},
						{Path: "updatedAt", Value: firestore.ServerTimestamp}})
				if err != nil {
					return &utils.HandledError{
						CapturedError: err,
						Code:          "update-transfer-payment-request-failed",
						Message:       fmt.Sprintf("Error updating transfer payment requests - reason: %v", err),
					}
				}
			}
			utils.ChangeStep(&ctx, "writes-in-transaction", "19")
			err = transaction.Create(transferRef, map[string]interface{}{
				"additionalData": map[string]interface{}{
					"receivedAt":  *request.ReceivedAt,
					"traceNumber": *request.TraceNumber,
				},
				"amount":    *request.Amount,
				"createdAt": firestore.ServerTimestamp,
				"destinationAccountInfo": map[string]interface{}{
					"accountNumber": *request.DestinationAccountNumber,
					"accountTypeId": *request.DestinationAccountTypeID,
					"bankId":        "ripley-chile",
					"nationalId":    *request.DestinationRut,
				},
				"message": message,
				"originAccountInfo": map[string]interface{}{
					"accountNumber": *request.OriginAccountNumber,
					"accountTypeId": *request.OriginAccountTypeID,
					"bankId":        *request.OriginBankID,
					"nationalId":    *request.OriginRut,
				},
				"transferPaymentRequestId": transferPaymentRequestID,
				"updatedAt":                firestore.ServerTimestamp,
			})
			if err != nil {
				return &utils.HandledError{
					CapturedError: err,
					Code:          "create-transfer-failed",
					Message:       fmt.Sprintf("Error could not create transfer - reason: %v", err),
				}
			}
		}
		return nil
	})
	if err != nil {
		if handledError, ok := err.(*utils.HandledError); ok {
			utils.Errormf(ctx, Map{"eventName": "known-error", "code": handledError.Code}, handledError.Message, handledError.CapturedError)
			utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
			return
		}
		utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error}, "Error of transaction - reason: %v", err)
		utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
		return
	}
	utils.Infomf(ctx, Map{"eventName": "transfer-received-successfully", "transferId": transferID}, "Transfer was processed successfully: %v", transferID)
	utils.HttpRespond(ctx, r, w, http.StatusOK, Response{Code: "ok", ID: transferID, Message: "Transfer was processed successfully."})
}
