package identitytoolkitservice

import (
	"context"
	"heypay-cash-in-server/settings"
	"heypay-cash-in-server/utils"

	"google.golang.org/api/identitytoolkit/v3"
	"google.golang.org/api/option"
)

var RelyingPartyService *identitytoolkit.RelyingpartyService

func init() {
	utils.Info(nil, "Initializing identity toolkit service...")
	ctx := context.Background()
	identityToolkitService, err := identitytoolkit.NewService(ctx, option.WithAPIKey(settings.V.AppServerProject.APIKey))
	if err != nil {
		utils.Fatalf(nil, "Fatal error while initializing identity toolkit service - reason: %v", err)
	}
	RelyingPartyService = identitytoolkit.NewRelyingpartyService(identityToolkitService)
}

func VerifyCustomToken(token string) (*identitytoolkit.VerifyCustomTokenResponse, error) {
	verifyCustomTokenRequest := identitytoolkit.IdentitytoolkitRelyingpartyVerifyCustomTokenRequest{Token: token, ReturnSecureToken: true}
	call := RelyingPartyService.VerifyCustomToken(&verifyCustomTokenRequest)
	return call.Do()
}
