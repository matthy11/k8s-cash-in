package settings

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
)

type Config struct {
	Version                 string
	Env                     string
	Project                 Project                                 `json:"project"`
	AppServerProject        AppProject                              `json:"app_server"`
	AccountsEngine          AccountsEngine                          `json:"accounts_engine"`
	GRPCAccountsEngine      GRPCAccountsEngine                      `json:"grpc_accounts_engine"`
	UserAccountCategory     map[string]UserAccountCategoryRules     `json:"user_account_category"`
	CommerceAccountCategory map[string]CommerceAccountCategoryRules `json:"commerce_account_category"`
}

type Project struct {
	ID string `json:id`
	BannedOriginRuts []string `json:"banned_origin_ruts"`
}

type AppProject struct {
	ID      string `json:"id"`
	AppName string `json:"app_name"`
	APIKey  string `json:"key"`
}

type AccountsEngine struct {
	APIURL     string `json:"api_url"`
	ClientID   string `json:"client_id"`
	PrivateKey string `json:"private_key"`
}

type GRPCAccountsEngine struct {
	APIURL string `json:"api_url"`
}

type UserAccountCategoryRules struct {
	MaxDailyBankTransferDepositAmount   int64 `json:"max_daily_bank_transfer_deposit_amount"`
	MaxMonthlyBankTransferDepositAmount int64 `json:"max_monthly_bank_transfer_deposit_amount"`
}

type CommerceAccountCategoryRules struct {
	MaxDailyBankTransferPaymentAmount   int64 `json:"max_daily_bank_transfer_payment_amount"`
	MaxMonthlyBankTransferPaymentAmount int64 `json:"max_monthly_bank_transfer_payment_amount"`
}

var V Config

func init() {
	V.Version = "1.3.1"
	V.Env = "LOCAL"
	if temp := os.Getenv("GO_ENV"); temp != "" {
		V.Env = strings.ToUpper(temp)
	}
	switch V.Env {
	case "LOCAL":
		_, filename, _, _ := runtime.Caller(0)
		dir := path.Join(path.Dir(filename), "..")
		err := os.Chdir(dir)
		if err != nil {
			panic(err)
		}
		jsonFile, err := os.Open("config/local.json")
		if err != nil {
			log.Fatalf("Fatal error loading local environment config - reason: %v", err)
		}
		data, err := ioutil.ReadAll(jsonFile)
		err = json.Unmarshal(data, &V)
		if err != nil {
			log.Fatalf("Fatal error decoding json local environment config - reason: %v", err)
		}
		defer jsonFile.Close()
	default:
		jsonData := os.Getenv("CONFIG")
		if len(jsonData) == 0 {
			log.Fatal("Fatal error CONFIG environment variable not configured")
		}
		err := json.Unmarshal([]byte(jsonData), &V)
		if err != nil {
			log.Fatalf("Fatal error decoding json environment config - reason: %v", err)
		}
	}
}
