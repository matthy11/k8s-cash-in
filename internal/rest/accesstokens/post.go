package accesstokens

import (
	"encoding/json"
	"heypay-cash-in-server/constants"
	"heypay-cash-in-server/internal/services/firebaseservice"
	"heypay-cash-in-server/internal/services/identitytoolkitservice"
	"heypay-cash-in-server/utils"
	. "heypay-cash-in-server/utils/types"
	"net/http"
)

type Request struct {
	ClientID     *string `json:"clientId"`
	ClientSecret *string `json:"clientSecret"`
	PrivateKey   *string `json:"privateKey"`
	Token        *string `json:"token"`
	Scope        *string `json:"scope"`
}
type Response struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

// PostAccessTokens - Generates an access token needed to consume the other endpoints.
func PostAccessTokens(w http.ResponseWriter, r *http.Request) {
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
	if (request.ClientID == nil && request.Token == nil) ||
		(request.ClientID != nil && request.ClientSecret == nil && request.PrivateKey == nil) {
		utils.Errormf(ctx, Map{"eventName": "invalid-parameters"}, "Error request parameters are not valid")
		utils.HttpRespondError(ctx, r, w, http.StatusBadRequest, "Invalid request parameters", "invalid-parameters")
		return
	}

	utils.Infomf(ctx, Map{
		"eventName":    "new-access-token-generation",
		"clientId":     utils.MaskStringPointer(request.ClientID),
		"clientSecret": utils.MaskStringPointer(request.ClientSecret),
		"privateKey":   utils.MaskStringPointer(request.PrivateKey),
	}, "New access token generation: {ClientId: %v, ClientSecret: %v, PrivateKey: %v}", utils.MaskStringPointer(request.ClientID), utils.MaskStringPointer(request.ClientSecret), utils.MaskStringPointer(request.PrivateKey))

	var token string
	if request.ClientID != nil {
		utils.ChangeStep(&ctx, "get-client-claims", "03")
		privateKey := *request.ClientSecret
		if request.PrivateKey != nil {
			privateKey = *request.PrivateKey
		}
		if request.Scope == nil {
			scope := "admin"
			request.Scope = &scope
		}
		claims, err := firebaseservice.GetClientClaims(ctx, *request.ClientID, privateKey, *request.Scope)
		if err != nil || len(claims) == 0 {
			utils.Warnmf(ctx, Map{"eventName": "invalid-credentials"}, "Warning invalid credentials for client %s", *request.ClientID)
		}
		if err != nil {
			utils.Errorf(ctx, "Error getting client claims - reason: %v", err)
			switch err.Error() {
			case "client-not-found":
				utils.HttpRespondError(ctx, r, w, http.StatusUnauthorized, "Client not found", err.Error())
				return
			case "private-key-not-active":
				utils.HttpRespondError(ctx, r, w, http.StatusUnauthorized, "Private key not active", err.Error())
				return
			case "blocked":
				utils.HttpRespondError(ctx, r, w, http.StatusUnauthorized, "Secret key has been blocked", err.Error())
				return
			}
			utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
			return
		}
		utils.Infof(ctx, "clientClaims: %+v", claims)
		if len(claims) == 0 {
			utils.Error(ctx, "Error no secret found matching credentials")
			utils.HttpRespondError(ctx, r, w, http.StatusUnauthorized, "No secret found matching credentials", "invalid-credentials")
			return
		}
		utils.ChangeStep(&ctx, "create-token", "04")
		token, err = firebaseservice.Auth.CustomTokenWithClaims(ctx, *request.ClientID, claims)
		if err != nil {
			utils.Errormf(ctx, Map{"eventName": constants.ErrorCode_Unknown_Error}, "Error could not create custom token with claims - reason: %v", err)
			utils.HttpRespondError(ctx, r, w, http.StatusInternalServerError, constants.ErrorLabel_Error_In_Server, utils.GetCurrentStepCode(ctx))
			return
		}
	} else {
		token = *request.Token
	}

	utils.ChangeStep(&ctx, "get-id-token", "05")
	utils.Info(ctx, "Verifying custom token")
	verifyCustomTokenResponse, err := identitytoolkitservice.VerifyCustomToken(token)
	if err != nil {
		utils.Errorf(ctx, "Error could not verify token - reason: %v", err)
		utils.HttpRespondError(ctx, r, w, http.StatusBadRequest, "Could not verify token", "invalid-token")
		return
	}

	utils.Infomf(ctx, Map{"eventName": "access-token-generated-successfully"}, "Custom token was verified and access token generated")
	utils.HttpRespond(ctx, r, w, http.StatusOK, Response{verifyCustomTokenResponse.IdToken, verifyCustomTokenResponse.RefreshToken, verifyCustomTokenResponse.ExpiresIn * 1000})
}
