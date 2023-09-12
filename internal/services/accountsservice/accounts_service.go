package accountsservice

import (
	"context"
	"heypay-cash-in-server/internal/services/accountsservice/grpc"
	"heypay-cash-in-server/internal/services/accountsservice/http"
	"heypay-cash-in-server/models/account"
	"heypay-cash-in-server/utils"
	. "heypay-cash-in-server/utils/types"
)

var (
	GetAccount func(ctx context.Context, id string) (*account.Account, error)
)

func init() {
	GetAccount = getAccount
}

// getAccount gets an account from gRPC or http accounts engine service
func getAccount(ctx context.Context, id string) (*account.Account, error) {
	document, err := grpc.GetAccount(ctx, id)
	if err != nil {
		utils.Errormf(ctx, Map{"eventName": "grpc-error", "method": "GetAccount", "accountId": id}, "gRPC GetAccount(%s) failed with error %v, retrying with http protocol.", id, err)
		document, err = http.GetAccount(ctx, id)
		if err != nil {
			utils.Errormf(ctx, Map{"eventName": "http-error", "method": "GetAccount", "accountId": id}, "Http GetAccount(%s) failed with error %v.", id, err)
			return nil, err
		}
	}
	utils.Debugmf(ctx, Map{"eventName": "service-success", "method": "GetAccount", "accountId": id}, "GetAccount(%s): %+v", id, document)
	return document, nil
}
