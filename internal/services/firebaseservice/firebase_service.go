package firebaseservice

import (
	"context"
	"errors"
	"heypay-cash-in-server/models/oauthclient"
	"heypay-cash-in-server/settings"
	"heypay-cash-in-server/utils"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/messaging"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/option"
)

var (
	Db   *firestore.Client
	Auth *auth.Client
	Msg  *messaging.Client

	CheckMaxDailyBankTransferDepositAmount   func(ctx context.Context, depositAmount int64, accountNumber, accountCategory string) (RuleCheckResult, error)
	CheckMaxMonthlyBankTransferDepositAmount func(ctx context.Context, depositAmount int64, accountNumber, accountCategory string) (RuleCheckResult, error)
	CheckMaxDailyBankTransferPaymentAmount   func(ctx context.Context, depositAmount int64, accountNumber, accountCategory string) (RuleCheckResult, error)
	CheckMaxMonthlyBankTransferPaymentAmount func(ctx context.Context, depositAmount int64, accountNumber, accountCategory string) (RuleCheckResult, error)
)

func init() {
	utils.Info(nil, "Initializing firebase service...")
	ctx := context.Background()

	var app *firebase.App
	var err error
	if settings.V.Env == "LOCAL" {
		clientOption := option.WithCredentialsFile("k8s/service-account.json")
		app, err = firebase.NewApp(ctx, nil, clientOption)
	} else {
		app, err = firebase.NewApp(ctx, nil)
	}
	if err != nil {
		utils.Fatalf(nil, "Fatal error initializing firebase service - reason: %v", err)
	}
	Db, err = app.Firestore(ctx)
	if err != nil {
		utils.Fatalf(nil, "Fatal error initializing firebase firestore service - reason: %v", err)
	}
	Auth, err = app.Auth(ctx)
	if err != nil {
		utils.Fatalf(nil, "Fatal error initializing firebase auth service - reason: %v", err)
	}
	Msg, err = app.Messaging(ctx)
	if err != nil {
		utils.Fatalf(nil, "Fatal error initializing firebase messaging service - reason: %v", err)
	}

	CheckMaxDailyBankTransferDepositAmount = checkMaxDailyBankTransferDepositAmount
	CheckMaxMonthlyBankTransferDepositAmount = checkMaxMonthlyBankTransferDepositAmount
	CheckMaxDailyBankTransferPaymentAmount = checkMaxDailyBankTransferPaymentAmount
	CheckMaxMonthlyBankTransferPaymentAmount = checkMaxMonthlyBankTransferPaymentAmount
}

func GetClientClaims(ctx context.Context, clientId, privateKey, scope string) (map[string]interface{}, error) {
	clientRef := Db.Collection("OAuthClients").Doc(clientId)
	clientSnapshot, err := clientRef.Get(ctx)
	if clientSnapshot != nil && !clientSnapshot.Exists() {
		utils.Error(ctx, "Error client not found")
		return nil, errors.New("client-not-found")
	}
	if err != nil {
		utils.Errorf(ctx, "Error getting client - reason: %v", err)
		return nil, err
	}
	clientPrivateKeysSnapshots, err := clientRef.Collection("privateKeys").Documents(ctx).GetAll()
	if err != nil {
		utils.Errorf(ctx, "Error getting client private keys - reason: %v", err)
		return nil, err
	}
	if len(clientPrivateKeysSnapshots) == 0 && err == nil {
		return nil, nil
	}
	for _, doc := range clientPrivateKeysSnapshots {
		var privateKeyData oauthclient.PrivateKey
		err = doc.DataTo(&privateKeyData)
		if err != nil {
			utils.Errorf(ctx, "Error getting private key data - reason: %v", err)
			return nil, err
		}
		err = bcrypt.CompareHashAndPassword([]byte(privateKeyData.Hash), []byte(privateKey))
		if err == nil {
			if privateKeyData.Status != "active" {
				return nil, errors.New("private-key-not-active")
			}
			_, err = doc.Ref.Update(ctx, []firestore.Update{
				{Path: "lastSignIn", Value: firestore.ServerTimestamp},
				{Path: "updatedAt", Value: firestore.ServerTimestamp}})
			if err != nil {
				utils.Errorf(ctx, "Error updating private key - reason: %v", err)
				return nil, err
			}
			_, err = clientRef.Update(ctx, []firestore.Update{
				{Path: "failedAttempts", Value: 0},
				{Path: "lastFailedAttempt", Value: nil},
				{Path: "lastSignIn", Value: firestore.ServerTimestamp},
				{Path: "updateAt", Value: firestore.ServerTimestamp}})
			if err != nil {
				utils.Errorf(ctx, "Error updating client - reason: %v", err)
				return nil, err
			}
			if privateKeyData.Claims[scope] == nil {
				return nil, nil
			}
			return privateKeyData.Claims[scope], nil
		}
	}
	var clientData oauthclient.OAuthClient
	err = clientSnapshot.DataTo(&clientData)
	if err != nil {
		utils.Errorf(ctx, "Error getting client data - reason: %v", err)
		return nil, err
	}
	if clientData.FailedAttempts >= 2 {
		_, err = clientRef.Update(ctx, []firestore.Update{
			{Path: "status", Value: "blocked"},
			{Path: "blockedReason", Value: "too-many-attempts"},
			{Path: "blockedAt", Value: firestore.ServerTimestamp},
			{Path: "failedAttempts", Value: clientData.FailedAttempts + 1},
			{Path: "lastFailedAttempt", Value: firestore.ServerTimestamp},
			{Path: "updateAt", Value: firestore.ServerTimestamp}})
		if err != nil {
			utils.Errorf(ctx, "Error updating client - reason: %v", err)
			return nil, err
		}
		return nil, errors.New("blocked")
	}
	_, err = clientRef.Update(ctx, []firestore.Update{
		{Path: "failedAttempts", Value: clientData.FailedAttempts + 1},
		{Path: "lastFailedAttempt", Value: firestore.ServerTimestamp},
		{Path: "updateAt", Value: firestore.ServerTimestamp}})
	if err != nil {
		utils.Errorf(ctx, "Error updating client - reason: %v", err)
		return nil, err
	}
	return nil, nil
}

type RuleCheckResult struct {
	Check        bool
	LimitValue   interface{}
	CurrentValue interface{}
}

func checkTransferTotalAmountInTimePeriod(ctx context.Context, days int, depositAmount int64, accountNumber, accountCategory string, maxAmount int64) (RuleCheckResult, error) {
	date := time.Now()
	date = date.AddDate(0, 0, -days)
	transfersSnapshots, err := Db.Collection("transfers").
		Where("destinationAccountInfo.accountNumber", "==", accountNumber).
		Where("createdAt", ">=", date).
		Where("reversed", "==", false).Documents(ctx).GetAll()
	if err != nil {
		utils.Errorf(ctx, "Error getting deposits - reason: %v", err)
		return RuleCheckResult{Check: false}, err
	}
	var amount int64
	for _, doc := range transfersSnapshots {
		amount += doc.Data()["amount"].(int64)
	}
	return RuleCheckResult{
		Check:        amount+depositAmount <= maxAmount,
		LimitValue:   maxAmount,
		CurrentValue: amount,
	}, nil
}

func checkMaxDailyBankTransferDepositAmount(ctx context.Context, depositAmount int64, accountNumber, accountCategory string) (RuleCheckResult, error) {
	maxAmount := settings.V.UserAccountCategory[accountCategory].MaxDailyBankTransferDepositAmount
	return checkTransferTotalAmountInTimePeriod(ctx, 1, depositAmount, accountNumber, accountCategory, maxAmount)
}

func checkMaxMonthlyBankTransferDepositAmount(ctx context.Context, depositAmount int64, accountNumber, accountCategory string) (RuleCheckResult, error) {
	maxAmount := settings.V.UserAccountCategory[accountCategory].MaxMonthlyBankTransferDepositAmount
	return checkTransferTotalAmountInTimePeriod(ctx, 30, depositAmount, accountNumber, accountCategory, maxAmount)
}

func checkMaxDailyBankTransferPaymentAmount(ctx context.Context, depositAmount int64, accountNumber, accountCategory string) (RuleCheckResult, error) {
	maxAmount := settings.V.CommerceAccountCategory[accountCategory].MaxDailyBankTransferPaymentAmount
	return checkTransferTotalAmountInTimePeriod(ctx, 1, depositAmount, accountNumber, accountCategory, maxAmount)
}

func checkMaxMonthlyBankTransferPaymentAmount(ctx context.Context, depositAmount int64, accountNumber, accountCategory string) (RuleCheckResult, error) {
	maxAmount := settings.V.CommerceAccountCategory[accountCategory].MaxMonthlyBankTransferPaymentAmount
	return checkTransferTotalAmountInTimePeriod(ctx, 30, depositAmount, accountNumber, accountCategory, maxAmount)
}

func DecodeToken(ctx context.Context, header http.Header) (*auth.Token, error) {
	authorization := header.Get("Authorization")
	if strings.HasPrefix(authorization, "Bearer ") {
		authToken := authorization[len("Bearer "):]
		return Auth.VerifyIDToken(ctx, authToken)
	}
	return nil, errors.New("no bearer authorization found")
}

func PrintDecodedToken(ctx context.Context, token auth.Token) {
	utils.Infof(ctx, "decodedToken {Audience: %v, Claims: %v, Expires: %v, IssuedAt: %v, Issuer: %v, Subject: %v, UID: %v}", token.Audience, token.Claims, token.Expires, token.IssuedAt, token.Issuer, token.Subject, token.UID)
}
