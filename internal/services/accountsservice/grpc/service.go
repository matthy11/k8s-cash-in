package grpc

import (
	"context"
	"errors"
	"heypay-cash-in-server/models/account"
	"heypay-cash-in-server/utils"

	"google.golang.org/grpc"
)

type Client struct {
	client AccountsServiceClient
	conn   *grpc.ClientConn
}

var (
	GetAccount func(ctx context.Context, id string) (*account.Account, error)
)

func init() {
	GetAccount = getAccount
}

func NewClient() (*Client, error) {
	conn, err := utils.GetConnection()
	if err != nil {
		return nil, err
	}
	return &Client{
		client: NewAccountsServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *Client) Close() {
	_ = c.conn.Close()
}

func (c *Client) GetAccount(ctx context.Context, id string) (*account.Account, error) {
	response, err := c.client.GetDocument(ctx, &GetDocumentRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	if response.Code != "ok" {
		return nil, errors.New(response.Code)
	}
	return mapAccount(response.Data), nil
}

/*
This methods will automatically close the connection after getting a response
*/
func getAccount(ctx context.Context, id string) (*account.Account, error) {
	c, err := NewClient()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	return c.GetAccount(ctx, id)
}

func mapAccount(doc *Account) *account.Account {
	return &account.Account{
		AdditionalData: utils.MapStringToMapInterface(doc.GetAdditionalData()),
		Balance:        doc.GetBalance(),
		BlockedBalance: doc.GetBlockedBalance(),
		Category:       doc.GetCategory(),
		Channel:        doc.GetChannel(),
		CreatedAt:      utils.ParseStringISODate(doc.GetCreatedAt()),
		Currency:       doc.GetCurrency(),
		LastMovementAt: utils.ParseStringISODate(doc.GetLastMovementAt()),
		MaxBalance:     doc.GetMaxBalance(),
		OwnerID:        doc.GetOwnerId(),
		OwnerInfo: account.OwnerInfo{
			FirstName: doc.GetOwnerInfo().GetFirstName(),
			LastName:  doc.GetOwnerInfo().GetLastName(),
			Name:      doc.GetOwnerInfo().GetName(),
			Type:      doc.GetOwnerInfo().GetType(),
		},
		OwnerNationalID: doc.GetOwnerNationalId(),
		OwnerType:       doc.GetOwnerType(),
		Status:          doc.GetStatus(),
		TransactionID:   doc.GetTransactionId(),
		Type:            doc.GetType(),
		UpdatedAt:       utils.ParseStringISODate(doc.GetUpdatedAt()),
		UsageRules: account.UsageRules{
			AllowedFromAccount:   doc.UsageRules.GetAllowedFromAccount(),
			AllowedToAccount:     doc.UsageRules.GetAllowedToAccount(),
			ForbiddenFromAccount: doc.UsageRules.GetForbiddenFromAccount(),
			ForbiddenToAccount:   doc.UsageRules.GetForbiddenToAccount(),
		},
	}
}
