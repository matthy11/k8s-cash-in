package http

import (
	"context"
	"fmt"
	"heypay-cash-in-server/internal/services/accountsengineservice"
	"heypay-cash-in-server/models/account"
	"heypay-cash-in-server/settings"
	"heypay-cash-in-server/utils/httpclient"
	. "heypay-cash-in-server/utils/types"
)

var (
	GetAccount func(ctx context.Context, id string) (*account.Account, error)
)

func init() {
	GetAccount = getAccount
}

func getAccount(ctx context.Context, id string) (*account.Account, error) {
	if accountsengineservice.DidTokenExpire() {
		err := accountsengineservice.RefreshAccessToken(ctx)
		if err != nil {
			return nil, err
		}
	}
	var accountDocument account.Account
	err := httpclient.Get(ctx, fmt.Sprintf("%s/accounts/%s", settings.V.AccountsEngine.APIURL, id), httpclient.HttpConfig{
		Headers: MapString{
			"Authorization": fmt.Sprintf("Bearer %s", accountsengineservice.AccessToken.AccessToken),
		},
	}, &accountDocument)
	if err != nil {
		return nil, err
	}
	return &accountDocument, nil
}
