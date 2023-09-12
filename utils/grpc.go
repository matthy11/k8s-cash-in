package utils

import (
	"context"
	"crypto/x509"
	"fmt"
	"heypay-cash-in-server/settings"
	"os"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

var (
	nextTokenRefresh time.Time
	token            *oauth2.Token
)

const refreshWindow = 600000
const timeout = 3000

func GetConnection() (*grpc.ClientConn, error) {
	if didTokenExpire() {
		err := generateToken()
		if err != nil {
			return nil, err
		}
	}
	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	creds := credentials.NewClientTLSFromCert(certPool, "")
	perRPC := oauth.NewOauthAccess(token)
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:443", settings.V.GRPCAccountsEngine.APIURL),
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			ctx, cancel := context.WithTimeout(ctx, timeout*time.Millisecond)
			defer cancel()
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func didTokenExpire() bool {
	return token == nil || time.Now().After(nextTokenRefresh)
}

func generateToken() error {
	var clientOptions idtoken.ClientOption
	if settings.V.Env == "LOCAL" {
		clientOptions = idtoken.WithCredentialsFile("k8s/grpc-client-accounts-engine-service-account.json")
	} else {
		clientOptions = idtoken.WithCredentialsFile(os.Getenv("GRPC_CLIENT_ACCOUNTS_ENGINE_CREDENTIALS"))
	}
	tokenSource, err := idtoken.NewTokenSource(context.Background(), fmt.Sprintf("https://%s", settings.V.GRPCAccountsEngine.APIURL), clientOptions)
	if err != nil {
		return err
	}
	token, err = tokenSource.Token()
	if err != nil {
		return err
	}
	nextTokenRefresh = token.Expiry.Add(-refreshWindow * time.Millisecond)
	return nil
}
