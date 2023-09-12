package accountsengineservice

import (
	"context"
	"fmt"
	"heypay-cash-in-server/settings"
	"heypay-cash-in-server/utils"
	"heypay-cash-in-server/utils/httpclient"
	. "heypay-cash-in-server/utils/types"
	"time"
)

// RefreshWindow - AccessToken lasts 60m, so we subtract 10m (600s) of it in order to maintain a safe refresh window and refresh it after 50m instead
const RefreshWindow int64 = 600

type Token struct {
	AccessToken  string `json:"accessToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	RefreshToken string `json:"refreshToken"`
}

var AccessToken Token
var nextTokenRefresh int64

func init() {
	utils.Info(nil, "Initializing accounts engine service...")
	err := httpclient.Post(nil, fmt.Sprintf("%s/accessTokens", settings.V.AccountsEngine.APIURL), Map{
		"clientId":   settings.V.AccountsEngine.ClientID,
		"privateKey": settings.V.AccountsEngine.PrivateKey,
		"scope":      "admin",
	}, httpclient.HttpConfig{
		MaskRequestAttributes:  []string{"privateKey"},
		MaskResponseAttributes: []string{"accessToken", "refreshToken"},
	}, &AccessToken)
	if err != nil {
		utils.Fatalf(nil, "Fatal error initializing accounts engine service - reason: %v", err)
	}
	nextTokenRefresh = time.Now().Unix() + (AccessToken.ExpiresIn / 1000) - RefreshWindow
}

func DidTokenExpire() bool {
	return time.Now().Unix() > nextTokenRefresh
}

func RefreshAccessToken(ctx context.Context) error {
	utils.Debugf(ctx, "Refreshing accounts engine's access token.")
	err := httpclient.Post(ctx, fmt.Sprintf("%s/refreshTokens", settings.V.AccountsEngine.APIURL), Map{
		"refreshToken": AccessToken.RefreshToken,
	}, httpclient.HttpConfig{
		MaskRequestAttributes:  []string{"refreshToken"},
		MaskResponseAttributes: []string{"accessToken", "refreshToken"},
	}, &AccessToken)
	if err != nil {
		utils.Errormf(ctx, Map{"eventName": "accounts-engine-error", "error": fmt.Sprintf("%v", err)}, "Error refreshing accounts engine's access token - reason: %s", err)
		return err
	}
	nextTokenRefresh = time.Now().Unix() + (AccessToken.ExpiresIn / 1000) - RefreshWindow
	return nil
}
