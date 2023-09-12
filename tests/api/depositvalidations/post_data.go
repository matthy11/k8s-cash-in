package depositvalidations

import (
	"fmt"
	"heypay-cash-in-server/models/account"
	"heypay-cash-in-server/models/transferpaymentrequest"
	"heypay-cash-in-server/models/user"
	"heypay-cash-in-server/tests/utils"
)

var destinationRut = utils.GetRandomStringNumber(9)

var CorrectFullBody = map[string]interface{}{
	"amount":                   utils.GetRandomInt(1000, 10000),
	"destinationAccountNumber": fmt.Sprintf("00000000999%s", destinationRut[:len(destinationRut)-1]),
	"destinationAccountTypeId": "40",
	"destinationRut":           fmt.Sprintf("000%s", destinationRut),
	"message":                  "",
	"originAccountNumber":      fmt.Sprintf("0000000000%s", utils.GetRandomStringNumber(8)),
	"originAccountTypeId":      "20",
	"originBankId":             "0001",
	"originRut":                fmt.Sprintf("0000%s", destinationRut),
	"receivedAt":               "201222100000",
	"traceNumber":              "000000000001",
}

var mandatoryFields = []string{"amount", "destinationAccountNumber", "destinationAccountTypeId", "destinationRut", "message", "originAccountNumber", "originAccountTypeId", "originBankId", "originRut", "receivedAt", "traceNumber"}

var InvalidData, ValidData = utils.GetAllPossibleInvalidAndValidCreationData(CorrectFullBody, mandatoryFields, map[string]interface{}{}, utils.DataOptions{
	MaxInvalidItems: 10,
	MaxValidItems:   10,
})

var TransferPaymentRequestDocumentMock = &transferpaymentrequest.TransferPaymentRequest{
	ID: "qQFAUBFtQkBXSEGTdTgw",
}

var UserDocumentMock = user.User{
	ID:               "UXkQeUJsNuYncetFXEvJ",
	PrimaryAccountID: "NRbVslXrujlcQxUoXYcg",
}

var UserAccountDocumentMock = &account.Account{
	Balance:        0,
	BlockedBalance: 0,
	ID:             "NRbVslXrujlcQxUoXYcg",
	MaxBalance:     -1,
	OwnerType:      "user",
	Status:         "active",
}

var CommerceAccountDocumentMock = &account.Account{
	Balance:        0,
	BlockedBalance: 0,
	ID:             "NRbVslXrujlcQxUoXYcg",
	MaxBalance:     -1,
	OwnerType:      "commerce",
	Status:         "active",
}

var BlockedAccountDocumentMock = &account.Account{
	Status: "blocked",
}

var OverMaxBalanceAccountDocumentMock = &account.Account{
	Balance:    1000,
	MaxBalance: 1000,
	Status:     "active",
}
