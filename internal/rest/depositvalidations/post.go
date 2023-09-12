package depositvalidations

import (
	"context"
	"encoding/json"
	"fmt"
	"heypay-cash-in-server/constants"
	"heypay-cash-in-server/internal/services/accountsservice"
	"heypay-cash-in-server/internal/services/databaseservice"
	"heypay-cash-in-server/internal/services/databaseservice/transferpaymentrequestsservice"
	"heypay-cash-in-server/internal/services/databaseservice/usersservice"
	"heypay-cash-in-server/internal/services/firebaseservice"
	"heypay-cash-in-server/internal/services/notificationservice"
	"heypay-cash-in-server/settings"
	"heypay-cash-in-server/utils"
	. "heypay-cash-in-server/utils/types"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
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
	Valid   bool   `json:"valid"`
}

// PostDepositValidations - Validates if the destination account can receive the incoming deposit which has not been processed yet.
func PostDepositValidations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request Request
	utils.ChangeStep(&ctx, "decode-request-body", "01")
	utils.Info(ctx, "Decoding request body")
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error}, "Error decoding request body - reason: %v", err)
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

	utils.ChangeStep(&ctx, "check-banned-origin-ruts", "03")
	*request.OriginRut = utils.RemoveStartingZeros(*request.OriginRut)
	if utils.ArrayIncludes(settings.V.Project.BannedOriginRuts, *request.OriginRut) {
		utils.Warnmf(ctx, Map{"eventName": "origin-rut-banned", "originRut": *request.OriginRut}, "Warning origin rut %s trying to do deposit while being on banned origin rut list", *request.OriginRut)
		utils.HttpRespond(ctx, r, w, http.StatusOK, Response{"origin-rut-banned", "Origin rut banned", false})
		return
	}

	*request.Amount = *request.Amount / 100
	*request.OriginAccountNumber = utils.RemoveStartingZeros(*request.OriginAccountNumber)
	*request.DestinationRut = utils.RemoveStartingZeros(*request.DestinationRut)
	*request.DestinationAccountNumber = utils.RemoveStartingZeros(*request.DestinationAccountNumber)
	destinationNationalID := (*request.DestinationAccountNumber)[3:]
	destinationNationalID += utils.GenerateDigitVerifier(destinationNationalID)

	message := ""
	if request.Message != nil {
		message = *request.Message
	}
	utils.Infomf(ctx, Map{
		"eventName":                "new-deposit-validation",
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
	}, "New deposit validation: {Amount: %v, Message: %v, OriginRut: %v, OriginBankID: %v, OriginAccountNumber: %v, OriginAccountTypeID: %v, DestinationRut: %v, DestinationAccountNumber %v, DestinationAccountTypeID: %v, ReceivedAt: %v, TraceNumber: %v}",
		*request.Amount, message, *request.OriginRut, *request.OriginBankID, *request.OriginAccountNumber, *request.OriginAccountTypeID, *request.DestinationRut, *request.DestinationAccountNumber, *request.DestinationAccountTypeID, *request.ReceivedAt, *request.TraceNumber)

	utils.Infomf(ctx, Map{"eventName": "start-deposit-validation", "type": "payment-request"}, "Starting deposit validation")
	var accountID *string
	var userID *string
	now := time.Now()
	utils.ChangeStep(&ctx, "start-transfer-payment-requests-transaction", "04")
	err = firebaseservice.Db.RunTransaction(ctx, func(txCtx context.Context, transaction *firestore.Transaction) error {
		accountID = nil
		utils.ChangeStep(&ctx, "find-transfer-payment-requests", "05")
		utils.Infof(ctx, "Finding pending transfer payment requests by amount %v, destinationNationalId %v, originNationalId %v, expiresAt > %v and status pending", *request.Amount, *request.DestinationRut, *request.OriginRut, now)
		transferPaymentRequestDocument, err := transferpaymentrequestsservice.FindPendingTransferPaymentRequest(ctx, *request.Amount, *request.DestinationRut, *request.DestinationAccountNumber, *request.OriginRut, now, transaction)
		if err != nil {
			return &utils.HandledError{
				CapturedError: err,
				Code:          "find-payment-transfer-request-failed",
				Message:       fmt.Sprintf("Error finding transfer payment requests - reason: %v", err),
			}
		}
		if transferPaymentRequestDocument != nil {
			utils.ChangeStep(&ctx, "update-transfer-payment-request", "06")
			accountID = &transferPaymentRequestDocument.DestinationAccountID
			err := transferpaymentrequestsservice.UpdateDocument(ctx, transferPaymentRequestDocument.ID, []firestore.Update{
				{Path: "status", Value: "validated"},
				{Path: "transferId", Value: fmt.Sprintf("%s*%s*%s", *request.ReceivedAt, *request.TraceNumber, *request.OriginBankID)},
				{Path: "updatedAt", Value: firestore.ServerTimestamp},
			}, databaseservice.UpdateDocumentConfig{Transaction: transaction})
			if err != nil {
				return &utils.HandledError{
					CapturedError: err,
					Code:          "update-transfer-payment-request-failed",
					Message:       fmt.Sprintf("Error updating transfer payment requests - reason: %v", err),
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
		utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error}, "Error in transfer payment requests transaction - reason: %v", err)
		utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
		return
	}
	if accountID == nil {
		utils.ChangeStep(&ctx, "find-user", "07")
		utils.Infof(ctx, "Getting users by nationalID %v", destinationNationalID)
		usersDocumentList, err := usersservice.GetDocumentsList(ctx, map[string]databaseservice.QueryFieldConfig{
			"nationalId": {WhereFilter: "==", Value: destinationNationalID},
		}, nil)
		if err != nil {
			utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error}, "Error getting users - reason: %v", err)
			utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
			return
		}
		if len(usersDocumentList) == 0 {
			utils.Infof(ctx, "No user found that matches deposit's destination rut")
			utils.HttpRespond(ctx, r, w, http.StatusOK, Response{"user-not-found", "No user found that matches deposit's destination rut", false})
			return
		}
		accountID = &usersDocumentList[0].PrimaryAccountID
		userID = &usersDocumentList[0].ID
	}
	utils.ChangeStep(&ctx, "get-account", "08")
	accountDocument, err := accountsservice.GetAccount(ctx, *accountID)
	if err != nil {
		utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error}, "Error getting account %s on accounts engine - reason: %v", *accountID, err)
		utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
		return
	}
	if accountDocument.Status != "active" {
		utils.Infof(ctx, "Receiver account is not active - status: %v", accountDocument.Status)
		utils.HttpRespond(ctx, r, w, http.StatusOK, Response{fmt.Sprintf("account-%s", accountDocument.Status), "Receiver account is not active", false})
		return
	}
	// null MaxBalances are handled with 0 value
	if accountDocument.MaxBalance != 0 && accountDocument.Balance+accountDocument.BlockedBalance+*request.Amount > accountDocument.MaxBalance {
		utils.Infof(ctx, "Account cannot receive this amount, exceeds its max balance limit")
		if userID != nil {
			utils.ChangeStep(&ctx, "send-over-max-balance-push-notification", "09")
			_ = notificationservice.SendPushNotification(ctx, *userID, notificationservice.Payload{
				Data: map[string]string{
					"type":       "deposit-rejected-max-balance",
					"resourceId": *request.TraceNumber,
				},
				Title: "Esta transferencia supera tu máximo permitido",
				Body:  "Agranda tu billetera ingresando al menú de tu App Chek",
			})
		}
		utils.HttpRespond(ctx, r, w, http.StatusOK, Response{"over-max-balance", "Account cannot receive this amount, exceeds its limit", false})
		return
	}
	if accountDocument.OwnerType == "user" && userID != nil {
		utils.ChangeStep(&ctx, "check-max-daily-bank-transfer-deposit-amount-rule", "10")
		ruleCheckResult, err := firebaseservice.CheckMaxDailyBankTransferDepositAmount(ctx, *request.Amount, *request.DestinationAccountNumber, accountDocument.Category)
		if err != nil {
			utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error}, "Error checking max daily bank transfer deposit amount - reason: %v", err)
			utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
			return
		}
		if !ruleCheckResult.Check {
			utils.Infomf(ctx, Map{"check": ruleCheckResult.Check, "currentValue": ruleCheckResult.CurrentValue, "limitValue": ruleCheckResult.LimitValue}, "Account cannot receive this amount, max daily bank transfer deposit amount exceeded")
			utils.ChangeStep(&ctx, "send-over-max-daily-bank-transfer-deposit-amount-push-notification", "11")
			_ = notificationservice.SendPushNotification(ctx, *userID, notificationservice.Payload{
				Data: map[string]string{
					"type":       "deposit-rejected",
					"resourceId": *request.TraceNumber,
				},
				Title: "¡No pudimos aceptar tu carga de saldo!",
				Body:  fmt.Sprintf("%s es el máximo diario que puedes cargar a tu cuenta %s.", utils.FormatCurrency(ruleCheckResult.LimitValue), settings.V.AppServerProject.AppName),
			})
			utils.HttpRespond(ctx, r, w, http.StatusOK, Response{"over-daily-deposit-amount", "Account cannot receive this amount, exceeds its daily deposit amount limit", false})
			return
		}
		utils.ChangeStep(&ctx, "check-max-monthly-bank-transfer-deposit-amount-rule", "12")
		ruleCheckResult, err = firebaseservice.CheckMaxMonthlyBankTransferDepositAmount(ctx, *request.Amount, *request.DestinationAccountNumber, accountDocument.Category)
		if err != nil {
			utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error}, "Error checking max monthly bank transfer deposit amount - reason: %v", err)
			utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
			return
		}
		if !ruleCheckResult.Check {
			utils.Infomf(ctx, Map{"check": ruleCheckResult.Check, "currentValue": ruleCheckResult.CurrentValue, "limitValue": ruleCheckResult.LimitValue}, "Account cannot receive this amount, max monthly bank transfer deposit amount exceeded")
			utils.ChangeStep(&ctx, "send-over-max-monthly-bank-transfer-deposit-amount-push-notification", "13")
			_ = notificationservice.SendPushNotification(ctx, *userID, notificationservice.Payload{
				Data: map[string]string{
					"type":       "deposit-rejected",
					"resourceId": *request.TraceNumber,
				},
				Title: "¡No pudimos aceptar tu carga de saldo!",
				Body:  fmt.Sprintf("%s es el máximo mensual que puedes cargar a tu cuenta %s", utils.FormatCurrency(ruleCheckResult.LimitValue), settings.V.AppServerProject.AppName),
			})
			utils.HttpRespond(ctx, r, w, http.StatusOK, Response{"over-monthly-deposit-amount", "Account cannot receive this amount, exceeds its monthly deposit amount limit", false})
			return
		}
	} else if accountDocument.OwnerType == "commerce" {
		utils.ChangeStep(&ctx, "check-max-daily-bank-transfer-payment-amount-rule", "14")
		ruleCheckResult, err := firebaseservice.CheckMaxDailyBankTransferPaymentAmount(ctx, *request.Amount, *request.DestinationAccountNumber, accountDocument.Category)
		if err != nil {
			utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error}, "Error checking max daily bank transfer payment amount - reason: %v", err)
			utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
			return
		}
		if !ruleCheckResult.Check {
			utils.Infomf(ctx, Map{"check": ruleCheckResult.Check, "currentValue": ruleCheckResult.CurrentValue, "limitValue": ruleCheckResult.LimitValue}, "Account cannot receive this amount, max daily bank transfer payment amount exceeded")
			utils.HttpRespond(ctx, r, w, http.StatusOK, Response{"over-daily-deposit-payment-amount", "Account cannot receive this amount, exceeds its daily payment amount limit", false})
			return
		}
		utils.ChangeStep(&ctx, "check-max-monthly-bank-transfer-deposit-amount-rule", "15")
		ruleCheckResult, err = firebaseservice.CheckMaxMonthlyBankTransferPaymentAmount(ctx, *request.Amount, *request.DestinationAccountNumber, accountDocument.Category)
		if err != nil {
			utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error}, "Error checking max monthly bank transfer payment amount - reason: %v", err)
			utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
			return
		}
		if !ruleCheckResult.Check {
			utils.Infomf(ctx, Map{"check": ruleCheckResult.Check, "currentValue": ruleCheckResult.CurrentValue, "limitValue": ruleCheckResult.LimitValue}, "Account cannot receive this amount, max monthly bank transfer payment amount exceeded")
			utils.HttpRespond(ctx, r, w, http.StatusOK, Response{"over-monthly-deposit-payment-amount", "Account cannot receive this amount, exceeds its monthly payment amount limit", false})
			return
		}
	}
	utils.Infomf(ctx, Map{"eventName": "deposit-validated-successfully"}, "Deposit was validated and can be transferred")
	utils.HttpRespond(ctx, r, w, http.StatusOK, Response{"valid", "Deposit is valid", true})
}
