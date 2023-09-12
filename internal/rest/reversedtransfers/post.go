package reversedtransfers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"heypay-cash-in-server/constants"
	"heypay-cash-in-server/internal/services/firebaseservice"
	"heypay-cash-in-server/models/transfer"
	"heypay-cash-in-server/utils"
	. "heypay-cash-in-server/utils/types"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
)

type Request struct {
	OriginBankID *string `json:"originBankId"`
	ReceivedAt   *string `json:"receivedAt"`
	TraceNumber  *string `json:"traceNumber"`
}
type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

// PostReversedTransfers - Begins a reverse procedure of a received transfer.
func PostReversedTransfers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	utils.ChangeStep(&ctx, "decode-auth-token", "14")
	utils.Info(ctx, "Decoding request auth token")
	decodedToken, err := firebaseservice.DecodeToken(ctx, r.Header)
	if err != nil {
		if strings.HasPrefix(err.Error(), "auth") {
			utils.Errormf(ctx, Map{"eventName": "access-token-expired"}, "Access token expired")
			utils.HttpRespondError(ctx, r, w, http.StatusUnauthorized, "Access token expired", "access-token-expired")
			return
		}
		utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error}, "Error could not decode auth token - reason: %v", err)
		utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
		return
	}

	utils.ChangeStep(&ctx, "verify-authentication", "15")
	firebaseservice.PrintDecodedToken(ctx, *decodedToken)
	if decodedToken.Claims["admin"] == nil && decodedToken.Claims["cash_in"] == nil {
		utils.Warnmf(ctx, Map{"eventName": "unauthorized-attempt-to-create-reversed-transfer"}, "Warning decodedToken claims doesn't meet validations, unauthorized attempt to create reversed transfer: %+v", decodedToken)
		utils.HttpRespondError(ctx, r, w, http.StatusUnauthorized, "Unauthorized", "unauthorized")
		return
	}

	var request Request
	utils.ChangeStep(&ctx, "decode-request-body", "01")
	utils.Info(ctx, "Decoding request body")
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error}, "Error decoding request body - reason: %v", err)
		utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
		return
	}

	utils.ChangeStep(&ctx, "check-request-parameters", "02")
	utils.Info(ctx, "Checking request parameters")
	if request.OriginBankID == nil ||
		request.ReceivedAt == nil ||
		request.TraceNumber == nil {
		utils.Errormf(ctx, Map{"eventName": "invalid-parameters"}, "Error request parameters are not valid")
		utils.HttpRespondError(ctx, r, w, http.StatusBadRequest, "Invalid request parameters", "invalid-parameters")
		return
	}

	originBankID := *request.OriginBankID
	if value, exists := constants.BankID[*request.OriginBankID]; exists {
		*request.OriginBankID = value
	}

	utils.Infomf(ctx, Map{
		"eventName":    "new-reversed-transfer",
		"originBankId": *request.OriginBankID,
		"receivedAt":   *request.ReceivedAt,
		"traceNumber":  *request.TraceNumber,
	}, "New reversed transfer {OriginBankID: %s, ReceivedAt: %s, TraceNumber: %s}", *request.OriginBankID, *request.ReceivedAt, *request.TraceNumber)

	transferID := fmt.Sprintf("%s*%s*%s", *request.ReceivedAt, *request.TraceNumber, originBankID)
	utils.Infof(ctx, "transferID: %v", transferID)
	utils.ChangeStep(&ctx, "start-transaction", "16")
	transferRef := firebaseservice.Db.Doc(fmt.Sprintf("transfers/%s", transferID))
	_ = firebaseservice.Db.RunTransaction(ctx, func(txCtx context.Context, transaction *firestore.Transaction) error {
		utils.ChangeStep(&ctx, "reads-in-transaction", "17")
		transferSnapshot, err := transaction.Get(transferRef)
		if transferSnapshot != nil {
			if !transferSnapshot.Exists() {
				utils.Error(ctx, "Transfer not found")
				utils.HttpRespondError(ctx, r, w, http.StatusNotFound, "Document not found", "not-found")
				return errors.New("not-found")
			}
			var transferData transfer.Transfer
			err = transferSnapshot.DataTo(&transferData)
			if err != nil {
				utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error}, "Error getting transfer data - reason: %v", err)
				utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
				return err
			}
			if transferData.Reversed {
				utils.Warn(ctx, "Transfer was already reversed")
				utils.HttpRespond(ctx, r, w, http.StatusOK, Response{"Transfer was already reversed", "already-processed", transferID})
				return errors.New("already-processed")
			}
			utils.ChangeStep(&ctx, "writes-in-transaction", "18")
			err = transaction.Update(transferRef, []firestore.Update{
				{Path: "reversed", Value: true},
				{Path: "reverseStatus", Value: "pending"},
				{Path: "updatedAt", Value: firestore.ServerTimestamp}})
			if err != nil {
				utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error}, "Error updating transfer - reason: %v", err)
				utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
				return err
			}
			utils.Infomf(ctx, Map{"eventName": "transfer-reversed-successfully", "transferId": transferID}, "Transfer successfully reversed")
			utils.HttpRespond(ctx, r, w, http.StatusOK, Response{"Transfer successfully reversed", "success", transferID})
			return nil
		}
		utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error}, "Error getting transfer - reason: %v", err)
		utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
		return err
	})
}
